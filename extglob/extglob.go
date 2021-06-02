// Package extglob basically extglob implements https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html
package extglob

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/marguerite/go-stdlib/internal"
)

// isExtGlobPattern if a string is extglob pattern
func isExtGlobPattern(b []byte) bool {
	i := bytes.Index(b, []byte{'('})
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
	for i, v := range b {
		switch v {
		case '*', '?', '[', '{', '+':
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
		default:
		}
	}
	return false
}

// IsPattern whether a string is a valid pattern
func IsPattern(b []byte, extglob bool) bool {
	if isPlainShellPattern(b) {
		return true
	}
	if extglob {
		return isExtGlobPattern(b)
	}
	return false
}

func escaped(b byte) (bool, error) {
	if b == '\\' {
		if b != PATH_SEPARATOR {
			return true, nil
		}
		return false, errors.New("Windows does not support escaping")
	}
	return false, nil
}

// ListFunc list files, directories/sub-directories based on the globalstar and tailingsep boolean
type ListFunc func(p string, t, t1 bool) ([]string, error)

// list
// if not globalstar, it will return all the files/directories under the current directory
// if globalstar, it will return sub-directories and files in sub-directories too
// if tailingSeparator, it will return directories and sub-directories only
func list(p string, globalstar, tailing bool) ([]string, error) {
	var files []string

	if globalstar {
		// walk the current directory
		err := filepath.Walk(p, func(p1 string, info os.FileInfo, err1 error) error {
			if err1 != nil {
				if os.IsPermission(err1) {
					fmt.Printf("[extglob]: no permission to read %s\n", p1)
					return nil
				}
				return err1
			}
			if tailing {
				if info.IsDir() && p1 != p {
					files = append(files, p1)
				}
				return nil
			}
			// don't include self
			if p1 != p {
				files = append(files, p1)
			}
			return nil
		})

		return files, err
	}

	// return files and directories under the current directory
	f, err := os.Open(p)
	if err != nil {
		return files, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return files, err
	}

	switch mode := info.Mode(); {
	case mode.IsDir():
		names, err := f.Readdirnames(-1)
		for i := 0; i < len(names); i++ {
			names[i] = filepath.Join(p, names[i])
		}
		return names, err
	default:
		return []string{p}, nil
	}
}

// Expand expand extglob pattern to actual files/directories
func Expand(b []byte, options ...bool) ([]string, error) {
	extglob := true
	globalstar := true

	switch len(options) {
	case 0:
	case 1:
		extglob = options[0]
	case 2:
		extglob, globalstar = options[0], options[1]
	default:
		return []string{}, errors.New("only two available options: extglob and globalstar")
	}

	return expand(b, extglob, globalstar, list, valid)
}

// ValidFunc valid if a path exits
type ValidFunc func(p string) bool

