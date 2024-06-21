package main

import (
	"os"
	"strconv"
	"encoding/json"
)

type JsonOperator interface {
	operate(*JsonConfig, *os.File)
}

type JsonSaver struct {
}

type JsonLoader struct {
}

type JsonConfig struct {
	file string
	values map[string]string
}

func (*JsonSaver) operate(j *JsonConfig, f *os.File) {
	encoder := json.NewEncoder(f)
	encoder.Encode(j.values)
}

func (*JsonLoader) operate(j *JsonConfig, f *os.File) {
	decoder := json.NewDecoder(f)
	err := decoder.Decode(j)
	if err != nil {
		panic(errDecoding)
	}
}

func (re *JsonConfig) Init(file string) *JsonConfig {
	re.file = file
	re.values = make(map[string]string)
	return re
}

func (re *JsonConfig) Debug() string {
	output := `
## Config
___`
	return output
}

func (re *JsonConfig) GetString(key string) string {
	val, isExists := re.values[key]
	if !isExists {
		panic(errConfigKey + ": " + key)
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
			panic(errBoolCoerce + ": " + val)
	}
}

func (re *JsonConfig) SetBool(key string, value bool) {
	val := "false"
	if value {
		val = "true"
	}
	re.SetString(key, val)
}

func (re *JsonConfig) operate(j JsonOperator, f *os.File) {
	j.operate(re, f)
}

func (re *JsonConfig) Save(f *os.File) {
	j := &JsonSaver{}
	j.operate(re, f)
}

func (re *JsonConfig) Load(f *os.File) {
	j := &JsonLoader{}
	j.operate(re, f)
}
