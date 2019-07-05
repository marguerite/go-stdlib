package dirutils

import (
	"errors"
	"fmt"
	"github.com/marguerite/util/slice"
	"os"
	"path/filepath"
	"regexp"
)

// ErrNonExistTarget ErrNonExistTarget is used to indicate the target a symlink points to actually does not exist on the filesystem.
type ErrNonExistTarget struct {
	Path string
	Link string
}

func (e ErrNonExistTarget) Error() string {
	return e.Path + "points to an non-existent target " + e.Link
}

// ReadSymlink follows the path of the symlink recursively and finds out the target it finally points to.
func ReadSymlink(path string) (string, error) {
	link, err := os.Readlink(path)
	if err != nil {
		return path, err
	}
	if !filepath.IsAbs(link) {
		link = filepath.Join(filepath.Dir(path), link)
	}
	info, err := os.Stat(link)
	if err != nil {
		return link, ErrNonExistTarget{path, link}
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return ReadSymlink(link)
	}
	return link, nil
}

func ls(d string, kind string) ([]string, error) {
	files := []string{}
	e := filepath.Walk(string(d), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("WARNING: no permission to visit %s, skipped.\n", path)
				return nil
			}
			return err
		}
		switch kind {
		case "dir":
			if info.IsDir() {
				files = append(files, path)
			}
			// the symlinks to directories
			if info.Mode()&os.ModeSymlink != 0 {
				link, err := ReadSymlink(path)
				if err != nil {
					if _, ok := err.(ErrNonExistTarget); ok {
						// the symlink points to an non-existent target, ignore
						fmt.Printf("WARNING: %s points to an non-existent target %s.\n", path, link)
						return nil
					}
					return err
				}
				f, err := os.Stat(link)
				if err != nil {
					return err
				}
				if f.IsDir() {
					files = append(files, path)
				}
			}
		case "symlink":
			if info.Mode()&os.ModeSymlink != 0 {
				files = append(files, path)
			}
		default:
			if info.Mode().IsRegular() {
				files = append(files, path)
			}
			// the symlinks to actual files
			if info.Mode()&os.ModeSymlink != 0 {
				link, err := ReadSymlink(path)
				if err != nil {
					if _, ok := err.(ErrNonExistTarget); ok {
						// the symlink points to an non-existent target, ignore
						fmt.Printf("WARNING: %s points to an non-existent target %s.\n", path, link)
						return nil
					}
					return err
				}
				f, err := os.Stat(link)
				if err != nil {
					return err
				}
				if f.Mode().IsRegular() {
					files = append(files, link)
				}
			}
		}
		return nil
	})
	return files, e
}

// Ls Takes a directory and the kind of file to be listed, returns the list of file and the possible error. Kind supports: dir, symlink, defaults to file.
func Ls(d string, kinds ...string) ([]string, error) {

	if len(kinds) == 0 {
		return fn(d, "")
	}

	if len(kinds) == 1 {
		return fn(d, kinds[0])
	}

	f := []string{}
	for _, kind := range kinds {
		i, err := fn(d, kind)
		if err != nil {
			// f is incomplete
			return f, err
		}
		slice.Concat(&f, i)
	}
	return f, nil
}

// MkdirP create directories for path
func MkdirP(path string) error {
	p := filepath.Dir(path)
	fmt.Printf("Creating directory: %s\n", p)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.MkdirAll(p, os.ModeDir)
		if err != nil {
			fmt.Printf("Can not create directory %s\n", p)
			return err
		}
		fmt.Printf("%s created\n", p)
	} else {
		fmt.Printf("%s exists already\n", p)
		return nil
	}
	return nil
}

func Glob(dir string, pattern interface{}, fn ...func(string, ...string) ([]string, error)) ([]string, error) {
	if len(fn) == 0 {
		fn = append(fn, Ls)
	}

	files, err := fn[0](dir)
	if err != nil {
		return []string{}, err
	}
	var re []*regexp.Regexp
	switch v := pattern.(type) {
	case *regexp.Regexp:
		re = append(re, v)
	case []*regexp.Regexp:
		re = v
	case string:
		re = append(re, regexp.MustCompile(v))
	default:
		return []string{}, errors.New("Unsupported pattern type. Supported: *regexp.Regexp, []*regexp.Regexp, string.")
	}
	m := []string{}
	for _, f := range files {
		for _, r := range re {
			if r.MatchString(f) {
				m = append(m, f)
			}
		}
	}
	return m, nil
}
