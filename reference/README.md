# Advent of Code references

Collected from 40 Firefox toolbar bookmarks on 2 November 2026: **51 downloaded repositories**, including the original GitLab project and three general references. One additional repository is unavailable. Profile bookmarks expand to their AoC repositories.

Each repository has one shallow checkout in `.repos/`. Year directories link to that checkout, so multi-year repositories appear in every year with source evidence. These links open the whole repository; each index also links to an example source file for its year. General references have no assigned AoC year.

| Year | Repositories |
| --- | --- |
| [2015](2015/README.md) | 14 |
| [2016](2016/README.md) | 18 |
| [2017](2017/README.md) | 23 |
| [2018](2018/README.md) | 25 |
| [2019](2019/README.md) | 18 |
| [2020](2020/README.md) | 17 |
| [2021](2021/README.md) | 18 |
| [2022](2022/README.md) | 21 |
| [2023](2023/README.md) | 20 |
| [2024](2024/README.md) | 16 |
| [2025](2025/README.md) | 14 |
| [general](general/README.md) | 3 |

## Sources and gaps

- [Repository discovery](discovery.md): profile expansion, non-GitHub pages, and stale bookmarks.
- [Manifest](manifest.json): every original bookmark, resolved repository, checked-out commit, and year evidence.
- `williewillus/advent-of-code-2016` returns 404.
- `anthonyf/advent-of-code` no longer contains the bookmarked 2016 source; its current 2023/2024 source is collected.
- Rustfinity is platform source in `general/`; its former Advent 2024 challenge source was not verified.
- Norvig’s unrelated `Jotto.ipynb`/`jotto.ipynb` files collide on this case-insensitive filesystem. Both remain available with `git show`; the Advent notebooks are readable.

## Restore downloads

`.repos/` is ignored by this project’s Git repository. The indexes, links, manifest, and restore script can be committed without vendoring third-party source. From a fresh checkout, run:

```sh
python3 reference/restore.py
```

The script downloads the recorded commits with up to three concurrent Git fetches and a four-minute timeout per Git command. It preserves existing checkouts and reports mismatched commits. It does not build or run downloaded code. Shallow clones omit earlier history; each upstream retains its own license.
