package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func IsSymlink(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if f.Mode()&os.ModeSymlink != 0 {
		return true, nil
	}
	return false, nil
}

func IsRegular(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if f.Mode().IsRegular() {
		return true, nil
	}
	return false, nil
}

func Touch(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("WARNING: no permission to create %s, skipped...\n", path)
				return nil
			}
			return err
		}
		defer f.Close()
	}
	return nil
}

func Remove(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
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
