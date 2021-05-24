package extglob

import (
	"bytes"
	"errors"
	"regexp"
	"testing"
)

func BenchmarkExpand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		expand([]byte("/home/[mn]arguerite"), true, true, list1, valid1)
	}
}

func BenchmarkExpand1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		expand1([]byte("/home/[mn]arguerite"), true, true, list1, valid1)
	}
}

func BenchmarkExpand2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		expand2("/home/[mn]arguerite", list1)
	}
}

func expand2(s string, fn ListFunc) ([]string, error) {
	re := regexp.MustCompile(`^/home/(m|n)arguerite$`)
	files, _ := fn("", true, true)
	for i := 0; i < len(files); i++ {
		if !re.MatchString(files[i]) {
			files = append(files[:i], files[i+1:]...)
		}
	}
	return files, nil
}

func expand1(b []byte, extglob, globalstar bool, fn ListFunc, fn1 ValidFunc) ([]string, error) {
	// FIXME two lines below creates allocs
	paths := make([]*bytes.Buffer, 0, bytes.Count(b, []byte{PATH_SEPARATOR}))
	tmp := bytes.NewBuffer(make([]byte, 0, 50))

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
				var paths1 []*bytes.Buffer
				for _, p := range paths {
					files, err := fn(p.String(), globalstar && (tmp.String() == "**"), i != len(b)-1)
					if err != nil {
						return []string{}, err
					}

					var m []string

					if globalstar && (tmp.String() == "**") {
						m = files
					} else {
						m = match(files, tmp, extglob, 0)
					}

					for _, v1 := range m {
						paths1 = append(paths1, bytes.NewBufferString(v1))
					}
				}
				paths = paths1
			} else {
				// append the vanilla sub-path, eg paths has ["/home"], after the append it will be ["/home/marguerite"]
				if len(paths) > 0 {
					for _, p := range paths {
						err := p.WriteByte(PATH_SEPARATOR)
						if err != nil {
							return []string{}, err
						}
						n, err := p.Write(tmp.Bytes())
						if err != nil {
							return []string{}, err
						}
						if n != len(tmp.Bytes()) {
							return []string{}, errors.New("not fully written")
						}
					}
				} else {
					// allocs
					paths = append(paths, bytes.NewBuffer(tmp.Bytes()))
				}
			}
			// reset the buffer to store the next sub-path
			tmp.Reset()
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
		if fn1(p.String()) {
			arr = append(arr, p.String())
		}
	}

	return arr, nil
}
