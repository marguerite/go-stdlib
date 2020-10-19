package stringutils

import (
	"testing"
)

func TestContainsFullMatch(t *testing.T) {
	s := "NotoColorEmoji.ttf"
	ok, full, matched := Contains(s, "Noto", ".ttf")
	if !ok || !full || len(matched) != 2 {
		t.Error("[stringutils]: Contains full match test failed")
	}
}

func TestContainsPartialMatch(t *testing.T) {
	s := "NotoColorEmoji.ttf"
	ok, full, matched := Contains(s, ".pfa", ".pfb", ".otf", ".ttf")
	if !ok || full || len(matched) != 1 {
		t.Error("[stringutils]: Contains partial match test failed")
	}
}
