package archiver

import (
	"errors"
	"regexp"

	"path/filepath"
)

func checkPath(p string) error {
	if len(p) == 0 {
		return errors.New("empty path")
	}
	if filepath.IsAbs(p) {
		return errors.New("non-relative path: " + p)
	}
	return nil
}

var isArchiveRegex = regexp.MustCompile(`(\.tar|\.tar\.gz|\.zip)$`)

func ValidExtension(path string) bool {
	return isArchiveRegex.MatchString(path)
}

func Extension(path string) string {
	m := isArchiveRegex.FindStringSubmatch(path)
	if len(m) == 0 {
		return ""
	}
	return m[1]
}
