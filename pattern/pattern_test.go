package pattern

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marguerite/go-stdlib/slice"
)

// **/*/?/{}/[]/\\
func TestDoubleAsteriskPattern(t *testing.T) {
	cwd, _ := os.Getwd()
	mainDir := filepath.Dir(cwd)
	patt := filepath.Join(mainDir, "**", "pattern.go")
	correct := filepath.Join(mainDir, "pattern", "pattern.go")
	result := Expand(patt)
	if len(result) > 0 {
		if ok, err := slice.Contains(result, correct); !ok || err != nil {
			t.Errorf("[pattern]: double aterisk test failed, expecting %s, got %s", correct, result)
		}
	} else {
		t.Errorf("[pattern]: double asterisk test failed, expecting %s, got empty", correct)
	}
}

func TestSingleAsteriskPattern(t *testing.T) {
	cwd, _ := os.Getwd()
	mainDir := filepath.Dir(cwd)
	t.Logf("main directory is %s", mainDir)
	patt := filepath.Join(mainDir, "pattern", "pattern_un*.go")
	correct := filepath.Join(mainDir, "pattern", "pattern_unix.go")
	result := Expand(patt)
	if len(result) > 0 {
		if ok, err := slice.Contains(result, correct); !ok || err != nil {
			t.Errorf("[pattern]: single aterisk test failed, expecting %s, got %s", correct, result)
		}
	} else {
		t.Errorf("[pattern]: single asterisk test failed, expecting %s, got empty", correct)
	}
}

func TestQuestionMarkPattern(t *testing.T) {
	cwd, _ := os.Getwd()
	mainDir := filepath.Dir(cwd)
	t.Logf("main directory is %s", mainDir)
	patt := filepath.Join(mainDir, "pattern", "patter?.go")
	correct := filepath.Join(mainDir, "pattern", "pattern.go")
	result := Expand(patt)
	if len(result) > 0 {
		if ok, err := slice.Contains(result, correct); !ok || err != nil {
			t.Errorf("[pattern]: question mark test failed, expecting %s, got %s", correct, result)
		}
	} else {
		t.Errorf("[pattern]: question mark test failed, expecting %s, got empty", correct)
	}
}

func TestBracketPattern(t *testing.T) {
	cwd, _ := os.Getwd()
	mainDir := filepath.Dir(cwd)
	t.Logf("main directory is %s", mainDir)
	patt := filepath.Join(mainDir, "pattern", "patt[e,f]rn.go")
	correct := filepath.Join(mainDir, "pattern", "pattern.go")
	result := Expand(patt)
	if len(result) > 0 {
		if ok, err := slice.Contains(result, correct); !ok || err != nil {
			t.Errorf("[pattern]: bracket test failed, expecting %s, got %s", correct, result)
		}
	} else {
		t.Errorf("[pattern]: bracket test failed, expecting %s, got empty", correct)
	}
}

func TestBracePattern(t *testing.T) {
	cwd, _ := os.Getwd()
	mainDir := filepath.Dir(cwd)
	t.Logf("main directory is %s", mainDir)
	patt := filepath.Join(mainDir, "pattern", "patter{m,n}.go")
	correct := filepath.Join(mainDir, "pattern", "pattern.go")
	result := Expand(patt)
	if len(result) > 0 {
		if ok, err := slice.Contains(result, correct); !ok || err != nil {
			t.Errorf("[pattern]: brace test failed, expecting %s, got %s", correct, result)
		}
	} else {
		t.Errorf("[pattern]: brace test failed, expecting %s, got empty", correct)
	}
}
