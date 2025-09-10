# Repository Mirror Archive Action

A GitHub Action that clones the current repository with `--mirror`, creates a tar archive, and compresses it with zstd.

## Usage

```yaml
name: Mirror and Archive

on:
  workflow_dispatch:

jobs:
  archive:
    runs-on: ubuntu-latest
    steps:
      - name: Mirror and Archive
        uses: ./
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}