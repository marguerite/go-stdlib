package extglob

import (
	"testing"

	internal "github.com/marguerite/go-stdlib/internal"
)

var (
	extGlobTestPattern = map[string]string{
		"star1":       "/home/*(m)arguerite/",
		"question1":   "/home/?(m)arguerite/",
		"at":          "/home/mar@(g|h)uerite/",
		"add1":        "/home/mar+(g)uerite/",
		"exclamation": "/home/!(a)aguerite/",
	}

	plainShellTestPattern = map[string]string{
		"star":     "/home/marguerite/**",
		"question": "/home/marguer?te/",
		"add":      "/home/mar?+/",
		"bracket":  "/home/[^a]arguerite/",
		"bracket1": "/home/[h-m]arguerite",
		"curly":    "/home/{marguerite,allen}",
	}
)

func TestIsExtGlobPattern(t *testing.T) {
	for _, v := range extGlobTestPattern {
		bol := isExtGlobPattern(internal.Str2bytes(v))
		if !bol {
			t.Errorf("isExtGlobPattern %s failed, expected true, got %t", v, bol)
		}
	}
}

func TestIsPlainShellPattern(t *testing.T) {
	for _, v := range plainShellTestPattern {
		bol := isPlainShellPattern(internal.Str2bytes(v))
		if !bol {
			t.Errorf("isExtGlobPattern %s failed, expected true, got %t", v, bol)
		}
	}
}

func TestNewPatternNotPath(t *testing.T) {
	b := internal.Str2bytes("havealotoffun")
	if _, err := NewPattern(b, false); err == nil {
		t.Errorf("NewPattern %b failed, expected 'not a valid path at all' error, got nil", b)
	}
}
