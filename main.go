package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	root, exists := os.LookupEnv("TAPWRIGHT_TAPS_ROOT")
	if !exists {
		output, err := exec.Command("brew", "--taps").Output()
		if err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				fmt.Fprintln(os.Stderr, "Error: Homebrew not found in PATH")
			} else {
				fmt.Fprintf(os.Stderr, "Error: locate Homebrew taps directory: %v\n", err)

				if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
					if stderr := strings.TrimSpace(string(exitErr.Stderr)); stderr != "" {
						fmt.Fprintf(os.Stderr, "Homebrew stderr: %s\n", stderr)
					}
				}
			}
			os.Exit(1)
		}

		root = strings.TrimSpace(string(output))
	}

	if root == "" {
		fmt.Fprintln(os.Stderr, "Error: taps root cannot be empty")
		os.Exit(1)
	}

	info, err := os.Stat(root)
	created := errors.Is(err, os.ErrNotExist)

	if err != nil && !created {
		fmt.Fprintf(os.Stderr, "Error: check taps root %q: %v\n", root, err)
		os.Exit(1)
	}

	if err == nil && !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: taps root %q is not a directory\n", root)
		os.Exit(1)
	}

	if err := os.MkdirAll(root, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: create taps root %q: %v\n", root, err)
		os.Exit(1)
	}

	if created {
		fmt.Printf("Created taps root directory: %s\n", root)
	}

	if err := CreateTap(root, user, name); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, "Create:", filepath.Join(root, user, name))
}
