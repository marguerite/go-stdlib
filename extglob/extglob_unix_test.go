// +build linux freebsd openbsd netbsd darwin

package extglob

import (
	"testing"

	internal "github.com/marguerite/go-stdlib/internal"
)

func TestNewPattern(t *testing.T) {
	b := internal.Str2bytes("/home/[h-m]arguerite/Packages")
	pattern, err := NewPattern(b, false)
	if err != nil {
		t.Errorf("NewPattern failed, expected nil error, got %s", err)
	}
	if pattern.Prefix != "/home" {
		t.Errorf("NewPattern failed, expected prefix /home, got %s", pattern.Prefix)
	}
	if pattern.Suffix != "/Packages" {
		t.Errorf("NewPattern failed, expected suffix /Packages, got %s", pattern.Suffix)
	}
	if pattern.Pattern != "/[h-m]arguerite" {
		t.Errorf("NewPattern failed, expected middle pattern /[h-m]arguerite, got %s", pattern.Pattern)
	}
}
