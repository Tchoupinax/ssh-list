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
