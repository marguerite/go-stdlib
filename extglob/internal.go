package extglob

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/marguerite/go-stdlib/internal"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// intialize the collator with current environment LANG, LC_COLLATE or LC_ALL
func newCollator() (*collate.Collator, error) {
	var lang string
	for _, v := range []string{"LC_ALL", "LC_COLLATE", "LANG"} {
		// val: zh_CN.UTF-8
		val := os.Getenv(v)
		if len(val) > 0 {
			i := strings.Index(val, ".")
			if i > 0 {
				val = val[:i]
			}

			if val == "C" {
				lang = "en_us"
			} else {
				lang = val
			}
			break
		}
	}

	tag, err := language.Parse(lang)
	if err != nil {
		if inv, ok := err.(language.ValueError); ok {
			return nil, errors.New(inv.Subtag())
		}
		return nil, fmt.Errorf("can not found a language tag for lang %s", lang)
	}

	return collate.New(tag, collate.IgnoreCase), nil
}

// makeSlice make slice based on the total number of '-', '=', and '.'
func makeSlice(buf []byte) [][]byte {
	var num int
	for _, v := range []byte{'-', '=', '.'} {
		// buf: -a-z-=c=.!..-
		buf1 := buf
		for {
			i := bytes.Index(buf1, []byte{v})

			if i < 0 {
				break
			}

			switch v {
			case '-':
				if i != 0 && i != len(buf1)-1 && isAlphaNumberic(byte2rune(buf1[i-1])) && isAlphaNumberic(byte2rune(buf1[i+1])) {
					num++
				}
				buf1 = buf1[i+1:]
			case '=':
				j := bytes.Index(buf1[i+1:], []byte{'='})
				if j == 1 {
					num++
					buf1 = buf1[i+j+2:]
				} else {
					buf1 = buf1[i+1:]
				}
			case '.':
				j := bytes.Index(buf1[i+1:], []byte{'.'})
				if j == 1 {
					num++
					buf1 = buf1[i+j+2:]
				} else {
					buf1 = buf1[i+1:]
				}
			}
		}
	}
	return make([][]byte, 0, num)
}

func isAlphaNumberic(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

func byte2rune(b byte) rune {
	r, _ := utf8.DecodeRune([]byte{b})
	return r
}

// splitextglobpattern split ?(y|z) to y and z
func splitextglobpattern(buf *bytes.Buffer) []*bytes.Buffer {
	buf1 := buf.Bytes()
	bufs := make([]*bytes.Buffer, 0, bytes.Count(buf1, []byte{'|'}))
	// closed: if a pair of [] appears or an unclosed [
	var closed, previous int

	for i := 0; i < buf.Len(); i++ {
		if i == buf.Len()-1 {
			bufs = append(bufs, bytes.NewBuffer(buf1[previous:]))
			break
		}
		switch buf1[i] {
		case '[':
			closed = i
		case ']':
			closed = 0
		case '|':
			if closed == 0 {
				bufs = append(bufs, bytes.NewBuffer(buf1[previous:i]))
				previous = i + 1
			}
		default:
		}
	}
	return bufs
}

func joinbytes(b ...[]byte) []byte {
	var total int
	for _, v := range b {
		total += len(v)
	}

	b1 := make([]byte, total)

	l := 0
	for j := 0; j < len(b); j++ {
		for i := 0; i < len(b[j]); i++ {
			b1[l+i] = b[j][i]
		}
		l += len(b[j])
	}

	return b1
}

func basenamebytes(s string) []byte {
	i := len(s) - 1
	for i >= 0 && s[i] != PATH_SEPARATOR {
		i--
	}
	return internal.Str2bytes(s)[i+1:]
}
