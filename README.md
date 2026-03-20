![Logo](./.github/key-icon.svg)
# ssh-list

## Installation

### macOS

```bash
brew tap tchoupinax/brew
brew install ssh-list
```

Tutorial : https://betterprogramming.pub/a-step-by-step-guide-to-create-homebrew-taps-from-github-repos-f33d3755ba74

## Render

```
→ ssh-list --stats workload

  SSH connections

     Alias         User         Identity file                          Host                      CPU      RAM
  ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
  0  workload-1    tchoupinax   ~/.ssh-keys/proxmox-vms                10.210.0.48               1.81/2   77% 7.3G/9.5G
  1  workload-10   tchoupinax   ~/.ssh-keys/proxmox-vms                10.210.0.23               0.43/4   40% 3.8G/9.5G
  2  workload-2    tchoupinax   ~/.ssh-keys/proxmox-vms                10.210.0.12               3.29/2   80% 6.1G/7.6G
  3  workload-3    tchoupinax   ~/.ssh-keys/proxmox-vms                10.210.0.14               0.73/2   74% 5.6G/7.6G
  4  workload-4    tchoupinax   ~/.ssh-keys/proxmox-vms                10.210.0.10               0.89/4   57% 10.9G/19.1G
  5  workload-6    tchoupinax   ~/.ssh-keys/proxmox-vms                10.210.0.13               0.87/2   93% 7.0G/7.6G
  6  workload-8    tchoupinax   ~/.ssh-keys/proxmox-vms                10.210.0.21               0.20/7   51% 3.4G/6.6G
  7  workload-9    tchoupinax   ~/.ssh-keys/proxmox-vms                10.210.0.15               0.18/7   87% 5.8G/6.6G
  Metrics fetched in 234ms

performed in 235.543125ms
```

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
