package dirutils

import (
	"github.com/marguerite/util/slice"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

func TestLsFile(t *testing.T) {
	cwd, _ := os.Getwd()
	if files, e := Ls("."); !reflect.DeepEqual(files, []string{filepath.Join(cwd, "dirutils.go"), filepath.Join(cwd, "dirutils_test.go")}) || e != nil {
		t.Error("dirutils.Ls test failed.")
	} else {
		t.Log("dirutils.Ls test passed.")
	}
}

func TestWildcard(t *testing.T) {
	s := "/usr/lib/libutil.{a,    la}"
	r := []*regexp.Regexp{regexp.MustCompile("/usr/lib/libutil.a"), regexp.MustCompile("/usr/lib/libutil.la")}
	rs, rr := ParseWildcard(s)
	if rs == "/usr/lib" && reflect.DeepEqual(r, rr) {
		t.Log("ParseWildcard test succeed.")
	} else {
		t.Error("ParseWildcard test failed.")
	}
}

func TestWildcard1(t *testing.T) {
	s := "/usr/lib/libutil*"
	r := []*regexp.Regexp{regexp.MustCompile("/usr/lib/libutil.*")}
	rs, rr := ParseWildcard(s)
	if rs == "/usr/lib" && reflect.DeepEqual(r, rr) {
		t.Log("ParseWildcard test succeed.")
	} else {
		t.Error("ParseWildcard test failed.")
	}
}

func TestWildcard2(t *testing.T) {
	s := "/usr/lib/libutil.so"
	r := []*regexp.Regexp{}
	rs, rr := ParseWildcard(s)
	if rs == s && reflect.DeepEqual(r, rr) {
		t.Log("ParseWildcard test succeed.")
	} else {
		t.Error("ParseWildcard test failed.")
	}
}

func TestGlobRegexp(t *testing.T) {
	fn := func(dir string, kind ...string) ([]string, error) {
		return []string{"zhung.schema.yaml",
			"zhung_42.custom.yaml",
			"zhung_10.recipe.yaml",
			"recipe.yaml"}, nil
	}

	re := regexp.MustCompile(".*\\.yaml")
	re1 := regexp.MustCompile("(custom|recipe)\\.yaml")
	files, err := glob("fakeDir", re, fn, re1)
	if err != nil {
		t.Error("TestGlobRegexp failed.")
	}

	if b, err := slice.Contains(files, "zhung.schema.yaml"); b && err == nil {
		t.Log("TestGlobRegexp passed.")
	} else {
		t.Error("TestGlobRegexp failed.")
	}
}

func TestGlobRegexpGroup(t *testing.T) {
	fn := func(dir string, kind ...string) ([]string, error) {
		return []string{"zhung.dict.yaml",
			"opencc/s2t.json"}, nil
	}

	re := []*regexp.Regexp{regexp.MustCompile(".*\\.yaml"),
		regexp.MustCompile("opencc\\/.*"),
	}
	re1 := []*regexp.Regexp{regexp.MustCompile("(custom|recipe)\\.yaml"),
		regexp.MustCompile("\\.(json|txt|ocd)")}

	files, err := glob("fakeDir", re, fn, re1)
	if err != nil {
		t.Error("TestGlobRegexpGroup failed.")
	}

	b, err := slice.Contains(files, "zhung.dict.yaml")
	b1, err1 := slice.Contains(files, "opencc/s2t.json")
	if (b && err == nil) && (!b1 && err1 == nil) {
		t.Log("TestGlobRegexpGroup passed.")
	} else {
		t.Error("TestGlobRegexpGroup failed.")
	}
}
