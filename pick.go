package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runPick() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		osascriptAlert("Could not determine home directory.")
		return
	}
	historyFile := homeDir + "/.clipboard-history/history"

	data, err := os.ReadFile(historyFile)
	if os.IsNotExist(err) {
		osascriptAlert("No clipboard history yet. Start copying some text!")
		return
	}
	if err != nil {
		osascriptAlert("Could not read clipboard history.")
		return
	}

	var items []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		decoded, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			continue
		}
		items = append(items, string(decoded))
	}

	if len(items) == 0 {
		osascriptAlert("Clipboard history is empty.")
		return
	}

	var previews []string
	for i, item := range items {
		previews = append(previews, fmt.Sprintf("%2d. %s", i+1, safePreview(item)))
	}

	var asList strings.Builder
	asList.WriteString("{")
	for i, p := range previews {
		if i > 0 {
			asList.WriteString(", ")
		}
		fmt.Fprintf(&asList, `"%s"`, p)
	}
	asList.WriteString("}")

	script := fmt.Sprintf(`
activate
set choices to %s
set chosen to choose from list choices with prompt "Select item to paste:" with title "Clipboard History" default items {item 1 of choices}
if chosen is false then return ""
return item 1 of chosen
`, asList.String())

	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return
	}
	chosen := strings.TrimSpace(string(out))

	// Parse "N. preview text" to get index N
	parts := strings.SplitN(chosen, ".", 2)
	if len(parts) < 1 {
		return
	}
	idx, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return
	}
	idx-- // convert to 0-based
	if idx < 0 || idx >= len(items) {
		return
	}

	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(items[idx])
	if err := cmd.Run(); err != nil {
		return
	}

	exec.Command("osascript", "-e", `tell application "System Events" to keystroke "v" using command down`).Run()
}

func safePreview(text string) string {
	collapsed := strings.Join(strings.Fields(text), " ")
	runes := []rune(collapsed)
	if len(runes) > 72 {
		runes = runes[:72]
	}
	var b strings.Builder
	for _, r := range runes {
		if r >= 0x20 && r != 0x7F {
			b.WriteRune(r)
		}
	}
	s := b.String()
	s = strings.ReplaceAll(s, `\`, "/")
	s = strings.ReplaceAll(s, `"`, "'")
	return s
}

func osascriptAlert(msg string) {
	safe := strings.ReplaceAll(msg, `\`, "/")
	safe = strings.ReplaceAll(safe, `"`, "'")
	exec.Command("osascript", "-e", fmt.Sprintf(`display alert "%s"`, safe)).Run()
}
