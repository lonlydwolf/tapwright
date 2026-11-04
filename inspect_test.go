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
			name:  "valid url",
			input: "https://github.com/octocat/Hello-World",
			want: RepoInfo{
				Owner:        "octocat",
				Repository:   "Hello-World",
				CanonicalURL: "https://github.com/octocat/Hello-World",
			},
			wantErr: false,
		},

		{
			name:  "valid url",
			input: "https://github.com/octocat/repo_123.git",
			want: RepoInfo{
				Owner:        "octocat",
				Repository:   "repo_123.git",
				CanonicalURL: "https://github.com/octocat/repo_123.git",
			},
			wantErr: false,
		},

		{
			name:    "invalid url",
			input:   "https://github.com/octocat/repo!",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url",
			input:   "https://github.com/octocat/répo",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid url",
			input:   "https://github.com/octocat/repo~test",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (Percent-encoded path components)",
			input:   "https://github.com/octocat/%48ello",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (Percent-encoded path components)",
			input:   "https://github.com/%6Fctocat/repo",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (Percent-encoded path components)",
			input:   "https://github.com/octocat/repo%2Egit",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid scheme",
			input:   "http://github.com/octocat/Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid host",
			input:   "https://gitlab.com/octocat/Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (with query)",
			input:   "https://github.com/octocat/Hello-World?tab=readme",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (with fragmet)",
			input:   "https://github.com/octocat/Hello-World#readme",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (missing repository)",
			input:   "https://github.com/octocat/",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (missing owner)",
			input:   "https://github.com//Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (empty owner and repository)",
			input:   "https://github.com//",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (sub-page)",
			input:   "https://github.com/octocat/Hello-World/issues",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "embedded user credentials",
			input:   "https://user:pass@github.com/octocat/Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "empty query",
			input:   "https://github.com/octocat/Hello-World?",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "empty fragmet",
			input:   "https://github.com/octocat/Hello-World#",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid path component",
			input:   "https://github.com/./Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid path component",
			input:   "https://github.com/octocat/.",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid path component",
			input:   "https://github.com/../Hello-World",
			want:    RepoInfo{},
			wantErr: true,
		},
		{
			name:    "invalid path component",
			input:   "https://github.com/octocat/..",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (with port)",
			input:   "https://github.com:443/octocat/repo",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid path (trailing slashes)",
			input:   "https://github.com/octocat/repo/",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url (case sensitivity)",
			input:   "HTTPS://GitHub.COM/octocat/repo",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url",
			input:   "https:///octocat/repo",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url",
			input:   "https:///octocat/repo",
			want:    RepoInfo{},
			wantErr: true,
		},

		{
			name:    "invalid url",
			input:   "https://github.com//repo",
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
