package main

import (
	"strings"
)

type StringList struct {
	delimiter rune
	Members []string
}

func (re *StringList) Init() *StringList {
	re.InitR(';')
	return re
}

func (re *StringList) InitR(d rune) *StringList {
	re.delimiter = d
	return re
}

func (re *StringList) InitS(s string) *StringList {
	re.InitRS(';', s)
	return re
}

func (re *StringList) InitRS(d rune, s string) *StringList {
	re.InitR(d)
	re.Members = strings.Split(s, string(re.delimiter))
	return re
}

func (re *StringList) String() string {
	return strings.Join(re.Members, string(re.delimiter))
}

func (re *StringList) BContainsS(s string) bool {
	for _, g := range re.Members {
		if g == s {
			return true
		}
	}
	return false
}
