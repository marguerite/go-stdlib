package osrelease

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
)

var release map[string]string

func init() {
	f, err := ioutil.ReadFile("/etc/os-release")
	if err != nil {
		panic("no /etc/os-release")
	}

	scanner := bufio.NewScanner(bytes.NewReader(f))
  release = make(map[string]string)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), "=")
		release[strings.ToLower(arr[0])] = strings.Replace(arr[1], "\"", "", -1)
	}
}

// Name The name of your Linux OS
func Name() string {
	return release["pretty_name"]
}

// Version The version of your Linux OS
func Version() float64 {
	version, _ := strconv.ParseFloat(release["version_id"], 64)
	return version
}
