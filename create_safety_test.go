package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateTapRejectsUserSymlinkOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "alice")); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := CreateTap(root, "alice", "tools"); err == nil {
		t.Fatal("CreateTap() error = nil, want an error for a user directory that escapes the root")
	}

	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatalf("read %s: %v", outside, err)
	}
	for _, e := range entries {
		t.Errorf("CreateTap() wrote %q outside the root", e.Name())
	}
}

func TestCreateTapRollsBackInsideRootOnly(t *testing.T) {
	root := t.TempDir()

	cwd := t.TempDir()
	decoy := filepath.Join(cwd, "alice", "homebrew-tools", "keep.txt")
	if err := os.MkdirAll(filepath.Dir(decoy), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := os.WriteFile(decoy, []byte("keep"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	t.Chdir(cwd)

	errBoom := errors.New("boom")
	orig := makeContentsFunc
	makeContentsFunc = func(r *os.Root, tapDir string) error {
		if err := r.Mkdir(filepath.Join(tapDir, "Casks"), 0o755); err != nil {
			t.Fatalf("partial build: %v", err)
		}
		return errBoom
	}
	t.Cleanup(func() { makeContentsFunc = orig })

	err := CreateTap(root, "alice", "tools")
	if !errors.Is(err, errBoom) {
		t.Fatalf("CreateTap() error = %v, want %v", err, errBoom)
	}

	tapDir := filepath.Join(root, "alice", "homebrew-tools")
	if _, err := os.Stat(tapDir); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("stat %s: err = %v, want the partial tap removed", tapDir, err)
	}
	if info, err := os.Stat(filepath.Join(root, "alice")); err != nil || !info.IsDir() {
		t.Errorf("user directory: err = %v, want it kept", err)
	}
	if got, err := os.ReadFile(decoy); err != nil || string(got) != "keep" {
		t.Errorf("decoy %s: got %q, err = %v; rollback touched the working directory", decoy, got, err)
	}
}
