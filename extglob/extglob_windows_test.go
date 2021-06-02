// +build windows

package extglob

import (
	"bytes"
	"os"
	"testing"
)

var (
	extglobPattern = []string{"C:\\*(m)arguerite",
		"C:\\?(m)arguerite",
		"C:\\mar@(g|h)uerite",
		"C:\\mar+(g)uerite",
		"C:\\!(a)arguerite",
	}

	shellPattern = []string{
		"**\\marguerite",
		"C:\\marguer?te",
		"C:\\[^a]arguerite",
		"C:\\[h-m]arguerite",
		"C:\\{marguerite,allen}",
	}
)

func TestEscaped(t *testing.T) {
	_, err := escaped('\\')
	if err == nil {
		t.Error("test escaped failed, it should panic an error, but we got nil")
	}
}

func TestIsExtGlobPattern(t *testing.T) {
	for _, v := range extglobPattern {
		bol := isExtGlobPattern([]byte(v))
		if !bol {
			t.Errorf("isExtGlobPattern %s failed, expected true, got %t", v, bol)
		}
	}
}

func TestIsPlainShellPattern(t *testing.T) {
	for _, v := range shellPattern {
		bol := isPlainShellPattern([]byte(v))
		if !bol {
			t.Errorf("isExtGlobPattern %s failed, expected true, got %t", v, bol)
		}
	}
}

func list1(p string, t, t1 bool) ([]string, error) {
	if t {
		return []string{"C:\\"}, nil
	}
	return []string{"C:\\marguerite"}, nil
}

func valid1(p string) bool {
	return true
}

func TestExpand(t *testing.T) {
	for _, pattern := range [][]string{extglobPattern, shellPattern} {
		for _, v := range pattern {
			results, err := expand([]byte(v), true, true, list1, valid1)
			if err != nil {
				t.Errorf("expand %s failed, expected nil error, got %s", v, err)
			}
			if len(results) > 1 || len(results) == 0 {
				t.Errorf("expand %s failed, expected len 1, got %d %v", v, len(results), results)
			} else {
				if results[0] != "C:\\marguerite" {
					t.Errorf("expand %s failed, expected C:\\marguerite, got %s", v, results[0])
				}
			}
		}
	}
}

func TestShellmatchStar(t *testing.T) {
	files := []string{"C:\\marguerite", "C:\\zhou", "C:\\erite"}
	buf := bytes.NewBufferString("*erite")
	shellmatch(&files, buf, 0)

	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("shellmatch failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestShellmatchNormal(t *testing.T) {
	files := []string{"C:\\marguerite", "C:\\zhou", "C:\\wenxuetian"}
	buf := bytes.NewBufferString("xuetian")
	shellmatch(&files, buf, 3)

	if len(files) == 0 || files[0] != "C:\\wenxuetian" {
		t.Errorf("shellmatch failed, expected []string{\"C:\\wenxuetian\"}, got %v", files)
	}
}

func TestShellmatchQuestion(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("?uerite")
	shellmatch(&files, buf, 3)
	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("shellmatch failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestShellmatchBracketRange(t *testing.T) {
	files := []string{"C:\\marguerite", "C:\\zhou", "C:\\wenxuetian"}
	buf := bytes.NewBufferString("[^a]arguer[^m]te")
	shellmatch(&files, buf, 0)

	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("shellmatch failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestShellmatchCurlyRange(t *testing.T) {
	files := []string{"C:\\marguerite", "C:\\marzhouite", "C:\\marwuefrite"}
	buf := bytes.NewBufferString("mar{gue,wuef}rite")
	shellmatch(&files, buf, 0)

	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("shellmatch failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestShellmatchOr(t *testing.T) {
	files := []string{"C:\\marguerite", "C:\\wenxuetian"}
	bufs := []*bytes.Buffer{bytes.NewBufferString("marguerite"), bytes.NewBufferString("wenxuetian")}
	result := shellmatchor(files, bufs, 0)
	if len(result) != 2 {
		t.Errorf("shellmatchor failed, expected %v, got %v", files, result)
	}
}

func TestParseBracketsWithHyphen(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar[-gh-]uerite")
	shellmatch(&files, buf, 0)
	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("parseBracketRange failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestParseBracketsReverse(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar[!az-]uerite")
	shellmatch(&files, buf, 0)
	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("parseBracketRange  failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestParseBracketsWithHyphen1(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar[a-z]uerite")
	shellmatch(&files, buf, 0)
	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("parseBracketRange  failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestParseBracketsWithHyphenAndLocaleC(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar[a-z]uerite")
	os.Setenv("LC_ALL", "C")
	shellmatch(&files, buf, 0)
	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("parseBracketRange  failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestParseBracketsWithRightBracket(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar[]gh]uerite")
	shellmatch(&files, buf, 0)
	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("parseBracketRange  failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestParseBracketsWithEqual(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar[=g=]uerite")
	shellmatch(&files, buf, 0)
	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("parseBracketRange  failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestParseBracketsWithClass(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar[:alpha:]uerite")
	shellmatch(&files, buf, 0)
	if len(files) == 0 || files[0] != "C:\\marguerite" {
		t.Errorf("parseBracketRange  failed, expected []string{\"C:\\marguerite\"}, got %v", files)
	}
}

func TestMatchWithQuestion(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar?(g|h)u?(y|z)erite")
	result := match(files, buf, true, 0)
	if len(result) == 0 || result[0] != files[0] {
		t.Errorf("match failed, expected []string{\"C:\\marguerite\"}, got %v", result)
	}
}

func TestMatchWithAt(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar@(g|h)uerite")
	result := match(files, buf, true, 0)
	if len(result) == 0 || result[0] != files[0] {
		t.Errorf("match failed, expected []string{\"C:\\marguerite\"}, got %v", result)
	}
}

func TestMatchWithExclamation(t *testing.T) {
	files := []string{"C:\\marguerite"}
	buf := bytes.NewBufferString("mar!(y|z)uerite")
	result := match(files, buf, true, 0)
	if len(result) == 0 || result[0] != files[0] {
		t.Errorf("match failed, expected []string{\"C:\\marguerite\"}, got %v", result)
	}
}

func TestMatchWithStar(t *testing.T) {
	files := []string{"C:\\marggguerite"}
	buf := bytes.NewBufferString("mar*(g|h)uerite")
	result := match(files, buf, true, 0)
	if len(result) == 0 || result[0] != files[0] {
		t.Errorf("match failed, expected []string{\"C:\\marggguerite\"}, got %v", result)
	}
}

func TestMatchWithAdd(t *testing.T) {
	files := []string{"C:\\marggguerite"}
	buf := bytes.NewBufferString("mar+(g|h)uerite")
	result := match(files, buf, true, 0)
	if len(result) == 0 || result[0] != files[0] {
		t.Errorf("match failed, expected []string{\"C:\\marggguerite\"}, got %v", result)
	}
}
