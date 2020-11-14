package stringutils

import (
	"regexp"
	"testing"
)

const uri = "https://github.com/marguerite/diagnose"

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

func TextGsubRegexpMatchWithNormalTextReplacer(t *testing.T) {
	re := regexp.MustCompile(`(m)arguerite|diagnos(e)`)
	if Gsub(uri, re, "opensuse") != "https://github.com/opensuse/opensuse" {
		t.Error("[stringutils]: Gsub regexp match with normal text replacer failed")
	}
}

func TextGsubRegexpMatchWithMatchedTargetReferenceReplacer(t *testing.T) {
	re := regexp.MustCompile(`(m)arguerite|diagnos(e)`)
	if Gsub(uri, re, "(\\1)") != "https://github.com/(m)/(e)" {
		t.Error("[stringutils]: Gsub regexp match with matched target reference replacer failed")
	}
}

func TextGsubStringMatchWithNormalTextReplacer(t *testing.T) {
	if Gsub(uri, "marguerite", "opensuse") != "https://github.com/opensuse/diagnose" {
		t.Error("[stringutils]: Gsub string match with normal text replacer failed")
	}
}

func TextGsubStringMatchWithMapReplacer(t *testing.T) {
	m := map[string]string{"marguerite": "opensuse", "diagnose": "diagnostic"}
	if Gsub(uri, "marguerite|diagnose", m) != "https://github.com/opensuse/diagnostic" {
		t.Error("[stringutils]: Gsub string match with map replacer failed")
	}
}
