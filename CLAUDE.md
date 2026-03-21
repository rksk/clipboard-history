# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

A minimal macOS clipboard history manager. Three scripts work together:

1. **`monitor.sh`** — background daemon that polls `pbpaste` every second, base64-encodes entries, deduplicates, and appends to `~/.clipboard-history/history` (capped at 50 items)
2. **`chooser.py`** — reads the history file, shows a native macOS `osascript` "choose from list" dialog, then restores the selected item to clipboard via `pbcopy` and auto-pastes with `Cmd+V`
3. **`build_service.py`** — installs a macOS Automator Quick Action (`~/.workflow/`) that triggers `chooser.py` from the Services menu

## Running the Scripts

```bash
# Start the monitor daemon manually
bash monitor.sh &

# Launch the chooser (normally triggered via Services menu)
python3 chooser.py

# Build/reinstall the Automator Quick Action
python3 build_service.py
```

The monitor is intended to run as a LaunchAgent (background daemon). The history file lives at `~/.clipboard-history/history` as newline-separated base64 strings.

## Architecture Notes

- No external dependencies — only macOS built-ins (`pbpaste`, `pbcopy`, `osascript`)
- Base64 encoding handles multiline text and special characters safely in the flat history file
- `chooser.py` collapses whitespace and truncates to 72 chars for display, but restores the full original content
- The Automator workflow plist structure is generated programmatically by `build_service.py`
