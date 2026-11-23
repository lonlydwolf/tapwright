package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type createTestCase struct {
	name       string
	user       string
	tap        string
	setup      func(t *testing.T, root string)
	wantErr    error
	wantFiles  map[string]string
	wantDirs   []string
	wantAbsent []string
}

const wantMarker = "tapwright-tap 1\n"

func TestCreateTap(t *testing.T) {
	tests := []createTestCase{
		{
			name: "success",
			user: "lonlydwolf",
			tap:  "tap",
			wantFiles: map[string]string{
				"lonlydwolf/homebrew-tap/.tapwright": wantMarker,
			},
			wantDirs: []string{
				"lonlydwolf",
				"lonlydwolf/homebrew-tap",
				"lonlydwolf/homebrew-tap/Casks",
			},
		},

		{
			name: "success - single lowercase letter",
			user: "a",
			tap:  "a",
			wantFiles: map[string]string{
				"a/homebrew-a/.tapwright": wantMarker,
			},
			wantDirs: []string{
				"a",
				"a/homebrew-a",
				"a/homebrew-a/Casks",
			},
		},

		{
			name: "success - digits only",
			user: "123",
			tap:  "123",
			wantFiles: map[string]string{
				"123/homebrew-123/.tapwright": wantMarker,
			},
			wantDirs: []string{
				"123",
				"123/homebrew-123",
				"123/homebrew-123/Casks",
			},
		},

		{
			name: "success - multiple hyphen-separated segments",
			user: "alice-2-team",
			tap:  "tools-2-extra",
			wantFiles: map[string]string{
				"alice-2-team/homebrew-tools-2-extra/.tapwright": wantMarker,
			},
			wantDirs: []string{
				"alice-2-team",
				"alice-2-team/homebrew-tools-2-extra",
				"alice-2-team/homebrew-tools-2-extra/Casks",
			},
		},

		{
			name: "success - letters and digits mixed in one segment",
			user: "alice2",
			tap:  "v2-tools",
			wantFiles: map[string]string{
				"alice2/homebrew-v2-tools/.tapwright": wantMarker,
			},
			wantDirs: []string{
				"alice2",
				"alice2/homebrew-v2-tools",
				"alice2/homebrew-v2-tools/Casks",
			},
		},

		{
			name:    "invalid user - dot",
			user:    "alice.team",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - empty user",
			user:    "",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - uppercase letter",
			user:    "Alice",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - leading hyphen",
			user:    "-alice",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - trailing hyphen",
			user:    "alice-",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - consecutive hyphens",
			user:    "alice--team",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - underscore",
			user:    "alice_team",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - trailing space",
			user:    "alice ",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - trailing newline",
			user:    "alice\n",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - non-ASCII letter",
			user:    "álice",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user - embedded slash",
			user:    "alice/tools",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},
		{
			name:       "invalid user - parent directory traversal",
			user:       "../alice",
			tap:        "tap",
			wantErr:    ErrInvalidUserName,
			wantAbsent: []string{"../alice"},
		},

		{
			name:    "invalid user - special character",
			user:    "al!ce",
			tap:     "tap",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid tap - dot",
			user:    "alice",
			tap:     "tools.v2",
			wantErr: ErrInvalidTapName,
		},

		{
			name:    "invalid tap - empty tap name",
			user:    "alice",
			tap:     "",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - uppercase letter",
			user:    "alice",
			tap:     "Tap",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - leading hyphen",
			user:    "alice",
			tap:     "-tap",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - trailing hyphen",
			user:    "alice",
			tap:     "tap-",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - consecutive hyphens",
			user:    "alice",
			tap:     "dev--tools",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - underscore",
			user:    "alice",
			tap:     "dev_tools",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - embedded space",
			user:    "alice",
			tap:     "dev tools",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - trailing newline",
			user:    "alice",
			tap:     "tap\n",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - non-ASCII letter",
			user:    "alice",
			tap:     "tóols",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - slash",
			user:    "alice",
			tap:     "dev/tools",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - special character",
			user:    "alice",
			tap:     "t@p",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid tap - `homwbrew-` prefix",
			user:    "alice",
			tap:     "homebrew-tools",
			wantErr: ErrInvalidTapName,
			wantAbsent: []string{
				"alice",
			},
		},

		{
			name:    "invalid user and tap - both empty, user error wins",
			user:    "",
			tap:     "",
			wantErr: ErrInvalidUserName,
		},

		{
			name:    "invalid user and tap - both malformed, user error wins",
			user:    "Alice",
			tap:     "t@p",
			wantErr: ErrInvalidUserName,
		},

		{
			name: "existing empty target",
			user: "alice",
			tap:  "tools",
			setup: func(t *testing.T, root string) {
				tapDir := filepath.Join(root, "alice", "homebrew-tools")
				if err := os.MkdirAll(tapDir, 0o755); err != nil {
					t.Fatalf("setup: %v", err)
				}
			},
			wantErr: ErrTapExists,
			wantDirs: []string{
				"alice",
				"alice/homebrew-tools",
			},
		},

		{
			name: "existing target",
			user: "alice",
			tap:  "tools",
			setup: func(t *testing.T, root string) {
				userDir := filepath.Join(root, "alice")
				if err := os.MkdirAll(userDir, 0o755); err != nil {
					t.Fatalf("setup: %v", err)
				}

				tapDir := filepath.Join(userDir, "homebrew-tools")
				if err := os.Mkdir(tapDir, 0o755); err != nil {
					t.Fatalf("setup: %v", err)
				}

				dataPath := filepath.Join(tapDir, "keep.txt")
				dataContent := []byte("original")
				if err := os.WriteFile(dataPath, dataContent, 0o644); err != nil {
					t.Fatalf("setup: %v", err)
				}
			},
			wantErr: ErrTapExists,
			wantFiles: map[string]string{
				"alice/homebrew-tools/keep.txt": "original",
			},
			wantDirs: []string{
				"alice",
				"alice/homebrew-tools",
			},
			wantAbsent: []string{
				"alice/homebrew-tools/Casks",
				"alice/homebrew-tools/.tapwright",
			},
		},

		{
			name: "user directory already exists (with 1 tap)",
			user: "lonlydwolf",
			tap:  "tools",
			setup: func(t *testing.T, root string) {
				userDir := filepath.Join(root, "lonlydwolf")
				if err := os.MkdirAll(userDir, 0o755); err != nil {
					t.Fatalf("setup: %v", err)
				}

				tapDir := filepath.Join(userDir, "homebrew-tap")
				if err := os.Mkdir(tapDir, 0o755); err != nil {
					t.Fatalf("setup: %v", err)
				}

				sentinelFile := filepath.Join(tapDir, "sentinel.txt")
				if err := os.WriteFile(sentinelFile, []byte("sentinel"), 0o644); err != nil {
					t.Fatalf("setup: %v", err)
				}

			},
			wantFiles: map[string]string{
				"lonlydwolf/homebrew-tools/.tapwright": wantMarker,
				"lonlydwolf/homebrew-tap/sentinel.txt": "sentinel",
			},
			wantDirs: []string{
				"lonlydwolf",
				"lonlydwolf/homebrew-tap",
				"lonlydwolf/homebrew-tools",
				"lonlydwolf/homebrew-tools/Casks",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if tc.setup != nil {
				tc.setup(t, root)
			}
			err := CreateTap(root, tc.user, tc.tap)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("CreateTap() error = %v, want %v", err, tc.wantErr)
			}

			for rel, want := range tc.wantFiles {
				path := filepath.Join(root, rel)
				if got, err := os.ReadFile(path); err != nil {
					t.Errorf("read %s: %v", path, err)
				} else if string(got) != want {
					t.Errorf("CreateTap() got = %q, want %v", string(got), want)
				}
			}

			for _, rel := range tc.wantDirs {
				path := filepath.Join(root, rel)
				info, err := os.Stat(path)
				switch {
				case errors.Is(err, fs.ErrNotExist):
					t.Errorf("CreateTap() error = %v, expected %v directory", err, path)
				case err != nil:
					t.Errorf("stat %s: %v", path, err)
				case !info.IsDir():
					t.Errorf("CreateTap() got = %q, want directory %v", info.Name(), path)
				}
			}

			for _, rel := range tc.wantAbsent {
				path := filepath.Join(root, rel)
				if info, err := os.Stat(path); err == nil {
					t.Errorf("CreateTap() got = %q, want directory %v to not exist", info.Name(), path)
				} else if !errors.Is(err, fs.ErrNotExist) {
					t.Errorf("stat %s: %v", path, err)
				}
			}

			allowed := map[string]bool{}
			for rel := range tc.wantFiles {
				allowed[rel] = true
			}
			for _, rel := range tc.wantDirs {
				allowed[rel] = true
			}

			err = fs.WalkDir(os.DirFS(root), ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if path != "." && !allowed[path] {
					t.Errorf("CreateTap() left unexpected entry %q in the tap root", path)
				}
				return nil
			})
			if err != nil {
				t.Fatalf("walk %s: %v", root, err)
			}
		})
	}
}
