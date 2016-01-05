package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ImportCsvFile(filename string) []*SshConfig {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rawData, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	nodes := make([]*SshConfig, 0)
	for _, node := range rawData {
		nodes = append(nodes, &SshConfig{Name: node[0], User: node[1], Password: node[2], Port: node[3]})
	}

	return nodes

}

func PathChomp(path string) string {
	likePosix, err := regexp.Compile(".*/.*")
	if err != nil {
		return ""
	}

	if likePosix.MatchString(path) {
		s := strings.Split(path, "/")
		return s[len(s)-1]
	}

	return ""
}

func PathCutFile(path string) string {
	likePosix, err := regexp.Compile(".*/.*")
	if err != nil {
		return ""
	}

	if likePosix.MatchString(path) {
		s := strings.Split(path, "/")
		sl := s[:len(s)-1]
		return strings.Join(sl, "/")
	}

	return ""
}

func PrependParent(root, path string) string {
	split := strings.Split(path, "/")
	droptop := strings.Join(split[2:], "/")

	return filepath.ToSlash(filepath.Join(root, droptop))
}

func GlobIn(s string) bool {
	for _, b := range s {
		if string(b) == "*" {
			return true
		}
	}
	return false
}

func JoinPath(path, file string) string {
	return filepath.ToSlash(filepath.Join(PathCutFile(path), file))
}

func IsRegularFile(f os.FileInfo) bool {
	if f.Mode().IsRegular() {
		return true
	}
	return false
}

func ApplyFilter(text, searchKey string) string {
	pattern, _ := regexp.Compile(searchKey)
	if pattern.MatchString(text) {
		return text
	} else {
		return ""
	}
}

func PrependHost(host string, data string) string {
	return fmt.Sprintf("[ %s ]: %s", host, data)
}

func ProcessOutput(data string, host *SshConfig, env *Enviroment) []string {
	out := make([]string, 0)
	var filtered string

	for _, d := range strings.Split(data, "\n") {
		if env.filter.on {
			filtered = ApplyFilter(d, env.filter.key)

			if !Empty(filtered) {
				out = append(out, PrependHost(host.Name, filtered))
			}

		} else {
			out = append(out, PrependHost(host.Name, d))
		}
	}

	return out
}

func LineMatch(line []byte, searchKey string) bool {
	pattern, _ := regexp.Compile(searchKey)
	if pattern.Match(line) {
		return true
	} else {
		return false
	}
}

func Empty(str ...string) bool {
	var count int
	for _, s := range str {
		if s != "" {
			count++
		}
	}

	if count > 0 {
		return false
	}

	return true
}

func WithBash(command string) string {
	return "/bin/bash -l -c " + "'" + command + "'"
}

func ImportPubKeyConfig(hosts, key, userName, port string) []*SshConfig {
	nodes := make([]*SshConfig, 0)
	if port == "" {
		port = "22"
	}

	for _, host := range strings.Split(hosts, ",") {
		nodes = append(nodes, &SshConfig{Name: host, User: userName, Key: key, Port: port})
	}

	return nodes
}