func valid(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

func expand(b []byte, extglob, globalstar bool, fn ListFunc, fn1 ValidFunc) ([]string, error) {
	var paths [][]byte
	tmp := bytes.NewBuffer([]byte{})

	for i, v := range b {
		if v == PATH_SEPARATOR || i == len(b)-1 {
			if i != 0 {
				ok, err := escaped(b[i-1])
				if err != nil {
					return []string{}, err
				}
				if ok {
					err := tmp.WriteByte(v)
					if err != nil {
						return []string{}, err
					}
					if i != len(b)-1 {
						continue
					}
				} else {
					if i == len(b)-1 {
						err := tmp.WriteByte(v)
						if err != nil {
							return []string{}, err
						}
					}
				}
			}

			if IsPattern(tmp.Bytes(), extglob) {
				var paths1 [][]byte

				for _, p := range paths {
					files, err := fn(internal.Bytes2str(p), globalstar && (tmp.String() == "**"), i != len(b)-1)

					if err != nil {
						return []string{}, err
					}

					if globalstar && (tmp.String() == "**") {
						for _, v1 := range files {
							paths1 = append(paths1, internal.Str2bytes(v1))
						}
					} else {
						tmp2 := bytes.NewBuffer(tmp.Bytes())
						m := match(files, tmp2, extglob, 0)
						for _, v1 := range m {
							paths1 = append(paths1, internal.Str2bytes(v1))
						}
					}
				}
				paths = paths1
			} else {
				// append the vanilla sub-path, eg paths has ["/home"], after the append it will be ["/home/marguerite"]
				if len(paths) > 0 {
					for i := 0; i < len(paths); i++ {
						p := joinbytes(paths[i], []byte{PATH_SEPARATOR}, tmp.Bytes())
						if fn1(internal.Bytes2str(p)) {
							paths[i] = p
						} else {
							paths = append(paths[:i], paths[i+1:]...)
							i--
						}
					}
				} else {
					// allocs
					paths = append(paths, tmp.Bytes())
				}
			}
			// reset the buffer to store the next sub-path
			tmp = bytes.NewBuffer([]byte{})
			continue
		}
		// simply write none path separator character
		err := tmp.WriteByte(v)
		if err != nil {
			return []string{}, err
		}
	}

	arr := make([]string, 0, len(paths))
	for _, p := range paths {
		if fn1(internal.Bytes2str(p)) {
			arr = append(arr, internal.Bytes2str(p))
		}
	}

	return arr, nil
}

// Match if basenames of file match pattern
func Match(files []string, pattern string) []string {
	return match(files, bytes.NewBufferString(pattern), true, 0)
}

func match(files []string, buf *bytes.Buffer, extglob bool, skip int) []string {
	// n: the number of bytes read from extglob buf
	// n1: the actual byte position of the to-be-match filename
	// length: the total length of extglob patterns without quotes and indicators
	//   eg: ?(g|h), the length is 'g|h' which is 3. we use it to move n
	// useless: whether an extglob pattern is an useless pattern. because we
	//          should not move n1 (the actual position) if no byte matched
	//   eg: filename is marguerite, pattern is mar?(y|z)guerite, the '?(y|z)'
	//       is zero matched, we shouldn't move n1 from 3 to 4, the next match
	//       should still be 'g'.
	// addition: if an extglob pattern can match multiple times. eg mar*(g|h)uerite
	//           the extra matches counts
	//   eg: filename is marggguerite, pattern is mar*(g|h)uerite, we should move
	//       n1 from 3 to 6 not 4 because there're two extra 'g's matched.
	var n, n1, length, useless, addition int
	// indicator: the previous read indicator byte
	var indicator byte

	if len(files) == 0 {
		return files
	}

	for {
		b, err := buf.ReadByte()
		if err == io.EOF {
			break
		}
		switch b {
		case '*':
			if buf.Len() > 0 && buf.Bytes()[0] == '(' {
				indicator = b
			} else {
				if buf.Len() == 0 {
					// all files are valid
					return files
				}

				if len(files) == 0 {
					return files
				}

				if len(files) > 1 {
					files8 := make([]string, 0, len(files))
					for i := 0; i < len(files); i++ {
						files7 := match(files[i:i+1], bytes.NewBuffer(joinbytes([]byte{'*'}, buf.Bytes())), true, n1)
						if len(files7) > 0 {
							files8 = append(files8, files7[0])
						}
					}

					return files8
				}

				files0basename := basenamebytes(files[0])

				var found bool
				for i := n1; i < len(files0basename); i++ {
					m := match(files, bytes.NewBuffer(buf.Bytes()), true, n1+skip+i)
					if len(m) > 0 {
						n1 += i
						found = true
						break
					}
				}

				if !found {
					return []string{}
				}
			}
		case '?':
			if buf.Len() > 0 && buf.Bytes()[0] == '(' {
				indicator = b
			} else {
				// skip
				n1++
			}
		case '+', '@', '!':
			if buf.Len() > 0 && buf.Bytes()[0] == '(' {
				indicator = b
			} else {
				// normal match
				for i := 0; i < len(files); i++ {
					v1 := basenamebytes(files[i])
					if n > len(v1)-1 || b != v1[n+skip] {
						files = append(files[:i], files[i+1:]...)
					}
				}
				n1++
			}
		case '(':
			buf1 := bytes.NewBuffer([]byte{})
			for {
				b2, err := buf.ReadByte()
				if err == io.EOF || b2 == ')' {
					break
				}
				err = buf1.WriteByte(b2)
				if err != nil {
					panic(err)
				}
			}

			// can't use buf1.Len() after matchextglob, because buf1 was read
			tmp := buf1.Len()
			files, useless, addition = matchextglob(files, buf1, indicator, n1+skip)
			// name n is 4, actual n is 8
			// name n is 11, actual n is 14

			n = n + tmp + addition + 1
			length += tmp
			// when indicator is '!', no match is what we want
			if indicator == '!' || useless == 0 {
				n1 += addition + 1
			}
		case '[':
			buf1 := buf.Bytes()
			var n2 int

			// buf1: "gu]erite"
			for {
				b1, _ := buf.ReadByte()
				// ']' can be the first char in the set
				if b1 == ']' && n2 != 0 {
					buf1 = buf1[:n2]
					break
				}
				n2++
			}

			parseBracketRange(&files, buf1, n1+skip)
			n += n2 + 1
			n1++
		case '{':
			buf1 := buf.Bytes()
			var n2 int

			for {
				b1, _ := buf.ReadByte()
				if b1 == '}' {
					buf1 = buf1[:n2]
					break
				}
				n2++
			}

			// buf1: gue,wuef
			// mar{gue,wuef}ite, we can't determine how many bytes are handled
			// so we split the match pattern
			files8 := make([]string, 0, len(files))
			for _, buf2 := range bytes.Split(buf1, []byte(",")) {
				// make a new buffer: guerite
				buf3 := bytes.NewBuffer(joinbytes(buf2, buf.Bytes()))
				files7 := match(files, buf3, extglob, n1+skip)
				for _, f := range files7 {
					files8 = append(files8, f)
				}
			}

			return files8
		default:
			// normal match
			for i := 0; i < len(files); i++ {
				v1 := basenamebytes(files[i])
				if n1+skip > len(v1)-1 || b != v1[n1+skip] {
					files = append(files[:i], files[i+1:]...)
					i--
				}
			}
			n1++
		}
		n++
	}
	return files
}

// matchextglob return 3 values: matched files, useless pattern numbers and the number addtionaly handled.
// useless pattern numbers: eg marguerite vs mar?(y|z)guerite, ?(y|z) matched nothing, it's useless
// the number addtionaly handled: eg marggguerite vs mar+(g|h)uerite, +(g|h) matched 3 'g's, two 'g's are "addtionaly" handled
// because normally we just handle 1 byte
func matchextglob(files []string, buf *bytes.Buffer, indicator byte, n int) ([]string, int, int) {
	// /home/marguerite
	// mar?(y|h)guerite
	// bufs: [y, h]
	bufs := splitextglobpattern(buf)

	switch indicator {
	case '?':
		// match zero or 1 occurrence
		files1 := shellmatchor(files, bufs, n)

		if len(files1) == 0 {
			return files, 1, 0
		}

		return files1, 0, 0
	case '*':
		// match zero or more occurrence
		var bufs1 []*bytes.Buffer
		for _, v := range bufs {
			bufs1 = append(bufs1, bytes.NewBuffer(v.Bytes()))
		}

		files1 := shellmatchor(files, bufs1, n)

		if len(files1) == 0 {
			return files, 1, 0
		}

		var n1 int
		tmp := files1

		// loop until no match
		for {
			var bufs2 []*bytes.Buffer
			for _, v := range bufs {
				bufs2 = append(bufs2, bytes.NewBuffer(v.Bytes()))
			}

			files2 := shellmatchor(files, bufs2, n+n1+1)

			if len(files2) == 0 {
				return tmp, 0, n1
			}
			tmp = files2
			n1++
		}
	case '+':
		// match one or more occurrence
		var bufs1 []*bytes.Buffer
		for _, v := range bufs {
			bufs1 = append(bufs1, bytes.NewBuffer(v.Bytes()))
		}

		files1 := shellmatchor(files, bufs1, n)

		if len(files1) == 0 {
			return []string{}, 1, 0
		}

		var n1 int
		tmp := files1

		// loop until no match
		for {
			var bufs2 []*bytes.Buffer
			for _, v := range bufs {
				bufs2 = append(bufs2, bytes.NewBuffer(v.Bytes()))
			}
			files2 := shellmatchor(files, bufs2, n+n1+1)

			if len(files2) == 0 {
				return tmp, 0, n1
			}
			tmp = files2
			n1++
		}
	case '@':
		// match exactly one occurrence
		files1 := shellmatchor(files, bufs, n)

		if len(files1) == 0 {
			return []string{}, 0, 0
		}

		return files1, 0, 0
	case '!':
		// match anything except one occurrence
		files1 := shellmatchor(files, bufs, n)

		if len(files1) == 0 {
			return files, 1, 0
		}

		return []string{}, 0, 0
	default:
		// unknown indicator, the whole pattern is useless
		return []string{}, 1, 0
	}

}

func shellmatchor(files []string, bufs []*bytes.Buffer, skip int) []string {
	m := make(map[string]struct{})

	for _, buf := range bufs {
		files1 := make([]string, len(files))
		copy(files1, files)
		shellmatch(&files1, buf, skip)
		for _, v := range files1 {
			if _, ok := m[v]; !ok {
				m[v] = struct{}{}
			}
		}
	}

	files3 := make([]string, len(m))
	var j int
	for k := range m {
		files3[j] = k
		j++
	}

	return files3
}

// shellmatch match usual shell pattern, return the matched files
//   skip: skip the 1st n bytes of filename. so the first few bytes should identical in number
func shellmatch(files *[]string, buf *bytes.Buffer, skip int) {
	// files: []string{"/home/marguerite"}
	// buf.String(): /home/mar[gh]uerite
	var n, n1 int

	for {
		b, err := buf.ReadByte()
		if err == io.EOF {
			break
		}

		switch b {

		case '*':
			// any string including null
			if buf.Len() == 0 {
				// all remaining files are valid
				return
			}

			files8 := make([]string, 0, len(*files))
			for i := 0; i < len(*files); i++ {
				files7 := match((*files)[i:i+1], bytes.NewBuffer(joinbytes([]byte{'*'}, buf.Bytes())), true, n1)
				if len(files7) > 0 {
					files8 = append(files8, files7[0])
				}
			}
			for i := 0; i < len(*files); i++ {
				var found bool
				for _, f := range files8 {
					if f == (*files)[i] {
						found = true
						break
					}
				}

				if !found {
					*files = append((*files)[:i], (*files)[i+1:]...)
					i--
				}
			}
			return
		case '?':
			// any byte
			n1++
		case '[':
			// mar[gu]erite
			buf1 := buf.Bytes()
			var n2 int

			// buf1: "gu]erite"
			for {
				b1, _ := buf.ReadByte()
				// ']' can be the first char in the set
				if b1 == ']' && n2 != 0 {
					buf1 = buf1[:n2]
					break
				}
				n2++
			}
			// n not n+1 because in case mar[gh]uerite, n == 3, which is '[', but in actual file /home/marguerite, there is no '[', position 3 is g.
			parseBracketRange(files, buf1, n1+skip)
			n = n + n2 + 1
			n1++
		case '{':
			// {marguerite,wengxuetian}
			buf1 := buf.Bytes()
			var n2 int

			for {
				b1, _ := buf.ReadByte()
				if b1 == '}' {
					buf1 = buf1[:n2]
					break
				}
				n2++
			}

			// buf1: gue,wuef
			// mar{gue,wuef}ite, we can't determine how many bytes are handled
			// so we split the match pattern
			files8 := make([]string, 0, len(*files))
			for _, buf2 := range bytes.Split(buf1, []byte(",")) {
				// make a new buffer: guerite
				buf3 := bytes.NewBuffer(joinbytes(buf2, buf.Bytes()))
				files6 := make([]string, 0, len(*files))
				for _, v := range *files {
					files6 = append(files6, v)
				}

				files7 := match(files6, buf3, true, n1)
				for _, f := range files7 {
					files8 = append(files8, f)
				}
			}

			for i := 0; i < len(*files); i++ {
				var found bool
				for _, f := range files8 {
					if f == (*files)[i] {
						found = true
						break
					}
				}

				if !found {
					*files = append((*files)[:i], (*files)[i+1:]...)
					i--
				}
			}

			return
		default:
			// normal match
			for i := 0; i < len(*files); i++ {
				v1 := basenamebytes((*files)[i])
				if n1+skip > len(v1)-1 || b != v1[n1+skip] {
					*files = append((*files)[:i], (*files)[i+1:]...)
					i--
				}
			}
			n1++
		}
		n++
	}
}

// parseBracketRange parse the range '[]'
func parseBracketRange(files *[]string, buf []byte, n int) {
	// map containing the matching bytes and occurrence times
	m := make(map[byte]int)
	// slice containing the starting and ending byte
	s := makeSlice(buf)
	// the previously looped byte
	var b byte
	// if '!' or '^' is following '[', do a reverse match
	var reverse bool

	// initialize the collator
	collator, err := newCollator()
	if err != nil {
		panic(err)
	}

	fn := func(j int) {
		if _, ok := m[buf[j]]; ok {
			m[buf[j]]++
		} else {
			m[buf[j]] = 1
		}
	}

	// ]abcd-
	for i := 0; i < len(buf); i++ {
		switch buf[i] {
		case '-':
			var next bool
			// the left and right can't be another symbol, eg [a-!] or [!-a]
			if isAlphaNumberic(byte2rune(b)) && (i < len(buf)-1 && isAlphaNumberic(byte2rune(buf[i+1]))) {
				next = true
			}

			// '-' can be first/last in the set
			if i == 0 || i == len(buf)-1 || !next {
				fn(i)
			} else {
				// collating range

				// first remove the just added byte, eg:
				// 'a-z', when meeting '-', 'a' has been added into the map already, drop it first
				if m[b] > 1 {
					m[b]--
				} else {
					delete(m, b)
				}

				// add the starting and ending byte, eg 'a' and 'z' into slice
				s = append(s, []byte{b, buf[i+1]})
				// jump i
				i++
			}
		case ':':
			// unicode class
			if i < len(buf)-1 {
				parseUnicodeClass(files, buf[i+1:], n, reverse)
				// :alpha:m
				idx := bytes.Index(buf[i+1:], []byte{':'})
				i = i + idx + 1
			} else {
				fn(i)
			}
		case '=':
			// collating char
			if i < len(buf)-2 && buf[i+2] == '=' {
				s = append(s, []byte{buf[i+1], buf[i+1]})
				i = i + 2
			} else {
				fn(i)
			}
		case '.':
			// collating symbol
			if i < len(buf)-2 && buf[i+2] == '.' {
				s = append(s, []byte{buf[i+1], buf[i+1]})
				i = i + 2
			} else {
				fn(i)
			}
		case '!', '^':
			if i == 0 {
				reverse = true
			} else {
				fn(i)
			}
		default:
			fn(i)
		}
		b = buf[i]
	}

	for j := 0; j < len(*files); j++ {
		v1 := basenamebytes((*files)[j])
		var found, found1 bool
		// matching bytes
		if len(m) > 0 {
			if _, ok := m[v1[n]]; ok {
				found = true
			}

			if (reverse && found) || (!reverse && !found) {
				*files = append((*files)[:j], (*files)[j+1:]...)
				j--
			}
		}

		// matching range
		if len(s) > 0 {
			for _, v := range s {
				if collator.Compare([]byte{v1[n]}, []byte{v[0]}) >= 0 &&
					collator.Compare([]byte{v1[n]}, []byte{v[1]}) <= 0 {
					found1 = true
					break
				}
			}

			if (reverse && found1) || (!reverse && !found1) {
				*files = append((*files)[:j], (*files)[j+1:]...)
				j--
			}
		}
	}
}

func parseUnicodeClass(files *[]string, buf []byte, n int, reverse bool) {
	// alnum   alpha   ascii   blank   cntrl   digit   graph   lower
	// print   punct   space   upper   word    xdigit
	idx := bytes.Index(buf, []byte{':'})
	if idx < 0 {
		return
	}
	buf = buf[:idx]

	m := make(map[string]func(r rune) bool)

	m["alnum"] = isAlphaNumberic
	m["alpha"] = unicode.IsLetter
	m["ascii"] = func(r rune) bool {
		if r > unicode.MaxASCII {
			return false
		}
		return true
	}
	m["blank"] = func(r rune) bool {
		return r == '\t' || r == ' '
	}
	m["cntrl"] = unicode.IsControl
	m["digit"] = unicode.IsDigit
	m["graph"] = unicode.IsGraphic
	m["lower"] = unicode.IsLower
	m["print"] = unicode.IsPrint
	m["punct"] = unicode.IsPunct
	m["space"] = unicode.IsSpace
	m["upper"] = unicode.IsUpper
	m["word"] = func(r rune) bool { return r == '_' || isAlphaNumberic(r) }
	m["xdigit"] = func(r rune) bool {
		lower := strings.ToLower(string([]rune{r}))
		return unicode.IsDigit(r) || (lower >= "a" && lower <= "f")
	}

	var fn func(r rune) bool

	if val, ok := m[string(buf)]; ok {
		fn = val
	}

	for i := 0; i < len(*files); i++ {
		v1 := basenamebytes((*files)[i])
		var found bool

		if fn != nil && fn(byte2rune(v1[n])) {
			found = true
		}

		if (reverse && found) || (!reverse && !found) {
			*files = append((*files)[:i], (*files)[i+1:]...)
			i--
		}
	}
}
