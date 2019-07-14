package fileutils

import (
	"errors"
	"fmt"
	"github.com/marguerite/dirutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// Touch a file or directory
func Touch(path string, isDir ...bool) error {
	// process isDir, we don't allow > 1 arguments.
	if len(isDir) > 1 {
		return errors.New("isDir is a symbol indicating whether the path is a target directory. You shall not pass two arguments.")
	}
	ok := false
	if len(isDir) == 1 {
		ok = isDir[0]
	}

	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			if ok {
				err := dirutil.MkdirP(path)
				return err
			}

			dir := filepath.Dir(path)
			if dir != "." {
				err := dirutils.MkdirP(dir)
				if err != nil {
					fmt.Println("Can not create containing directory " + dir)
					return error
				}
			}
			f, err := os.Create(path)
			defer f.Close()
			if err != nil {
				if os.IsPermission(err) {
					fmt.Println("WARNING: no permission to create " + path + ", skipped...")
					return nil
				}
				return err
			}
		} else {
			return fmt.Errorf("Another unhandled non-IsNotExist PathError occurs: %s", err.Error())
		}
	}

	return nil
}

func Copy(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}

	// source is a directory
	if stat.IsDir() {
		return fmt.Errorf("%s is a directory", src)
	}

	// source is a symlink
	if stat.Mode()&os.ModeSymlink == os.ModeSymlink {
		fmt.Printf("%s is a symlink, following the original one", src)
		org, err := os.Readlink(src)
		if err != nil {
			return err
		}
		src = org
	}

	if info, err := os.Stat(dst); !os.IsNotExist(err) {
		// dst is a directory
		if info.IsDir() {
			basename := filepath.Base(src)
			dst = filepath.Join(dst, basename)
		} else {
			err = os.Remove(dst)
			if err != nil {
				return err
			}
		}
	}

	in, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, in, stat.Mode())
	if err != nil {
		return err
	}
	return nil
}

// HasPrefixSuffixInGroup if a string's prefix/suffix matches one in group
// b trigger's prefix match
func HasPrefixSuffixInGroup(s string, group []string, b bool) bool {
	prefix := "(?i)"
	suffix := ""
	if b {
		prefix += "^"
	} else {
		suffix += "$"
	}

	for _, v := range group {
		re := regexp.MustCompile(prefix + regexp.QuoteMeta(v) + suffix)
		if re.MatchString(s) {
			return true
		}
	}
	return false
}
