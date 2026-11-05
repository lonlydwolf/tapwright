package main

import "testing"

type testCase struct {
	name    string
	input   string
	want    RepoInfo
	wantErr bool
}

func TestInspect(t *testing.T) {
	tests := []testCase{
		{
			name:  "valid url / standard repo",
			input: "https://github.com/octocat/Hello-World",
			want: RepoInfo{
				Owner:        "octocat",
				Repository:   "Hello-World",
				CanonicalURL: "https://github.com/octocat/Hello-World",
			},
			wantErr: false,
		},
		{
			name:  "valid url / repo with extension",
			input: "https://github.com/octocat/repo_123.git",
			want: RepoInfo{
				Owner:        "octocat",
				Repository:   "repo_123.git",
				CanonicalURL: "https://github.com/octocat/repo_123.git",
			},
			wantErr: false,
		},
		{
			name:    "invalid repo / special characters exclamation",
			input:   "https://github.com/octocat/repo!",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid repo / unicode character",
			input:   "https://github.com/octocat/répo",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid repo / tilde character",
			input:   "https://github.com/octocat/repo~test",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path / percent-encoded repository",
			input:   "https://github.com/octocat/%48ello",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path / percent-encoded owner",
			input:   "https://github.com/%6Fctocat/repo",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path / percent-encoded dot in repository",
			input:   "https://github.com/octocat/repo%2Egit",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid scheme / http instead of https",
			input:   "http://github.com/octocat/Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid host / gitlab instead of github",
			input:   "https://gitlab.com/octocat/Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url / with query parameters",
			input:   "https://github.com/octocat/Hello-World?tab=readme",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url / with fragment",
			input:   "https://github.com/octocat/Hello-World#readme",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path / missing repository",
			input:   "https://github.com/octocat/",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path / missing owner",
			input:   "https://github.com//Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path / empty owner and repository",
			input:   "https://github.com//",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path / sub-page path",
			input:   "https://github.com/octocat/Hello-World/issues",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url / embedded user credentials",
			input:   "https://user:pass@github.com/octocat/Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url / empty query",
			input:   "https://github.com/octocat/Hello-World?",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url / empty fragment",
			input:   "https://github.com/octocat/Hello-World#",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path component / dot owner",
			input:   "https://github.com/./Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path component / dot repository",
			input:   "https://github.com/octocat/.",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path component / double-dot owner",
			input:   "https://github.com/../Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path component / double-dot repository",
			input:   "https://github.com/octocat/..",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url / explicit port number",
			input:   "https://github.com:443/octocat/repo",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path / trailing slash",
			input:   "https://github.com/octocat/repo/",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url / uppercase scheme and host",
			input:   "HTTPS://GitHub.COM/octocat/repo",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url / empty host missing authority",
			input:   "https:///octocat/repo",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path / empty owner segment",
			input:   "https://github.com//repo",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid scheme / uppercase HTTPS",
			input:   "HTTPS://github.com/octocat/Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid scheme / mixed-case HtTpS",
			input:   "HtTpS://github.com/octocat/Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url / malformed percent escape",
			input:   "https://github.com/octocat/repo%GG",
			want:    RepoInfo{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := inspect(tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("inspect() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if got != tc.want {
				t.Errorf("inspect() got = %v, want %v", got, tc.want)
			}
		})
	}
}
