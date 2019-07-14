package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// Environ safely get an environment variable
func Environ(env string) (string, error) {
	val, ok := os.LookupEnv(env)
	if !ok {
		return "", fmt.Errorf("%s not set.", env)
	}
	if len(val) == 0 {
		return val, fmt.Errorf("%s is empty.", env)
	}
	return val, nil
}

func search(cmd string) (string, error) {
	// the cmd definitely exists
	f, _ := os.Stat(cmd)
	if f.IsDir() {
		return "", fmt.Errorf("%s is a directory", cmd)
	}
	if strings.Contains(f.Mode().Perm().String(), "-rwxr-") {
		return cmd, nil
	}
	return "",f
}

// Search if an executable exists
func Search(cmd string) (string, error) {
  f, err := os.Stat(cmd)
	if err != nil {
		if os.IsNotExist(err) {
      // add $PATH and try again
		} else {
			return "", fmt.Errof("Another unhandled non-IsNotExist PathError occurs %s", err.Error())
		}
	}

	// whether to add path.
	if filepath.Dir(cmd) == "." {
		if _, err := os.Stat(cmd); os.IsNotExist(err) {
			pathEnv, err := Environ("PATH")
			if err != nil {
				return "", errors.New("System $PATH not set or empty.")
			}
			ok := false
			command := ""
			for _, v := range strings.Split(pathEnv, ":") {
				c := filepath.Join(v, cmd)
				if _, err := os.Stat(c); !os.IsNotExist(err) {
					command = c
					ok = true
					break
				}
			}
			if ok {
				cmd = command
			} else {
				return "", fmt.Errorf("Can not find executable %s", cmd)
			}
		}
	} else {
		if _, err := os.Stat(cmd); os.IsNotExist(err) {
			return "", fmt.Errorf("Can not find executable %s", cmd)
		}
	}

	if f, _ := os.Stat(cmd); strings.Contains(f.Mode().String(), "-rwxr-") {
		return cmd, nil
	}
	return cmd, fmt.Errorf("%s exists but not executable.", cmd)
}

// Run run command with options, returns output, ExitStatus and error
func Run(cmd string, opts ...string) (string, int, error) {
	c, err := Search(cmd)
	if err != nil {
		return "", -1, err
	}
	out, err := exec.Command(c, opts...).Output()

	fmt.Printf("Executing: %s %s\n", c, strings.Join(opts, " "))

	if err != nil {
		if msg, ok := err.(*exec.Error); ok {
			return string(out), -1, err
		}

		if msg, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := msg.Sys().(syscall.WaitStatus); ok {
				return string(out), waitStatus.ExitStatus(), err
			}
		}

		return string(out), -1, err
	}

	return string(out), 0, nil
}
