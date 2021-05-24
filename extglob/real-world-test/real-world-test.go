// real-word-test: do the extglob tests in your real home

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/marguerite/go-stdlib/extglob"
	"github.com/marguerite/go-stdlib/fileutils"
)

func random2stars(path string) string {
	arr := strings.Split(path, string(extglob.PATH_SEPARATOR))
	if len(arr[0]) == 0 {
		arr = arr[1:]
	}
	rand.Seed(time.Now().UnixNano())
	p := rand.Intn(len(arr) - 1)
	arr[p] = "**"
	str := strings.Join(arr, string(extglob.PATH_SEPARATOR))
	if str[0] != '/' {
		str = "/" + str
	}
	return str
}

func main() {
	// touch the dummy file
	wd, _ := os.Getwd()

	test := filepath.Join(wd, "miku.ogv")
	fileutils.Touch(test)
	defer os.Remove(test)

	patterns := []string{
		// *
		filepath.Join(random2stars(wd), "mi*u.ogv"),
		filepath.Join(wd, "mi*ku.ogv"),
		filepath.Join(wd, "miku.og*"),
		// ?
		filepath.Join(wd, "mi?u.o?v"),
		// []
		filepath.Join(wd, "mi[kg]u.ogv"),
		filepath.Join(wd, "mi[a-k]u.ogv"),
		filepath.Join(wd, "mi[-k-]u.ogv"),
		filepath.Join(wd, "mi[]kg]u.ogv"),
		filepath.Join(wd, "mi[^u-z]u.ogv"),
		filepath.Join(wd, "mi[!u-z]u.ogv"),
		filepath.Join(wd, "mi[:alpha:]u.ogv"),
		filepath.Join(wd, "mi[=k=]u.ogv"),
		filepath.Join(wd, "miku[...]ogv"),
		// {}
		filepath.Join(wd, "{miku.ogv, real-word-test.go}"),
		// extglob ?
		filepath.Join(wd, "?(m|g)iku.ogv"),
		filepath.Join(wd, "?(n|g)miku.ogv"),
		// extglob *
		filepath.Join(wd, "*(m|g)iku.ogv"),
		filepath.Join(wd, "*(n|g)miku.ogv"),
		// extglob +
		filepath.Join(wd, "+(m|g)iku.ogv"),
		// extglob @
		filepath.Join(wd, "@(m|g)iku.ogv"),
		// extglob !
		filepath.Join(wd, "!(n|k)iku.ogv"),
	}

	for i, pattern := range patterns {
		var t bool
		if i == 0 {
			t = true
		}
		files, _ := extglob.Expand([]byte(pattern), true, t)
		if len(files) == 0 || filepath.Base(files[0]) != "miku.ogv" {
			fmt.Printf("pattern: %s\n", pattern)
			fmt.Printf("expected %s, got %v\n", test, files)
		} else {
			fmt.Printf("pattern %s passed\n", pattern)
		}
	}
}
