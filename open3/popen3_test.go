package open3

import (
	"bytes"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPopen3(t *testing.T) {
	cmd := exec.Command("/usr/bin/ls")
	stdoutbuf := bytes.NewBuffer([]byte{})
	stderrbuf := bytes.NewBuffer([]byte{})
	chdir, _ := filepath.Abs("../dir")
	wt, err := Popen3(cmd, chdir, func(stdin io.WriteCloser, stdout, stderr io.ReadCloser, wt Wait_thr) error {
		stdin.Close()
		stdoutbuf.ReadFrom(stdout)
		stderrbuf.ReadFrom(stderr)
		return nil
	})

	if wt.Value != 0 {
		t.Errorf("popen3 failed, expected 0, got %d", wt.Value)
	}

	if err != nil {
		t.Errorf("popen3 failed, expected nil, got %s", err)
	}

	if !strings.Contains(stdoutbuf.String(), "dir.go") {
		t.Errorf("popen3 failed, expected dir.go dir_test.go, got %s", stdoutbuf.String())
	}

	if stderrbuf.String() != "" {
		t.Errorf("popen3 failed, expected nil, got %s", stderrbuf.String())
	}
}

func TestPopen3Pipe(t *testing.T) {
	cmd := exec.Command("/usr/bin/ls")
	cmd1 := exec.Command("/usr/bin/grep", "test")

	stdoutbuf := bytes.NewBuffer([]byte{})
	stdout1buf := bytes.NewBuffer([]byte{})

	Popen3(cmd, "", func(stdin io.WriteCloser, stdout, stderr io.ReadCloser, wt Wait_thr) error {
		stdin.Close()
		stdoutbuf.ReadFrom(stdout)
		return nil
	})

	Popen3(cmd1, "", func(stdin io.WriteCloser, stdout, stderr io.ReadCloser, wt Wait_thr) error {
		_, err := io.Copy(stdin, bytes.NewReader(stdoutbuf.Bytes()))
		stdin.Close()
		if err != nil {
			return err
		}
		stdout1buf.ReadFrom(stdout)
		return nil
	})

	if stdout1buf.String() != "popen3_test.go\n" {
		t.Errorf("popen3 pipe failed, expected popen3_test.go, got %s", stdout1buf.String())
	}
}
