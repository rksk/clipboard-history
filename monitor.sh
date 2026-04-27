#!/bin/bash
# Clipboard history monitor — runs in background via LaunchAgent
# Polls pbpaste every second; stores up to MAX_ITEMS in ~/.clipboard-history/history

HISTORY_DIR="$HOME/.clipboard-history"
HISTORY_FILE="$HISTORY_DIR/history"
MAX_ITEMS=50
MAX_ENTRY_BYTES=65536

mkdir -p -m 700 "$HISTORY_DIR"
[ -f "$HISTORY_FILE" ] || install -m 600 /dev/null "$HISTORY_FILE"
[ -f "$HISTORY_DIR/monitor.log" ] || install -m 600 /dev/null "$HISTORY_DIR/monitor.log"

last=""
while true; do
    current=$(pbpaste 2>/dev/null)
    if [[ -n "$current" && "$current" != "$last" ]]; then
        byte_count=$(printf '%s' "$current" | wc -c | tr -d ' ')
        if (( byte_count > MAX_ENTRY_BYTES )); then
            last="$current"
            sleep 1
            continue
        fi
        last="$current"
        # base64-encode so multiline/special chars are safe in the flat history file
        encoded=$(printf '%s' "$current" | base64 | tr -d '\n')
        tmpfile=$(mktemp "$HISTORY_DIR/.tmp.XXXXXX")
        chmod 600 "$tmpfile"
        # Put newest item first, remove duplicates, cap at MAX_ITEMS
        echo "$encoded" > "$tmpfile"
        grep -Fxv "$encoded" "$HISTORY_FILE" 2>/dev/null | head -n $((MAX_ITEMS - 1)) >> "$tmpfile"
        mv "$tmpfile" "$HISTORY_FILE"
        chmod 600 "$HISTORY_FILE"
    fi
    sleep 1
done
