package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func WriteConfig(config V2rayConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	//p := "/etc/v2ray"
	p := path.Join(homeDir, ".config/v2ray")
	_, err = os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(p, 0644)
			if err != nil {
				return err
			}
		}
	}
	file := path.Join(p, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LogWriter(logPath string) (io.Writer, error) {
	if logPath != "" {
		_, err := os.Stat(logPath)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir(logPath, 0644)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	filePath := path.Join(logPath, "access.log")
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}
