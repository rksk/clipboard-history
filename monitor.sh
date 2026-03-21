#!/bin/bash
# Clipboard history monitor — runs in background via LaunchAgent
# Polls pbpaste every second; stores up to MAX_ITEMS in ~/.clipboard-history/history

HISTORY_FILE="$HOME/.clipboard-history/history"
MAX_ITEMS=50

mkdir -p "$(dirname "$HISTORY_FILE")"
touch "$HISTORY_FILE"

last=""
while true; do
    current=$(pbpaste 2>/dev/null)
    if [[ -n "$current" && "$current" != "$last" ]]; then
        last="$current"
        # Encode to base64 (single line) so we can safely store multiline content
        encoded=$(printf '%s' "$current" | python3 -c \
            "import sys, base64; print(base64.b64encode(sys.stdin.buffer.read()).decode())")
        tmpfile=$(mktemp)
        # Put newest item first, remove duplicates, cap at MAX_ITEMS
        echo "$encoded" > "$tmpfile"
        grep -Fxv "$encoded" "$HISTORY_FILE" 2>/dev/null | head -n $((MAX_ITEMS - 1)) >> "$tmpfile"
        mv "$tmpfile" "$HISTORY_FILE"
    fi
    sleep 1
done
