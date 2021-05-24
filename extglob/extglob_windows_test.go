// +build windows

package extglob

import (
	"testing"
)

func TestEscaped(t *testing.T) {
	_, err := escaped('\\')
	if err == nil {
		t.Error("test escaped failed, it should panic an error, but we got nil")
	}
}
