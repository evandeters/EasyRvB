package main

import (
	"os"
	"path/filepath"
	"strings"
)

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

func getRoleType(path string) (string, error) {
	var fileData string
	error := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		pathArr := strings.Split(filePath, string(os.PathSeparator))

		if strings.HasSuffix(pathArr[len(pathArr)-1], ".toml") {
			fileData = pathArr[len(pathArr)-1]
		}

		return nil
	})

	if error != nil {
		return "", error
	}

	return strings.Split(fileData, ".")[0], nil
}
