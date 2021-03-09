// Package extglob basically extglob implements https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html
package extglob

import (
	"bytes"
	"errors"

	internal "github.com/marguerite/go-stdlib/internal"
)

// isExtGlobPattern if a string is extglob pattern
func isExtGlobPattern(b []byte) bool {
	i := bytes.Index(b, internal.Str2bytes("("))
	if i <= 0 {
		return false
	}
	switch b[i-1] {
	case '?', '*', '+', '@', '!':
		return true
	default:
		return false
	}
}

// isPlainShellPattern if a string is plain shell pattern
func isPlainShellPattern(b []byte) bool {
	// pattern chars are:
	// * ? [] {} +
	// usually they will not occur in path
	for _, v := range []byte{'*', '?', '[', '{', '+'} {
		i := bytes.Index(b, []byte{v})
		if i >= 0 {
			if i == len(b)-1 {
				if v == '[' || v == '}' {
					// no pattern can have unclosed brackets
					return false
				}
				return true
			}
			if v == '[' || v == '{' {
				return true
			}
			// it is a ExtGlob pattern
			if b[i+1] == '(' {
				return false
			}
			return true
		}
	}
	return false
}

// IsPattern whether a string is a valid pattern
func IsPattern(b []byte, extglob bool) bool {
	if extglob {
		return isExtGlobPattern(b)
	}
	return isPlainShellPattern(b)
}

// ExtGlob a struct representing pattern like "/home/{marguerite,allen}/Packages",
type ExtGlob struct {
	Prefix  string // "home"
	Suffix  string // "Packages"
	Pattern string // "{marguerite, allen}"
}

type d struct {
	b bool
	s []byte
}

// NewPattern initialize a new pattern from bytes
func NewPattern(b []byte, extglob bool) (ExtGlob, error) {
	// not a path at all
	if !bytes.Contains(b, internal.Str2bytes(PATH_SEPARATOR)) {
		return ExtGlob{}, errors.New("not a valid path at all")
	}

	if !IsPattern(b, extglob) {
		return ExtGlob{internal.Bytes2str(b), "", ""}, errors.New("not a pattern")
	}

	arr := bytes.Split(b, internal.Str2bytes(PATH_SEPARATOR))
	arr1 := make([]d, 0, len(arr))

	for _, v := range arr {
		if IsPattern(v, extglob) {
			arr1 = append(arr1, d{true, v})
		} else {
			arr1 = append(arr1, d{false, v})
		}
	}

	var first, last int

	for i := 0; i < len(arr); i++ {
		if arr1[i].b {
			first = i
			break
		}
	}

	for i := len(arr) - 1; i > -1; i-- {
		if arr1[i].b {
			last = i
			break
		}
	}

	var prefix, pattern, suffix bytes.Buffer

	if first > 0 {
		for i := 0; i < first; i++ {
			if i == 0 {
				writeBytes(&prefix, arr1[i].s, false)
				continue
			}
			writeBytes(&prefix, arr1[i].s, true)
		}
	}

	if last < len(arr1)-1 {
		for i := last + 1; i < len(arr1); i++ {
			if i == len(arr1)-1 {
				writeBytes(&suffix, arr1[i].s, true)
				break
			}
			writeBytes(&suffix, arr1[i].s, true)
		}
	}

	for i := first; i <= last; i++ {
		writeBytes(&pattern, arr1[i].s, true)
	}

	return ExtGlob{prefix.String(), suffix.String(), pattern.String()}, nil
}
