// +build windows
package extglob

import (
	"bytes"
	"os"
	"path/filepath"

	internal "github.com/marguerite/go-stdlib/internal"
)

const PATH_SEPARATOR = "\\"

var (
	root = filepath.VolumeName(os.GetEnv("SYSTEMROOT")) + PATHSEPARATOR
)

func writeBytes(buf *bytes.Buffer, b []byte, sep bool) {
	// windows system has path like "C:\\Program Files(x64)\Tencent QQ"
	// the path separator will not occur in the first place of path at all
	(*buf).Write(b)
	(*buf).Write(internal.Str2bytes(PATH_SEPARATOR))
}
