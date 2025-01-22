package util

import (
	"fmt"
	"io/fs"
	"path"
	"strings"
)

// ResolveFilename resolves a module name to an actual file path following SCSS import rules.
// source is the path of the file containing the import statement
// name is the module name to resolve
// fsys is the filesystem to use for file operations
func ResolveFilename(source string, name string, fsys fs.FS) (string, error) {
	base := path.Dir(source)

	// List of possible extensions and prefixes
	importExtensions := []string{".import.scss", ".import.sass"}
	extensions := []string{".scss", ".sass", ".css"}
	allExtensions := append(importExtensions, extensions...)
	prefixes := []string{"", "_"}

	// Helper function to check if file exists
	fileExists := func(path string) bool {
		fi, err := fs.Stat(fsys, path)
		return err == nil && !fi.IsDir()
	}

	// Helper function to check if directory exists
	dirExists := func(path string) bool {
		fi, err := fs.Stat(fsys, path)
		return err == nil && fi.IsDir()
	}

	// Helper function to check if path has a supported extension
	hasExtension := func(p string) bool {
		ext := path.Ext(p)
		for _, supportedExt := range allExtensions {
			if ext == supportedExt {
				return true
			}
		}
		return false
	}

	// If name already has a supported extension, check that exact file first
	if hasExtension(name) {
		fullPath := path.Join(base, name)
		if fileExists(fullPath) {
			return fullPath, nil
		}
		// If it's a directory with extension, return error
		if dirExists(fullPath) {
			return "", fmt.Errorf("cannot import directory with extension '%s'", name)
		}
		// If exact file not found with explicit extension,
		// return error immediately as we shouldn't try other variants
		return "", fmt.Errorf("could not find file '%s' relative to '%s'", name, source)
	}

	// First try all possible combinations for the main file
	// Check import extensions first
	for _, ext := range importExtensions {
		for _, prefix := range prefixes {
			fullPath := path.Join(base, path.Dir(name), prefix+path.Base(name)+ext)
			if fileExists(fullPath) {
				return fullPath, nil
			}
		}
	}

	// If we have a nested path, check for index files first
	if strings.Contains(name, "/") {
		basePath := path.Join(base, path.Dir(name))
		fileName := path.Base(name)

		nestedDirPath := path.Join(basePath, fileName)
		if dirExists(nestedDirPath) {
			// First check for index.import.* files
			for _, prefix := range []string{"_index", "index"} {
				for _, ext := range importExtensions {
					indexPath := path.Join(nestedDirPath, prefix+ext)
					if fileExists(indexPath) {
						return indexPath, nil
					}
				}
			}
			// Then check for regular index files
			for _, prefix := range []string{"_index", "index"} {
				for _, ext := range extensions[:2] { // Only check .scss and .sass for index files
					indexPath := path.Join(nestedDirPath, prefix+ext)
					if fileExists(indexPath) {
						return indexPath, nil
					}
				}
			}
		}
	}

	// Then check regular extensions
	for _, ext := range extensions {
		for _, prefix := range prefixes {
			fullPath := path.Join(base, path.Dir(name), prefix+path.Base(name)+ext)
			if fileExists(fullPath) {
				return fullPath, nil
			}
		}
	}

	// Finally check for regular index files in directories
	dirPath := path.Join(base, name)
	if dirExists(dirPath) {
		// First check for index.import.* files
		for _, prefix := range []string{"_index", "index"} {
			for _, ext := range importExtensions {
				indexPath := path.Join(dirPath, prefix+ext)
				if fileExists(indexPath) {
					return indexPath, nil
				}
			}
		}
		// Then check for regular index files
		for _, prefix := range []string{"_index", "index"} {
			for _, ext := range extensions[:2] { // Only check .scss and .sass for index files
				indexPath := path.Join(dirPath, prefix+ext)
				if fileExists(indexPath) {
					return indexPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not resolve import path '%s' relative to '%s': no such file or directory", name, source)
}
