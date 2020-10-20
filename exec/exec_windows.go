// +build windows

package exec

import (
	"os"
)

func isExecutable(f os.FileInfo) bool {
	return true
}
