package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func WriteConfig(configFile string, config V2rayConfig) error {
	configName := path.Base(configFile)
	dir := configFile[:strings.LastIndex(configFile, configName)]

	err := MkdirIfNotExist(dir)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LogWriter(logPath string) (io.Writer, error) {
	err := MkdirIfNotExist(logPath)
	if err != nil {
		return nil, err
	}

	filePath := path.Join(logPath, "access.log")
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func MkdirIfNotExist(dir string) error {
	if dir != "" {
		_, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir(dir, 0755)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
