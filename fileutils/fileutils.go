package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/marguerite/go-stdlib/dir"
	"github.com/marguerite/go-stdlib/extglob"
	"github.com/marguerite/go-stdlib/internal"
)

//Touch touch a file
func Touch(path string) error {
	_, err := os.Stat(path)

	if err == nil {
		return err
	}
	if os.IsNotExist(err) {
		// create containing directory
		parent := filepath.Dir(path)
		err = dir.MkdirP(parent)
		if err != nil {
			return err
		}
		_, err = os.Create(path)
		if err != nil {
			return err
		}
	}
	return err
}

//cp copy a single file to another file or directory
func cp(source, destination, original string) error {
	// source always exists and can be file only
	s, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	// destination can be non-existent target, file or directory.
	di, err := os.Stat(destination)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil && di.Mode().IsDir() {
		destination = filepath.Join(destination, filepath.Base(source))
		original = ""
	}

	if os.IsNotExist(err) {
		err = dir.MkdirP(filepath.Dir(destination))
		if err != nil {
			return err
		}
	}

	fi, _ := os.Stat(source)
	err = ioutil.WriteFile(destination, s, fi.Mode())
	if err != nil {
		return err
	}

	if len(original) > 0 {
		err := os.RemoveAll(original)
		if err != nil {
			return err
		}
		err = os.Symlink(destination, original)
		if err != nil {
			return err
		}
	}
	return nil
}

func copy(source, destination string, fn func(s, d, o string) error) error {
	// check source status
	si, _ := os.Stat(source)

	// source is a symlink, copy its original content
	if si.Mode()&os.ModeSymlink != 0 {
		link, err := dir.FollowSymlink(source)
		if err != nil {
			return err
		}
		tmp, _ := os.Stat(link)
		si = tmp
		source = link
	}

	// check destination status
	di, err := os.Stat(destination)
	// destination can be non-existent target
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var orig string

	if err == nil && di.Mode()&os.ModeSymlink != 0 {
		// copy to its original file and symlink back
		link, err := dir.FollowSymlink(destination)
		if err != nil {
			return err
		}
		orig = destination
		destination = link
	}

	// copy single file
	if si.Mode().IsRegular() {
		err := fn(source, destination, orig)
		if err != nil {
			return err
		}
		return nil
	}
	// copy directory
	if si.Mode().IsDir() {
		// files can be symlink or actual file
		files, err := dir.Ls(source, true, true)
		if err != nil {
			return err
		}

		for _, f := range files {
			fi, err := os.Stat(f)
			if err != nil {
				// skipped
				continue
			}

			if fi.Mode().IsDir() {
				// skip directory, they'll be created in cp function
				continue
			}

			// keep hierarchy
			p, _ := filepath.Rel(source, filepath.Dir(f))
			dest := filepath.Join(destination, p, filepath.Base(f))

			// f is a symlink, copy its original content
			if fi.Mode()&os.ModeSymlink != 0 {
				link, err := dir.FollowSymlink(f)
				if err != nil {
					continue
				}
				f = link
			}

			err = fn(f, dest, "")
			if err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("source %s has unknown filemode %v", source, si)
}

// Copy like Linux's cp command, copy a file/dirctory to another place.
func Copy(src, dest string) error {
	sources, err := extglob.Expand(internal.Str2bytes(src))
	if err != nil {
		return err
	}
	// sources are always valid files, the check is in extglob's validFunc
	for _, v := range sources {
		err1 := copy(v, dest, cp)
		if err1 != nil {
			return err1
		}
	}
	return nil
}
