package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func errChk(e error) {
	if e != nil {
		panic(e)
	}
}

type NonExistTargetError struct {
  Name string
  Err error
}

func (e NonExistTargetError) Error() string {
  return e.Name
}

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
		//return link, fmt.Errorf("%s may be broken.", link)
    return link, NonExistTargetError{path+" points to an non-existent target "+link, err}
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return ReadSymlink(link)
	}
	return link, nil
}

func Ls(dir, format string) ([]string, error) {
	var files []string
	e := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      if os.IsPermission(err) {
        fmt.Printf("WARNING: no permission to visit %s, skipped.\n", path)
        return nil
      }
      return err
    }
		switch format {
		case "dir":
			if info.IsDir() {
				files = append(files, path)
			}
			// the symlinks to directories
			if info.Mode()&os.ModeSymlink != 0 {
				link, err := ReadSymlink(path)
        if err != nil {
				  if _, ok := err.(NonExistTargetError); ok {
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
          if _, ok := err.(NonExistTargetError); ok {
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

func Lsd(dir string) []string {
	dirs, e := Ls(dir, "dir")
	errChk(e)
	return dirs
}

func Lsl(dir string) []string {
	links, e := Ls(dir, "symlink")
	errChk(e)
	return links
}

func Lsf(dir string) []string {
	files, e := Ls(dir, "file")
	errChk(e)
	return files
}

func Glob(dir string, r *regexp.Regexp) []string {
	m := []string{}
	for _, v := range append(Lsd(dir), Lsf(dir)...) {
		if r.MatchString(v) {
			m = append(m, v)
		}
	}
	return m
}
