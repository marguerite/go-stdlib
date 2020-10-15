package dir

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/marguerite/util/slice"
)

func TestLs(t *testing.T) {
	cwd, _ := os.Getwd()
	correct := []string{filepath.Join(cwd, "dir_test.go"), filepath.Join(cwd, "dir.go")}
	if files, err := Ls(cwd, true, true); !reflect.DeepEqual(files, correct) || err != nil {
		t.Errorf("[dir]Ls test failed, expecting %s, got %s, err %s", correct, files, err)
	}
}

func TestGlobString(t *testing.T) {
	cwd, _ := os.Getwd()
	patt := filepath.Join(filepath.Dir(cwd), "**", "dir*.go")
	t.Logf("[dir]Glob pattern is %s", patt)
	correct := filepath.Join(filepath.Dir(cwd), "dir", "dir_test.go")
	result, err := Glob(patt)
	if err != nil {
		t.Errorf("[di]: Glob test failed with %s", err.Error())
	}
	if len(result) > 0 {
		if ok, err := slice.Contains(result, correct); !ok || err != nil {
			t.Errorf("[dir]: Glob test failed, expecting %s, got %s", correct, result)
		}
	} else {
		t.Errorf("[dir]: Glob test failed, expecting %s, got empty", correct)
	}
}

func TestGlobRegex(t *testing.T) {
	cwd, _ := os.Getwd()
	re := regexp.MustCompile(`dir.*\.go`)
	correct := filepath.Join(cwd, "dir_test.go")
	result, err := Glob(re, cwd)
	if err != nil {
		t.Errorf("[di]: Glob test failed with %s", err.Error())
	}
	if len(result) > 0 {
		if ok, err := slice.Contains(result, correct); !ok || err != nil {
			t.Errorf("[dir]: Glob test failed, expecting %s, got %s", correct, result)
		}
	} else {
		t.Errorf("[dir]: Glob test failed, expecting %s, got empty", correct)
	}
}

func TestGlobRegexWithExclusion(t *testing.T) {
	cwd, _ := os.Getwd()
	re := regexp.MustCompile(`dir.*\.go`)
	re1 := regexp.MustCompile(`dir_test\.go`)
	correct := filepath.Join(cwd, "dir.go")
	result, err := Glob(re, cwd, re1)
	if err != nil {
		t.Errorf("[di]: Glob test failed with %s", err.Error())
	}
	if len(result) > 0 {
		if ok, err := slice.Contains(result, correct); !ok || err != nil {
			t.Errorf("[dir]: Glob test failed, expecting %s, got %s", correct, result)
		}
	} else {
		t.Errorf("[dir]: Glob test failed, expecting %s, got empty", correct)
	}
}
