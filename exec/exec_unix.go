// +build aix darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package exec

import (
	"os"
	"strings"
)

func isExecutable(f os.FileInfo) bool {
	if strings.Contains(f.Mode().String(), "-rwxr-") {
		return true
	}
	return false
}
