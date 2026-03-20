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

Hosts that are down, non‑Linux, or unreachable show `—` in those columns. Queries run in parallel (up to 8 at a time) with a timeout per host.
