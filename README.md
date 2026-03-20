![Logo](./.github/key-icon.svg)
# ssh-list

## Installation

### macOS

```bash
brew tap tchoupinax/brew
brew install ssh-list
```

Tutorial : https://betterprogramming.pub/a-step-by-step-guide-to-create-homebrew-taps-from-github-repos-f33d3755ba74

## Features

### Include file

SSH native config allows to include folder or file dynamically. `ssh-list` handle this and also read config from any files found.

```bash
Include ~/.ssh/dynamic-config/*
```

### Remote CPU & RAM (`--stats`)

When listing hosts, you can ask `ssh-list` to connect to each server over SSH and read Linux `/proc` metrics (load average, cores, memory). This uses the same keys as interactive SSH.

```bash
ssh-list --stats
```

- **CPU** column: 1‑minute load average / logical CPU count (e.g. `0.42/4`).
- **RAM** column: approximate used percentage and used/total size (from `MemAvailable` when present).

Entries whose **alias** or **hostname** contains `git` (case-insensitive, e.g. GitHub/GitLab hosts) **do not** get a metrics SSH: CPU/RAM show `–` for those rows.

Other failures (host down, non‑Linux, timeout) show `⚠`; missing values use `∅`. Queries run in parallel (up to 8 at a time) with a timeout per host.

The table appears **immediately** with `…` in the CPU/RAM columns while each host is queried; the view refreshes as results arrive (alternate screen on a real TTY). When everything is done, the final table is printed again on the main screen so it stays visible, followed by how long the metric fetch took. If stdout is not a terminal, the tool falls back to waiting for all metrics before printing once.
