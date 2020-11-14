package runtime

import (
	"os/exec"
	"runtime"
	"strings"

	cmd "github.com/marguerite/go-stdlib/exec"
)

// Is64Bit if the operation system is 64bit system
func Is64Bit() bool {
	arch := runtime.GOARCH
	if strings.Contains(arch, "64") || strings.Contains(arch, "390") || strings.Contains(arch, "wasm") {
		return true
	}
	return false
}

// LogName the current log in user's name
func LogName() string {
	out, err := exec.Command("/usr/bin/logname").Output()
	if err != nil {
		env, err := cmd.Env("LOGNAME")
		if err != nil {
			return ""
		}
		return env
	}
	return string(out)
}
