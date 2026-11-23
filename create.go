package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var nameRule = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrInvalidTapName  = errors.New("invalid tap name")
	ErrTapExists       = errors.New("tap already exists")
)

const markerContent = "tapwright-tap 1\n"

func CreateTap(root, user, name string) error {
	if !nameRule.MatchString(user) {
		return fmt.Errorf("%w %q: use lowercase letters, digits and hyphens", ErrInvalidUserName, user)
	}

	if !nameRule.MatchString(name) {
		return fmt.Errorf("%w %q: use lowercase letters, digits and hyphens", ErrInvalidTapName, name)
	}

	userDir := filepath.Join(root, user)
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		return fmt.Errorf("create user directory %q: %w", userDir, err)
	}

	tapDir := filepath.Join(userDir, "homebrew-"+name)
	if err := os.Mkdir(tapDir, 0o755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("%w: %w", ErrTapExists, err)
		}
		return fmt.Errorf("create tap directory %q: %w", tapDir, err)
	}

	if err := makeContents(tapDir); err != nil {
		if rmErr := os.RemoveAll(tapDir); rmErr != nil {
			return errors.Join(err, rmErr)
		}
		return err
	}

	return nil
}

func makeContents(tapDir string) error {
	caskDir := filepath.Join(tapDir, "Casks")
	if err := os.Mkdir(caskDir, 0o755); err != nil {
		return fmt.Errorf("create casks directory %q: %w", caskDir, err)
	}

	tapMarker := filepath.Join(tapDir, ".tapwright")
	if err := os.WriteFile(tapMarker, []byte(markerContent), 0o644); err != nil {
		return fmt.Errorf("create tap marker %q: %w", tapMarker, err)
	}

	return nil
}
