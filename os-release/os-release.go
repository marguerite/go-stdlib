package osrelease

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
)

// OsRelease map refer to /etc/os-release
type OsRelease map[string]string

func (release *OsRelease) init() {
	f, _ := ioutil.ReadFile("/etc/os-release")
	scanner := bufio.NewScanner(bytes.NewReader(f))
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), "=")
		(*release)[strings.ToLower(arr[0])] = strings.Replace(arr[1], "\"", "", -1)
	}
}

// Name The name of your Linux OS
func Name() string {
	release := OsRelease{}
	release.init()
	return release["pretty_name"]
}

// Version The version of your Linux OS
func Version() int {
	release := OsRelease{}
	release.init()
	version, _ := strconv.Atoi(release["version_id"])
	return version
}
