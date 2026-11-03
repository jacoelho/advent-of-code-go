#!/usr/bin/env python3
import argparse
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import time


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    parser = argparse.ArgumentParser(description="Replay the matched 2019 Day 18 solver comparison.")
    parser.add_argument("output", type=Path, help="new output directory")
    parser.add_argument("--prepare-only", action="store_true", help="verify inputs and create overlays without Go execution")
    args = parser.parse_args()
    artifacts = Path(__file__).resolve().parent
    root = artifacts.parents[2]
    plan = json.loads((artifacts / "plan.json").read_text())
    for name, expected in json.loads((artifacts / "SHA256.json").read_text()).items():
        if digest(artifacts / name) != expected:
            raise SystemExit(f"artifact changed: {name}")
    for name in ("go.mod", "go.sum"):
        if digest(root / name) != digest(artifacts / (name + ".txt")):
            raise SystemExit(f"module changed: {name}; review before a new comparison")
    if digest(root / "inputs/2019/18.txt") != plan["fixture_sha256"]:
        raise SystemExit("actual fixture changed")
    output = args.output.resolve()
    output.mkdir(parents=True, exist_ok=False)
    package = root / "internal/aoc2019"
    for side in ("baseline", "candidate"):
        replacements = {
            str(package / "day18.go"): str(artifacts / f"{side}.go.txt"),
            str(package / "day18_test.go"): str(artifacts / f"{side}_test.go.txt"),
            str(package / "zz_day18_performance_test.go"): str(artifacts / "benchmark_test.go.txt"),
        }
        (output / f"{side}-overlay.json").write_text(json.dumps({"Replace": replacements}, indent=2) + "\n")
    if args.prepare_only:
        print(output)
        return
    env = os.environ.copy()
    env.update(GOCACHE="/tmp/aoc-performance-cache", GOMODCACHE="/tmp/aoc-performance-modcache")
    records = []

    def run(command, cwd, timeout):
        started = time.monotonic()
        result = subprocess.run(command, cwd=cwd, env=env, text=True, capture_output=True, timeout=timeout)
        records.append({"args": command, "cwd": str(cwd), "elapsed_seconds": time.monotonic() - started,
                        "returncode": result.returncode, "stdout": result.stdout, "stderr": result.stderr})
        (output / "commands.json").write_text(json.dumps(records, indent=2) + "\n")
        if result.returncode:
            raise SystemExit(f"command failed; evidence retained in {output}; no retry")
        return result

    run(["go", "version"], root, 30)
    run(["go", "env", "GOOS", "GOARCH", "GOARM64", "GOFLAGS", "CGO_ENABLED", "GOTOOLCHAIN"], root, 30)
    (output / "runtime-env.json").write_text(json.dumps({name: env.get(name, "<unset>") for name in
        ("GOMAXPROCS", "GOGC", "GOMEMLIMIT", "GODEBUG")}, indent=2) + "\n")
    for side in ("baseline", "candidate"):
        run(["go", "test", "-c", "-pgo=off", "-overlay", str(output / f"{side}-overlay.json"),
             "-o", str(output / f"{side}.test"), "./internal/aoc2019"], root, 60)
    started = time.monotonic()
    for round_number in range(1, plan["samples_per_side"] + 1):
        order = ("baseline", "candidate") if round_number % 2 else ("candidate", "baseline")
        for side in order:
            remaining = plan["budget_seconds"] - (time.monotonic() - started)
            if remaining <= 0:
                raise SystemExit("measurement budget exhausted; partial results retained")
            result = run([str(output / f"{side}.test"), "-test.run=^$", "-test.bench=^BenchmarkDay18$",
                          "-test.benchmem", "-test.benchtime=100ms", "-test.count=1", "-test.cpu=1",
                          "-test.timeout=60s"], package, min(65, remaining))
            with (output / f"{side}.txt").open("a") as stream:
                stream.write(result.stdout)
            if result.stdout.count("BenchmarkDay18/") != 4:
                raise SystemExit("benchmark selector did not run all four workloads")
            print(round_number, side, flush=True)
    if benchstat := shutil.which("benchstat"):
        result = run([benchstat, str(output / "baseline.txt"), str(output / "candidate.txt")], root, 30)
        (output / "benchstat.txt").write_text(result.stdout)
        print(result.stdout)
    else:
        print("benchstat unavailable; raw samples retained")


if __name__ == "__main__":
    main()
