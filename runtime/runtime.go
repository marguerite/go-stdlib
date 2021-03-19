package runtime

import (
	"bytes"
	"runtime"
	"strings"

	"github.com/marguerite/go-gnulib/login"
	"github.com/marguerite/go-stdlib/internal"
)

// Is64Bit if the operation system is 64bit system
func Is64Bit() bool {
	arch := runtime.GOARCH
	if strings.Contains(arch, "64") || strings.Contains(arch, "390") || strings.Contains(arch, "wasm") {
		return true
	}
	return false
}

// LogName the current login user's name
func LogName() string {
	name, err := login.GetLogin()
	if err != nil {
		panic(err)
	}
	// name is fixed-width, may has many tailing null bytes, thus os.Open may fail
	return internal.Bytes2str(bytes.Trim(internal.Str2bytes(name), "\x00"))
}
