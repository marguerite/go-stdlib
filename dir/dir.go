package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/marguerite/go-stdlib/pattern"
)

// FollowSymlink follows the path of the symlink recursively and finds out the target it finally points to.
func FollowSymlink(path string) (link string, err error) {
	link, err = os.Readlink(path)
	if err != nil {
		return link, err
	}
	if !filepath.IsAbs(link) {
		link, err = filepath.Abs(filepath.Join(filepath.Dir(path), link))
		if err != nil {
			return link, err
		}
	}
	f, err := os.Stat(link)
	if err != nil {
		return link, err
	}
	if f.Mode()&os.ModeSymlink != 0 {
		return FollowSymlink(link)
	}
	return link, nil
}

// Ls get the file list of directory
// symlink: whether to include symlinks
// recursive: whether to recursively list the second level file list
// kind: if set, will only list the direcories
func Ls(directory string, symlink, recursive bool, kind ...string) (files []string, err error) {
	directories := pattern.Expand(directory)

	for _, v := range directories {
		f, err := os.Open(v)

		if err == nil {
			i, err := f.Stat()
			if err != nil {
				return files, err
			}

			if i.Mode()&os.ModeSymlink != 0 {
				if !symlink {
					// skip
					f.Close()
					continue
				}
				// redirect f to actual file
				link, err := FollowSymlink(v)
				f.Close()
				if err != nil {
					return files, err
				}
				f, err = os.Open(link)
				if err != nil {
					f.Close()
					return files, err
				}
			}

			if i.Mode().IsDir() {
				items, err := f.Readdir(-1)
				if err != nil {
					f.Close()
					return files, err
				}
				for _, j := range items {
					path := filepath.Join(v, j.Name())

					if j.IsDir() {
						files = append(files, path)
					} else if len(kind) == 0 {
						files = append(files, path)
					}

					if recursive && j.IsDir() {
						subfiles, err := Ls(path, symlink, recursive, kind...)
						if err != nil {
							return files, err
						}
						for _, sub := range subfiles {
							files = append(files, sub)
						}
					}
				}
				f.Close()
				continue
			}

			if len(kind) == 0 {
				files = append(files, v)
			}
		}

		f.Close()
	}

	return files, nil
}

// MkdirP create directories for path
func MkdirP(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return err
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

// Glob glob actual files via the pattern, pattern can be *regexp.Regexp or string
// when *regexp.Regexp is used, base is a must.
func Glob(patt interface{}, opts ...interface{}) (files []string, err error) {
	if len(opts) > 2 {
		return files, fmt.Errorf("opts just have two values: base and exclusion")
	}

	var base string
	if len(opts) > 0 {
		if val, ok := opts[0].(string); ok {
			base = val
		}
	}

	switch val := patt.(type) {
	case *regexp.Regexp:
		matches, err := Ls(base, true, true)
		if err != nil {
			return files, err
		}

		for _, v := range matches {
			if val.MatchString(v) {
				if len(opts) > 1 {
					if val1, ok := opts[1].(*regexp.Regexp); ok {
						if val1.MatchString(v) {
							continue
						}
					}
				}
				files = append(files, v)
			}
		}
	case string:
		// string match
		matches := pattern.Expand(val)
		for _, v := range matches {
			if _, err := os.Stat(v); err == nil {
				files = append(files, v)
			}
		}
	}
	return files, err
}
