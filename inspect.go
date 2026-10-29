package main

import (
	"fmt"
	"net/url"
	"strings"
)

type RepoInfo struct {
	owner        string
	repository   string
	canonicalURL string
}

func inspect(rawURL string) (RepoInfo, error) {
	parsedURL, err := url.Parse(rawURL)

	if err != nil {
		return RepoInfo{}, err
	}

	if parsedURL.Scheme != "https" {
		return RepoInfo{}, fmt.Errorf("invalid scheme %q: must be https", parsedURL.Scheme)
	}

	if parsedURL.Host != "github.com" {
		return RepoInfo{}, fmt.Errorf("invalid host %q: must be github.com", parsedURL.Host)
	}

	if parsedURL.RawQuery != "" {
		return RepoInfo{}, fmt.Errorf("URL must not include query")
	}

	if parsedURL.Fragment != "" {
		return RepoInfo{}, fmt.Errorf("URL must not include fragment")
	}

	pathParts := strings.Split(parsedURL.Path, "/")

	if len(pathParts) != 3 {
		return RepoInfo{}, fmt.Errorf("invalid path %q: must be /<owner>/<repository>", parsedURL.Path)
	}

	return RepoInfo{pathParts[1], pathParts[2], rawURL}, err
}
