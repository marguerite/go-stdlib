package dirutils

import (
	"reflect"
	"testing"
)

func TestLs(t *testing.T) {
	if dirs, e := Ls(".", "file"); !reflect.DeepEqual(dirs, []string{"dir.go", "dir_test.go"}) || e != nil {
		t.Error("dir.Ls test failed.")
	} else {
		t.Log("dir.Ls test passed.")
	}
}

func TestLsf(t *testing.T) {
	if dirs, e := Ls(".", "file"); !reflect.DeepEqual(dirs, []string{"dir.go", "dir_test.go"}) ||
		e != nil {
		t.Error("dir.Lsf test failed.")
	} else {
		t.Log("dir.Lsf test passed.")
	}
}
