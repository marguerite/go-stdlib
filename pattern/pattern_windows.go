package pattern

import (
	"os"
	"path/filepath"
)

const PATH_SEPARATOR = "\\"

var (
	root = filepath.VolumeName(os.GetEnv("SYSTEMROOT")) + PATHSEPARATOR
)
