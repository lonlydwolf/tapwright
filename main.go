package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: tapwright <command>")
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case "inspect":
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "Usage: tapwright inspect <argument>")
			os.Exit(1)
		}

		inspectCMD(args[1])

	case "create":
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, "Usage: tapwright create <user> <name>")
			os.Exit(1)
		}

		createCMD(args[1], args[2])

	default:
		fmt.Fprintln(os.Stderr, "Unknown command:", command)
		os.Exit(1)
	}
}

func inspectCMD(arg string) {

	if repoInfo, err := inspect(arg); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	} else if err := json.NewEncoder(os.Stdout).Encode(repoInfo); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

}

func createCMD(user, name string) {
	root := os.Getenv("TAPWRIGHT_TAPS_ROOT")
	if root == "" {
		root = "/opt/homebrew/Library/Taps"
	}

	if err := CreateTap(root, user, name); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, filepath.Join(root, user, "homebrew-"+name))
}
