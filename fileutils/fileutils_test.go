package fileutils

import (
	"testing"
)

func TestCopy(t *testing.T) {
	fn := func(s, d, o string) error { return nil }
	err := copy("fileutils.go", "fileutils.go.new", fn)
	if err != nil {
		t.Error("fileutils.Copy test failed")
	}
}
