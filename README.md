# Clipboard History

A lightweight macOS clipboard history manager. Stores up to 50 clipboard entries and lets you browse and paste them via a native macOS dialog. 

No third-party dependencies, no cloud sync, just local scripts using built-in tools.

## Requirements

- macOS 12 or later
- Python 3 (included with macOS)

## Installation

### 1. Clone the repo

```bash
git clone <repo-url> ~/.clipboard-history
```

### 2. Install the Services menu trigger

```bash
python3 ~/.clipboard-history/build_service.py
```

This creates an Automator Quick Action at `~/Library/Services/Clipboard History.workflow`.

### 3. Set up the background monitor

Copy the LaunchAgent plist, replacing `YOUR_USERNAME` with your actual username:

```bash
sed "s/YOUR_USERNAME/$(whoami)/g" \
    ~/.clipboard-history/com.clipboard-history.monitor.plist \
    > ~/Library/LaunchAgents/com.clipboard-history.monitor.plist
```

Load it so it starts now and on every login:

```bash
launchctl load ~/Library/LaunchAgents/com.clipboard-history.monitor.plist
```

### 4. Grant permissions

On first use macOS will prompt for:

- **Accessibility** — allows `chooser.py` to simulate `Cmd+V` to paste
- **Automation** — allows `osascript` to control System Events

You can also grant these in advance under **System Settings → Privacy & Security**.

### 5. Assign a keyboard shortcut (recommended)

1. Open **System Settings → Keyboard → Keyboard Shortcuts**
2. Select **Services** in the left sidebar
3. Scroll to find **Clipboard History** (under General)
4. Double-click the empty space to the right of it
5. Press your desired shortcut (e.g. `⌘ + ⇧ + V`)
6. Click **Done**

If the shortcut doesn't trigger, another app may have a conflicting binding — try a different combination.

## Usage

Trigger via the assigned keyboard shortcut, or from any app's **Services** menu. A dialog lists your recent clipboard items — select one to paste it into the active app.

## How it works

Three scripts work together:

- **`monitor.sh`** — background daemon that polls `pbpaste` every second, base64-encodes entries, deduplicates, and appends to `~/.clipboard-history/history` (capped at 50 items)
- **`chooser.py`** — reads the history file, shows a native macOS `osascript` "choose from list" dialog, then restores the selected item to clipboard via `pbcopy` and auto-pastes with `Cmd+V`
- **`build_service.py`** — generates the Automator Quick Action plist at `~/Library/Services/Clipboard History.workflow`

No external dependencies — only macOS built-ins (`pbpaste`, `pbcopy`, `osascript`). Base64 encoding handles multiline text and special characters safely in the flat history file. `chooser.py` collapses whitespace and truncates to 72 chars for display, but restores the full original content.

To run the scripts directly without the Services menu:

```bash
bash monitor.sh &       # start the monitor daemon
python3 chooser.py      # launch the chooser dialog
```

## Uninstall

```bash
launchctl unload ~/Library/LaunchAgents/com.clipboard-history.monitor.plist
rm ~/Library/LaunchAgents/com.clipboard-history.monitor.plist
rm -rf ~/Library/Services/Clipboard\ History.workflow
rm -rf ~/.clipboard-history
```
