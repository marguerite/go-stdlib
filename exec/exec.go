package exec

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

var (
	ErrIsDirectory = errors.New("Is a directory")
)

// Env get an environment variable
func Env(env string) (val string, err error) {
	val, ok := os.LookupEnv(env)
	if !ok {
		return val, os.ErrNotExist
	}
	if len(val) == 0 {
		return val, os.ErrNotExist
	}
	return val, nil
}

// Search if an executable exists
func Search(cmd string) (val string, err error) {
	f, err := os.Stat(cmd)
	if err != nil {
		if os.IsNotExist(err) {
			// look in $Path
			val, err = exec.LookPath(cmd)
			if err != nil {
				return val, os.ErrNotExist
			}
			return val, nil
		}
		return val, err
	}
	if f.IsDir() {
		return val, ErrIsDirectory
	}
	if isExecutable(f) {
		return cmd, nil
	}
	return val, os.ErrPermission
}

// Exec3 run command with options, returns stdout, ExitStatus and error
func Exec3(cmd string, options ...string) (out []byte, exit int, err error) {
	out, err = exec.Command(cmd, options...).Output()

	if err == nil {
		return out, 0, err
	}

	if _, ok := err.(*exec.Error); ok {
		return out, -1, os.ErrNotExist
	}

	if val, ok := err.(*exec.ExitError); ok {
		if stat, ok := val.Sys().(syscall.WaitStatus); ok {
			return out, stat.ExitStatus(), err
		}
	}

	return out, -1, err
}
