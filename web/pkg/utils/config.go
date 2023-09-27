package utils

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Load 会通过`filepath`的扩展名判断json还是yaml，然后调用相应的配置获取
func Load(filePath string, config interface{}) error {
	cl := configLoader{filePath: filePath}
	return cl.load(config)
}

type configLoader struct {
	filePath string
}

func (cl configLoader) load(config interface{}) error {
	file, err := os.OpenFile(cl.filePath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	fileName := file.Name()

	lastIndex := strings.LastIndex(fileName, ".")
	if lastIndex == -1 {
		return errors.New("file must have a type suffix")
	}

	fileType := fileName[lastIndex+1:]
	switch strings.ToLower(fileType) {
	case "json":
		return cl.loadJSON(file, config)
	case "yaml":
		fallthrough
	case "yml":
		return cl.loadYAML(file, config)
	default:
		return errors.New("file type not supported")
	}
}

func (cl configLoader) loadJSON(file io.Reader, config interface{}) error {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, config)
}

func (cl configLoader) loadYAML(file io.Reader, config interface{}) error {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, config)
}
