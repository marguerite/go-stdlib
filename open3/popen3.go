package open3

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Wait_thr the wait thread of process pid and exit code
type Wait_thr struct {
	Pid   int
	Value int
}

// Popen3 like ruby popen3 without block but func
// Note: you need to close stdin after use in fn
func Popen3(cmd *exec.Cmd, chdir string, fn func(stdin io.WriteCloser, stdout, stderr io.ReadCloser, wt Wait_thr) error, env ...string) (Wait_thr, error) {
	var wt Wait_thr
	// setup environment
	if len(env) > 0 {
		// check environment variables
		for _, v := range env {
			if strings.Contains(v, "=") {
				continue
			}
			return wt, errors.New("invalid environment variable")
		}
		cmd.Env = env
	}

	// if we switch to the working directory
	wd, err := os.Getwd()
	if err != nil {
		return wt, err
	}
	if len(chdir) > 0 {
		err := os.Chdir(chdir)
		if err != nil {
			return wt, err
		}
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return wt, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return wt, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdin.Close()
		return wt, err
	}

	err = cmd.Start()

	if err != nil {
		stdin.Close()
		if _, ok := err.(*exec.Error); ok {
			wt.Value = -1
		}

		if val, ok := err.(*exec.ExitError); ok {
			if stat, ok := val.Sys().(syscall.WaitStatus); ok {
				wt.Value = stat.ExitStatus()
			}
		}
		return wt, err
	}

	wt.Pid = cmd.Process.Pid

	err = fn(stdin, stdout, stderr, wt)

	if err != nil {
		stdin.Close()
		return wt, err
	}

	err = cmd.Wait()
	if err != nil {
		if _, ok := err.(*exec.Error); ok {
			wt.Value = -1
		}

		if val, ok := err.(*exec.ExitError); ok {
			if stat, ok := val.Sys().(syscall.WaitStatus); ok {
				wt.Value = stat.ExitStatus()
			}
		}
	}

	// switch back
	if len(chdir) > 0 {
		err = os.Chdir(wd)
		if err != nil {
			return wt, err
		}
	}

	return wt, nil
}
