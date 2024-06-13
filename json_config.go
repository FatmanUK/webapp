package main

import (
	"os"
	"strconv"
	"encoding/json"
	"errors"
)

type JsonSaver struct {
}

func (*JsonSaver) operate(file string, j *map[string]string) {
	mode := os.O_CREATE|os.O_WRONLY
	f, err := os.OpenFile(file, mode, os.ModePerm)
	defer f.Close()
	if err != nil {
		panic("Conf file not found")
	}
	encoder := json.NewEncoder(f)
	encoder.Encode(j)
}

type JsonLoader struct {
}

func (*JsonLoader) operate(file string, j *map[string]string) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		panic("Conf file not found")
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(j)
	if err != nil {
		panic("Decoding failed")
	}
}

type JsonOperator interface {
	operate(string, *map[string]string)
}

type JsonConfig struct {
	file string
	values map[string]string
}

func (re JsonConfig) debugOutput() string {
	output := `
## Config
  
___`
	return output
}

func (re *JsonConfig) FileExists() bool {
	_, err := os.Stat(re.file)
	return ! errors.Is(err, os.ErrNotExist)
}

func (re *JsonConfig) GetString(key string) string {
	val, isExists := re.values[key]
	if !isExists {
		panic("No such key: " + key)
	}
	return val
}

func (re *JsonConfig) SetString(key string, value string) {
	re.values[key] = value
}

func (re *JsonConfig) GetInt(key string) int {
	val := re.GetString(key)
	rv, _ := strconv.Atoi(val)
	return rv
}

func (re *JsonConfig) SetInt(key string, value int) {
	val := strconv.Itoa(value)
	re.SetString(key, val)
}

func (re *JsonConfig) GetBool(key string) bool {
	val := re.GetString(key)
	switch val {
		case "t", "T", "true", "TRUE", "True", "Y", "y",
			"yes", "YES", "Yes":
			return true
		case "f", "F", "false", "FALSE", "False", "N", "n",
			"no", "NO", "No":
			return false
		default:
			panic("Couldn't coerce boolean: " + val)
	}
}

func (re *JsonConfig) SetBool(key string, value bool) {
	val := "false"
	if value {
		val = "true"
	}
	re.SetString(key, val)
}

func (re *JsonConfig) operate(j JsonOperator) {
	j.operate(re.file, &(re.values))
}

func (re *JsonConfig) Save() {
	re.operate(&JsonSaver{})
}

func (re *JsonConfig) Load() {
	re.operate(&JsonLoader{})
}