package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func LoadConfig(configFile string) (io.ReadCloser, error) {
	fixedFile := os.ExpandEnv(configFile)
	file, err := os.Open(fixedFile)
	if err != nil {
		return nil, newError("config file not readable").Base(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, newError("failed to load config file: ", fixedFile).Base(err).AtWarning()
	}
	return ioutil.NopCloser(bytes.NewReader(content)), nil
}
