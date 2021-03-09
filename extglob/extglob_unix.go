// +build linux freebsd openbsd netbsd darwin

package extglob

import (
	"bytes"

	internal "github.com/marguerite/go-stdlib/internal"
)

const PATH_SEPARATOR = "/"

var (
	root = "/"
)

func writeBytes(buf *bytes.Buffer, b []byte, sep bool) {
	// unix system has path like "/home/marguerite"
	if sep {
		(*buf).Write(internal.Str2bytes(PATH_SEPARATOR))
	}
	(*buf).Write(b)
}
