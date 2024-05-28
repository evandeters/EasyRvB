package main

import (
	"os"
	"path/filepath"
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

func parseRoleMetadata(path string) error {
	f := os.Open(path + string(os.PathSeparator) + ".metadata")
	defer f.close()

}
