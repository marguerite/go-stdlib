package fileutils

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func TestCopy(t *testing.T) {
	fn := func(s, d, o string) ([]string, error) {
		return []string{s + ".new"}, nil
	}
	res, err := copy("fileutils.go", "fileutils.go.new", []*regexp.Regexp{}, fn)
	if err != nil {
		t.Error("fileutils.Copy test failed")
	}
	if reflect.DeepEqual(res, []string{"fileutils.go.new"}) {
		t.Log("fileutils.Copy test succeed")
	} else {
		t.Error("fileutils.Copy test failed")
	}
}
