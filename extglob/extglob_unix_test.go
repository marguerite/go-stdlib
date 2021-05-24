// +build linux freebsd openbsd netbsd darwin

package extglob

import (
	"testing"
)

func TestEscaped(t *testing.T) {
	ok, err := escaped('\\')
	if !ok || err != nil {
		t.Error("test escaped failed, expected true, got false")
	}
}
