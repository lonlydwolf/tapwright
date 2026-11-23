package main

import (
	"errors"
	"regexp"
)

var nameRule = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrInvalidTapName  = errors.New("invalid tap name")
	ErrTapExists       = errors.New("tap already exists")
)

func CreateTap(root, user, name string) error {
	return nil
}
