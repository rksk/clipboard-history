#!/usr/bin/env python3
"""
Clipboard History Chooser
Shows a native macOS 'choose from list' dialog with clipboard history.
The selected item is copied to clipboard and pasted into the frontmost app.
"""

import base64
import os
import subprocess
import sys

HISTORY_FILE = os.path.expanduser("~/.clipboard-history/history")


def osascript(script: str) -> str:
    result = subprocess.run(["osascript", "-e", script],
                            capture_output=True, text=True)
    return result.stdout.strip()


def alert(msg: str):
    osascript(f'display alert "{msg}"')


# ── Load history ──────────────────────────────────────────────────────────────
if not os.path.exists(HISTORY_FILE):
    alert("No clipboard history yet. Start copying some text!")
    sys.exit(0)

with open(HISTORY_FILE) as f:
    lines = [l.strip() for l in f if l.strip()]

if not lines:
    alert("Clipboard history is empty.")
    sys.exit(0)

# Decode each base64 entry
items: list[str] = []
for line in lines:
    try:
        decoded = base64.b64decode(line).decode("utf-8", errors="replace")
        items.append(decoded)
    except Exception:
        continue

if not items:
    alert("Could not read clipboard history.")
    sys.exit(0)

# ── Build numbered preview list ───────────────────────────────────────────────
previews: list[str] = []
for i, item in enumerate(items, 1):
    # Collapse whitespace for display; escape double-quotes for AppleScript
    preview = " ".join(item.split())[:72].replace("\\", "\\\\").replace('"', '\\"')
    previews.append(f"{i:>2}. {preview}")

as_list = "{" + ", ".join(f'"{p}"' for p in previews) + "}"

script = f"""
activate
set choices to {as_list}
set chosen to choose from list choices ¬
    with prompt "Select item to paste:" ¬
    with title "Clipboard History" ¬
    default items {{item 1 of choices}}
if chosen is false then return ""
return item 1 of chosen
"""

chosen = osascript(script)

if not chosen:
    sys.exit(0)

# ── Parse index and restore full content ──────────────────────────────────────
try:
    idx = int(chosen.split(".")[0].strip()) - 1
    selected = items[idx]
except (ValueError, IndexError):
    sys.exit(1)

# Copy full content back to clipboard
subprocess.run("pbcopy", input=selected.encode("utf-8"), check=True)

# Paste into the frontmost app
osascript('tell application "System Events" to keystroke "v" using command down')
