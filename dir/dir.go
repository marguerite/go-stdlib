package dir

import (
	"os"
	"path/filepath"
)

func errChk(e error) {
	if e != nil {
		panic(e)
	}
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
		case "symlink":
			if info.Mode()&os.ModeSymlink != 0 {
				files = append(files, path)
			}
		default:
			if info.Mode().IsRegular() {
				files = append(files, path)
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
	dirs, e := Ls(dir, "symlink")
	errChk(e)
	return dirs
}

func Lsf(dir string) []string {
	files, e := Ls(dir, "file")
	errChk(e)

	links, e := Ls(dir, "symlink")
	errChk(e)

	for _, v := range links {
		files = append(files, v)
	}

	return files
}
