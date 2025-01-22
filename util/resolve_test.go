package util

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestResolveFilename(t *testing.T) {
	// Create a mock filesystem with test files
	mockFS := fstest.MapFS{
		"src/styles/variables.scss":                                     &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/_partial.scss":                                      &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/components/_index.scss":                             &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/mixins.sass":                                        &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/theme.css":                                          &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/utils/_helpers.scss":                                &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/nested/main/_index.scss":                            &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/nested/_main.scss":                                  &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/config.import.scss":                                 &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/base.import.scss":                                   &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/_more.import.scss":                                  &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/lib.import.sass":                                    &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/plain/index.scss":                                   &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/with.dots/index.scss":                               &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/with.dots.scss/":                                    &fstest.MapFile{Data: []byte(""), Mode: 0644 | fs.ModeDir},
		"src/styles/normal_before_index/other/index.scss":               &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/normal_before_index/other.scss":                     &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/index_import_before_normal/other/index.scss":        &fstest.MapFile{Data: []byte(""), Mode: 0666},
		"src/styles/index_import_before_normal/other/index.import.scss": &fstest.MapFile{Data: []byte(""), Mode: 0666},
	}

	tests := []struct {
		name        string
		source      string
		importPath  string
		wantPath    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "Basic SCSS file",
			source:     "src/styles/main.scss",
			importPath: "variables",
			wantPath:   "src/styles/variables.scss",
		},
		{
			name:       "Partial file with underscore",
			source:     "src/styles/main.scss",
			importPath: "partial",
			wantPath:   "src/styles/_partial.scss",
		},
		{
			name:       "Directory with index",
			source:     "src/styles/main.scss",
			importPath: "components",
			wantPath:   "src/styles/components/_index.scss",
		},
		{
			name:       "SASS extension",
			source:     "src/styles/main.scss",
			importPath: "mixins",
			wantPath:   "src/styles/mixins.sass",
		},
		{
			name:       "CSS extension",
			source:     "src/styles/main.scss",
			importPath: "theme",
			wantPath:   "src/styles/theme.css",
		},
		{
			name:       "Nested partial",
			source:     "src/styles/main.scss",
			importPath: "utils/helpers",
			wantPath:   "src/styles/utils/_helpers.scss",
		},
		{
			name:       "Nested index",
			source:     "src/styles/main.scss",
			importPath: "nested/main",
			wantPath:   "src/styles/nested/main/_index.scss",
		},
		{
			name:       "Explicit extension",
			source:     "src/styles/main.scss",
			importPath: "variables.scss",
			wantPath:   "src/styles/variables.scss",
		},
		{
			name:       "Explicit sass extension",
			source:     "src/styles/main.scss",
			importPath: "mixins.sass",
			wantPath:   "src/styles/mixins.sass",
		},
		{
			name:       "Import SCSS extension",
			source:     "src/styles/main.scss",
			importPath: "config",
			wantPath:   "src/styles/config.import.scss",
		},
		{
			name:       "Import SCSS extension - check underscore",
			source:     "src/styles/main.scss",
			importPath: "more",
			wantPath:   "src/styles/_more.import.scss",
		},
		{
			name:       "Explicit import SCSS extension",
			source:     "src/styles/main.scss",
			importPath: "base.import.scss",
			wantPath:   "src/styles/base.import.scss",
		},
		{
			name:       "Import SASS extension",
			source:     "src/styles/main.scss",
			importPath: "lib",
			wantPath:   "src/styles/lib.import.sass",
		},
		{
			name:       "Directory with plain index",
			source:     "src/styles/main.scss",
			importPath: "plain",
			wantPath:   "src/styles/plain/index.scss",
		},
		{
			name:       "Directory with dots in name",
			source:     "src/styles/main.scss",
			importPath: "with.dots",
			wantPath:   "src/styles/with.dots/index.scss",
		},
		{
			name:        "Directory with extension in name",
			source:      "src/styles/main.scss",
			importPath:  "with.dots.scss",
			wantErr:     true,
			errContains: "cannot import directory with extension",
		},
		{
			name:        "Non-existent file",
			source:      "src/styles/main.scss",
			importPath:  "nonexistent",
			wantErr:     true,
			errContains: "no such file or directory",
		},
		{
			name:       "Precedence of .scss over folder",
			source:     "src/styles/normal_before_index/main.scss",
			importPath: "other",
			wantPath:   "src/styles/normal_before_index/other.scss",
		},
		{
			name:       "Precedence of .import.scss over .scss",
			source:     "src/styles/index_import_before_normal/main.scss",
			importPath: "other",
			wantPath:   "src/styles/index_import_before_normal/other/index.import.scss",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveFilename(tt.source, tt.importPath, mockFS)
			if tt.wantErr {
				if err == nil {
					t.Error("ResolveFilename() error = nil, wantErr = true")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ResolveFilename() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ResolveFilename() error = %v", err)
				return
			}
			if filepath.ToSlash(got) != tt.wantPath {
				t.Errorf("ResolveFilename() = %v, want %v", got, tt.wantPath)
			}
		})
	}
}
