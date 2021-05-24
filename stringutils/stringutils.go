package stringutils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Contains if s contains any of the substrings, all of the substrings,
// and index of their occurrences
func Contains(s string, substrings ...string) (bool, bool, map[string]int) {
	if len(s) == 0 {
		return false, false, nil
	}

	if len(substrings) == 0 {
		return true, true, nil
	}

	m := make(map[rune][]string, len(substrings))

	for _, v := range substrings {
		if len(v) == 0 {
			continue
		}
		r := []rune(v)[0]
		val := v[1:]
		if len(v) == 1 {
			val = "nil"
		}
		if _, ok := m[r]; !ok {
			m[r] = []string{val}
		} else {
			tmp := m[r]
			tmp = append(tmp, val)
			m[r] = tmp
		}
	}

	var complete int
	m1 := make(map[string]int)

	for i, v := range s {
		if vals, ok := m[v]; ok {
			for _, val := range vals {
				if val == "nil" {
					complete++
					m1[string(v)] = i
					continue
				}
				if len(s)-i < len(val) {
					// fail
					continue
				}
				n := 1
				for _, r := range val {
					if r != []rune(s)[i+n] {
						break
					}
					if n == len(val) {
						m1[string(v)+val] = i
						complete++
						break
					}
					n++
				}
			}
		}
	}

	if complete > 0 {
		if complete == len(substrings) {
			return true, true, m1
		}
		return true, false, m1
	}
	return false, false, nil
}

// Gsub like ruby gsub, match can be string or *regexp.Regexp, replacer can be map (hash in ruby) or string
// the string match can have "|", then the conditional sub strings can be replaced by the value specified in replacer map
// and the replacer string can refer to matched target in match regexp via "\1" or "\2"
func Gsub(orig string, match, replacer interface{}) string {
	// prepare source
	var source [][]string
	if val, ok := match.(*regexp.Regexp); ok {
		if !val.MatchString(orig) {
			return orig
		}
		m := val.FindAllStringSubmatch(orig, -1)
		source = m
	}
	if val, ok := match.(string); ok {
		s := strings.Split(val, "|")
		for _, v := range s {
			source = append(source, []string{v})
		}
	}

	// prepare replace
	if val, ok := replacer.(string); ok {
		var rplr []string
		for i, v := range val {
			if v == '\\' {
				r := []rune{v}
				for j := i + 1; j < len(val); j++ {
					if !unicode.IsDigit([]rune(val)[j]) {
						break
					}
					r = append(r, []rune(val)[j])
				}
				rplr = append(rplr, string(r))
			}
		}

		// actual replace
		for i, v := range source {
			r := val
			for _, v1 := range rplr {
				idxStr := strings.TrimPrefix(v1, "\\")
				idx, _ := strconv.Atoi(idxStr)
				// idx+i, eg: in regexp `(m)arguerite|diagnos(e)`, m is \1 and e is \2, but in gsub they should be all \1
				if idx+i > len(v)-1 {
					err := fmt.Errorf("%d is greater than %d, over-capacity", idx+i, len(v)-1)
					panic(err)
				}
				r = strings.Replace(r, v1, v[idx+i], 1)
			}
			orig = strings.Replace(orig, v[0], r, 1)
		}

		return orig
	}

	if val, ok := replacer.(map[string]string); ok {
		for _, v := range source {
			if val1, ok := val[v[0]]; ok {
				orig = strings.Replace(orig, v[0], val1, 1)
			} else {
				orig = strings.Replace(orig, v[0], "", 1)
			}
		}
	}

	return orig
}
