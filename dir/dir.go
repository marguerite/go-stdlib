package dir

import (
	"fmt"
	"os"
	"path/filepath"
)

func errChk(e error) {
	if e != nil {
		panic(e)
	}
}

func readSymlink(path string) (string, error) {
	link, err := os.Readlink(path)
	if err != nil {
		return path, fmt.Errorf("%s points to a broken target.", path)
	}
	link = filepath.Join(filepath.Dir(path), link)
	info, err := os.Stat(link)
	if err != nil {
		return link, fmt.Errorf("%s may be broken.", link)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return readSymlink(link)
	}
	return link, nil
}

func Ls(dir, format string) ([]string, error) {
	var files []string
	e := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		switch format {
		case "dir":
			if info.IsDir() {
				files = append(files, path)
			}
			// the symlinks to directories
			if info.Mode()&os.ModeSymlink != 0 {
				link, err := readSymlink(path)
				if err != nil {
					return err
				}
				f, err := os.Stat(link)
				if err != nil {
					return err
				}
				if f.IsDir() {
					files = append(files, link)
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
				link, err := readSymlink(path)
				if err != nil {
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
