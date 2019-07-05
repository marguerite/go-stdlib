package dirutils

import (
	"github.com/marguerite/util/slice"
	"reflect"
	"regexp"
	"testing"
)

func TestLsFile(t *testing.T) {
	if files, e := Ls("."); !reflect.DeepEqual(files, []string{"dirutils.go", "dirutils_test.go"}) || e != nil {
		t.Error("dirutils.Ls test failed.")
	} else {
		t.Log("dirutils.Ls test passed.")
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
	files, err := Glob("fakeDir", re, fn)
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
		regexp.MustCompile("opencc\\/.*\\.json"),
	}

	files, err := Glob("fakeDir", re, fn)
	if err != nil {
		t.Error("TestGlobRegexpGroup failed.")
	}

	b, err := slice.Contains(files, "zhung.dict.yaml")
	b1, err1 := slice.Contains(files, "opencc/s2t.json")
	if (b && err == nil) && (b1 && err1 == nil) {
		t.Log("TestGlobRegexpGroup passed.")
	} else {
		t.Error("TestGlobRegexpGroup failed.")
	}
}
