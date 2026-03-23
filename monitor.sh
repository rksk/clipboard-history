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
        # base64-encode so multiline/special chars are safe in the flat history file
        encoded=$(printf '%s' "$current" | base64 | tr -d '\n')
        tmpfile=$(mktemp)
        # Put newest item first, remove duplicates, cap at MAX_ITEMS
        echo "$encoded" > "$tmpfile"
        grep -Fxv "$encoded" "$HISTORY_FILE" 2>/dev/null | head -n $((MAX_ITEMS - 1)) >> "$tmpfile"
        mv "$tmpfile" "$HISTORY_FILE"
    fi
    sleep 1
done
