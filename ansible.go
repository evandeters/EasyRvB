package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Role struct {
	Name string
	Max  int
}

func getAnsibleRoles(path string) ([]string, error) {
	var roles []string
	error := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		pathArr := strings.Split(filePath, string(os.PathSeparator))

		if pathArr[len(pathArr)-2] == "roles" && !strings.HasPrefix(pathArr[len(pathArr)-1], ".") {
			roles = append(roles, filePath)
		}

		return nil
	})

	if error != nil {
		return nil, error
	}

	return roles, nil
}

func parseRoleMetadata(path string) (map, error {
	f, err := os.Open(path + string(os.PathSeparator) + ".metadata")
    if err != nil {
        return err
    }
	defer f.Close()

    scanner := bufio.NewScanner(f)
    categoryRegex := regexp.MustCompile(`^\[(.*)?\]$`)
    configRegex := regexp.MustCompile(`^(\w+?)\s+=\s+(.*)$`)
    configMap := make(map[string]map[string]string)
    for scanner.Scan() {
        category := ""
        line := scanner.Text()
        if categoryRegex.MatchString(line) {
            category = categoryRegex.FindStringSubmatch(line)[1]
            configMap[category] = make(map[string]string)
            fmt.Println(category)
        } else if configRegex.MatchString(line) {
            config := configRegex.FindStringSubmatch(line)
            configMap[category][config[1]] = config[2]
            fmt.Println(configMap[category])
        }
    }

    return configMap, nil
}
