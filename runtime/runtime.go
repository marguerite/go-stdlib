package runtime

import (
	"bufio"
	"os/exec"
	"runtime"
	"strings"

	cmd "github.com/marguerite/go-stdlib/exec"
	"github.com/marguerite/go-stdlib/ioutils"
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

// LinuxDistribution query /etc/os-release and return the distribution name
func LinuxDistribution() (distribution string) {
	scanner := bufio.NewScanner(ioutils.NewReaderFromFile("/etc/os-release"))
	pretty := "PRETTY_NAME=\""
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), pretty) {
			distribution = strings.ReplaceAll(strings.ReplaceAll(scanner.Text(), pretty, ""), "\"", "")
		}
	}
	return distribution
}
