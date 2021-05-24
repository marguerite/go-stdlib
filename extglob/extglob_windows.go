// +build windows
package extglob

import (
	"os"
	"path/filepath"
)

const PATH_SEPARATOR = '\\'

var (
	root = filepath.VolumeName(os.GetEnv("SYSTEMROOT")) + PATHSEPARATOR
)
