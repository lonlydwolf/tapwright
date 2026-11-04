package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type RepoInfo struct {
	Owner        string `json:"owner"`
	Repository   string `json:"repository"`
	CanonicalURL string `json:"canonical_url"`
}

func inspect(rawURL string) (RepoInfo, error) {
	parsedURL, err := url.Parse(rawURL)

	if err != nil {
		return RepoInfo{}, err
	}

	if parsedURL.Scheme != "https" {
		return RepoInfo{}, fmt.Errorf("invalid scheme %q: must be https", parsedURL.Scheme)
	}

	if parsedURL.User != nil {
		return RepoInfo{}, fmt.Errorf("URL must not include credentials")
	}

	if parsedURL.Host != "github.com" {
		return RepoInfo{}, fmt.Errorf("invalid host %q: must be github.com", parsedURL.Host)
	}

	if parsedURL.RawQuery != "" || parsedURL.ForceQuery {
		return RepoInfo{}, fmt.Errorf("URL must not include query")
	}

	if strings.Contains(rawURL, "#") {
		return RepoInfo{}, fmt.Errorf("URL must not include fragment")
	}

	if strings.Contains(parsedURL.EscapedPath(), "%") {
		return RepoInfo{}, fmt.Errorf("invalid path %q: must not include percent-encoded path components", parsedURL.EscapedPath())
	}

	pathParts := strings.Split(parsedURL.Path, "/")

	if len(pathParts) != 3 {
		return RepoInfo{}, fmt.Errorf("invalid path %q: must be /<owner>/<repository>", parsedURL.Path)
	}

	if pathParts[1] == "" || pathParts[1] == "." || pathParts[1] == ".." ||
		pathParts[2] == "" || pathParts[2] == "." || pathParts[2] == ".." {
		return RepoInfo{}, fmt.Errorf("invalid path %q: must be /<owner>/<repository>", parsedURL.Path)
	}

	var repoRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !repoRegex.MatchString(pathParts[1]) || !repoRegex.MatchString(pathParts[2]) {
		return RepoInfo{}, fmt.Errorf("inavlid path %q: each component may contain only ASCII letters, digits, hyphens, underscores, and periods", parsedURL.Path)
	}

	return RepoInfo{pathParts[1], pathParts[2], rawURL}, err
}
