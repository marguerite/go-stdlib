package pattern

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func toRegex(s string) *regexp.Regexp {
	if s == "**" {
		return regexp.MustCompile(`^.*$`)
	}
	// escape the dots first
	s = strings.Replace(s, ".", "\\.", -1)
	s = strings.Replace(s, "?", ".", -1)
	s = strings.Replace(s, "*", ".*", -1)
	s = strings.Replace(s, "{", "(", -1)
	s = strings.Replace(s, ",", "|", -1)
	s = strings.Replace(s, "}", ")", -1)
	return regexp.MustCompile("^" + s + "$")
}

// Expand a `pattern` to actual files. pattern indicators are :
// **, *, ?, {}, []
func Expand(s string) (branch []string) {
	if strings.ContainsAny(s, "*?{[") {
		var prefix string
		i := strings.IndexAny(s, "*?{[")

		if i > 0 {
			p, _ := filepath.Abs(filepath.Dir(s[:i]))
			prefix = p + PATH_SEPARATOR
		} else {
			prefix = root
		}

		segments := strings.Split(strings.TrimPrefix(s, prefix), PATH_SEPARATOR)

		for j, segment := range segments {
			var parents, tmp []string

			if j > 0 {
				parents = branch
			} else {
				parents = append(parents, prefix)
			}

			for _, parent := range parents {
				if strings.ContainsAny(segment, "*?{[") {
					f, err := os.Open(parent)
					if err != nil {
						if j > 0 {
							continue
						}
						panic(err)
					}

					i, err := f.Stat()
					if err != nil {
						panic(err)
					}

					if i.Mode().IsDir() {
						sub, err := f.Readdir(-1)
						if err != nil {
							panic(err)
						}

						f.Close()

						for _, v := range sub {
							if toRegex(segment).MatchString(v.Name()) {
								tmp = append(tmp, filepath.Join(parent, v.Name()))
							}
						}
					}
				} else {
					tmp = append(tmp, filepath.Join(parent, segment))
				}
			}

			branch = tmp
		}
	} else {
		branch = append(branch, s)
	}
	return branch
}
