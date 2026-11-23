package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var nameRule = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrInvalidTapName  = errors.New("invalid tap name")
	ErrTapExists       = errors.New("tap already exists")
)

const markerContent = "tapwright-tap 1\n"

var makeContentsFunc = makeContents

func CreateTap(rootPath, user, name string) error {
	if !nameRule.MatchString(user) {
		return fmt.Errorf("%w %q: use lowercase letters, digits and hyphens", ErrInvalidUserName, user)
	}

	if !nameRule.MatchString(name) {
		return fmt.Errorf("%w %q: use lowercase letters, digits and hyphens", ErrInvalidTapName, name)
	}

	if strings.HasPrefix(name, "homebrew-") {
		return fmt.Errorf("%w %q: tap name must not start with `homebrew-`", ErrInvalidTapName, name)
	}

	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return err
	}
	defer root.Close()

	if err := root.MkdirAll(user, 0o755); err != nil {
		return fmt.Errorf("create user directory %q: %w", filepath.Join(rootPath, user), err)
	}

	tapDir := filepath.Join(user, "homebrew-"+name)
	if err := root.Mkdir(tapDir, 0o755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("%w: %w", ErrTapExists, err)
		}
		return fmt.Errorf("create tap directory %q: %w", tapDir, err)
	}

	if err := makeContentsFunc(root, tapDir); err != nil {
		if rmErr := root.RemoveAll(tapDir); rmErr != nil {
			return errors.Join(err, rmErr)
		}
		return err
	}

	return nil
}

func makeContents(root *os.Root, tapDir string) error {
	caskDir := filepath.Join(tapDir, "Casks")
	if err := root.Mkdir(caskDir, 0o755); err != nil {
		return fmt.Errorf("create casks directory %q: %w", caskDir, err)
	}

	tapMarker := filepath.Join(tapDir, ".tapwright")
	if err := root.WriteFile(tapMarker, []byte(markerContent), 0o644); err != nil {
		return fmt.Errorf("create tap marker %q: %w", tapMarker, err)
	}

	return nil
}
