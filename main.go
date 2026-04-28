package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: clipboard-history <monitor|pick|build-service>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "monitor":
		runMonitor()
	case "pick":
		runPick()
	case "build-service":
		runBuildService()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
