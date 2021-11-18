package services

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func ReadFiles(relativeDir string) (map[string][]byte, error) {
	fileNameToDataMap := make(map[string][]byte)

	parentDir, err := getParentDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(parentDir, relativeDir)
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ext := filepath.Ext(info.Name())
		if info.Mode().IsRegular() && ext == ".yaml" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			fileNameToDataMap[info.Name()] = data
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return fileNameToDataMap, nil
}

func readFile(relativePath string) ([]byte, error) {
	parentDir, err := getParentDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(parentDir, relativePath)
	return ioutil.ReadFile(path)
}

func getParentDir() (string, error) {
	env, found := os.LookupEnv("ON_CLOUD_FUNCTIONS")
	if !found {
		env = "false"
	}
	deployed, err := strconv.ParseBool(env)
	if err != nil {
		return "", err
	}
	if deployed {
		return "/workspace/serverless_function_source_code", nil
	} else {
		return ".", nil
	}
}
