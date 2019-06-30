package dirutils

import (
	"reflect"
	"testing"
)

func TestLsFile(t *testing.T) {
	if files, e := Ls(".", "file"); !reflect.DeepEqual(files, []string{"dirutils.go", "dirutils_test.go"}) || e != nil {
		t.Error("dirutils.Ls test failed.")
	} else {
		t.Log("dirutils.Ls test passed.")
	}
}
