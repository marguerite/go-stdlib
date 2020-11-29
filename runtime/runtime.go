package runtime

import (
	"os/exec"
	"runtime"
	"strings"

  //"github.com/ericlagergren/go-gnulib/login"
  "github.com/marguerite/go-gnulib/login"
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
  return name
}
