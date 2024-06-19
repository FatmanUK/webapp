package main

import (
	"os"
	"errors"
	"strconv"
	"encoding/json"
)

var errDecoding = "Decoding failed"
var errNoKey = "No such key"
var errNoBoolCoerce = "Couldn't coerce boolean"

type JsonOperator interface {
	operate(*JsonConfig, *os.File)
}

type JsonSaver struct {
}

type JsonLoader struct {
}

type JsonConfig struct {
	File string
	Values map[string]string
}

func (*JsonSaver) operate(j *JsonConfig, f *os.File) {
	encoder := json.NewEncoder(f)
	encoder.Encode(j.Values)
}

func (*JsonLoader) operate(j *JsonConfig, f *os.File) {
	decoder := json.NewDecoder(f)
	err := decoder.Decode(j)
	if err != nil {
		panic(errDecoding)
	}
}

func (re *JsonConfig) Init(file string) *JsonConfig {
	re.File = file
	re.Values = make(map[string]string)
	return re
}

func (re *JsonConfig) Debug() string {
	output := `
## Config
___`
	return output
}

func (re *JsonConfig) GetString(key string) string {
	val, isExists := re.Values[key]
	if !isExists {
		panic(errNoKey + ": " + key)
	}
	return val
}

func (re *JsonConfig) SetString(key string, value string) {
	re.Values[key] = value
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
			panic(errNoBoolCoerce + ": " + val)
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
