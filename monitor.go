package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	maxItems      = 50
	maxEntryBytes = 65536
)

func runMonitor() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "home dir: %v\n", err)
		os.Exit(1)
	}
	historyDir := homeDir + "/.clipboard-history"
	historyFile := historyDir + "/history"

	if err := os.MkdirAll(historyDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}
	os.Chmod(historyDir, 0700)

	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		f, err := os.OpenFile(historyFile, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create history file: %v\n", err)
			os.Exit(1)
		}
		f.Close()
	}
	os.Chmod(historyFile, 0600)

	// Clean up any orphaned tmp files on exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		cleanTmpFiles(historyDir)
		os.Exit(0)
	}()

	var last string
	for {
		out, _ := exec.Command("pbpaste").Output()
		current := string(out)

		if current != "" && current != last {
			if len(current) > maxEntryBytes {
				last = current
				time.Sleep(time.Second)
				continue
			}
			last = current

			encoded := base64.StdEncoding.EncodeToString([]byte(current))
			if err := writeHistory(historyDir, historyFile, encoded); err != nil {
				fmt.Fprintf(os.Stderr, "write history: %v\n", err)
			}
		}
		time.Sleep(time.Second)
	}
}

func writeHistory(historyDir, historyFile, encoded string) error {
	data, _ := os.ReadFile(historyFile)

	var newLines []string
	newLines = append(newLines, encoded)
	count := 1
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line == encoded {
			continue
		}
		if count >= maxItems {
			break
		}
		newLines = append(newLines, line)
		count++
	}

	f, err := os.CreateTemp(historyDir, ".tmp.")
	if err != nil {
		return err
	}
	tmpName := f.Name()
	f.Chmod(0600)

	renamed := false
	defer func() {
		if !renamed {
			os.Remove(tmpName)
		}
	}()

	if _, err := f.WriteString(strings.Join(newLines, "\n") + "\n"); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, historyFile); err != nil {
		return err
	}
	renamed = true
	os.Chmod(historyFile, 0600)
	return nil
}

func cleanTmpFiles(historyDir string) {
	entries, err := os.ReadDir(historyDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".tmp.") {
			os.Remove(historyDir + "/" + e.Name())
		}
	}
}
