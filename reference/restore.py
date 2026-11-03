#!/usr/bin/env python3
"""Restore the reference checkouts at the commits recorded in manifest.json."""

import concurrent.futures
import json
import os
from pathlib import Path
import subprocess
import tempfile


ROOT = Path(__file__).resolve().parent
ENV = dict(os.environ, GIT_TERMINAL_PROMPT="0", GIT_LFS_SKIP_SMUDGE="1")


def git(directory, *arguments):
    return subprocess.run(
        ["git", "-c", "credential.helper=", "-c", "core.hooksPath=/dev/null",
         "-c", "core.fsmonitor=false", "-C", str(directory), *arguments],
        env=ENV, check=True, capture_output=True, text=True, timeout=240,
    ).stdout.strip()


def restore(repository):
    name = repository["name"]
    checkout = ROOT / repository["checkout"]
    try:
        if checkout.exists():
            if git(checkout, "rev-parse", "HEAD") != repository["commit"]:
                raise ValueError("existing checkout has a different commit; left unchanged")
        else:
            # Publish only complete checkouts; failed fetches leave no partial target.
            with tempfile.TemporaryDirectory(prefix=".fetch-", dir=checkout.parent) as work:
                staging = Path(work) / "repository"
                staging.mkdir()
                git(staging, "init", "--quiet")
                git(staging, "remote", "add", "origin", repository["url"] + ".git")
                git(staging, "fetch", "--quiet", "--depth=1", "origin", repository["commit"])
                git(staging, "checkout", "--quiet", "--detach", repository["commit"])
                staging.rename(checkout)
        for relative in repository["links"]:
            link = ROOT / relative
            link.parent.mkdir(exist_ok=True)
            expected = os.path.relpath(checkout, link.parent)
            if link.is_symlink():
                if os.readlink(link) != expected:
                    raise ValueError(f"unexpected link {relative}; left unchanged")
            elif link.exists():
                raise ValueError(f"path already exists: {relative}; left unchanged")
            else:
                link.symlink_to(expected, target_is_directory=True)
        return True, f"OK {name}"
    except (OSError, ValueError, subprocess.SubprocessError) as error:
        detail = error.stderr.strip() if isinstance(error, subprocess.CalledProcessError) else str(error)
        return False, f"FAILED {name}: {detail}"


def main():
    manifest = json.loads((ROOT / "manifest.json").read_text())
    repositories = [r for r in manifest["repositories"] if r["status"] == "collected"]
    (ROOT / ".repos").mkdir(exist_ok=True)
    failed = False
    with concurrent.futures.ThreadPoolExecutor(max_workers=3) as pool:
        for success, message in pool.map(restore, repositories):
            print(message, flush=True)
            failed |= not success
    return int(failed)


if __name__ == "__main__":
    raise SystemExit(main())
