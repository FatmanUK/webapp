package main

import (
	"testing"
)

type M struct {
	*testing.T
}

func (t *M) LogTestSucceedOrFailB(b bool) {
	if b {
		t.Logf("Test succeeded")
	} else {
		t.Fatalf("Test failed")
	}
}

func TestInit(t *testing.T) {
	{
		var s StringList
		s.Init()
		t.Logf(string(s.delimiter))
	}
	{
		var s StringList
		s.InitR('#')
		t.Logf(string(s.delimiter))
	}
	{
		var s StringList
		s.InitS("x1;y2;z3")
		t.Logf(s.String())
	}
	{
		var s StringList
		s.InitRS('#', "x1#y2#z3")
		t.Logf(s.String())
	}
	{
		var s StringList
		(&s).Init()
		t.Logf(string(s.delimiter))
	}
	{
		var s StringList
		(&s).InitR('#')
		t.Logf(string(s.delimiter))
	}
	{
		var s StringList
		(&s).InitS("x1;y2;z3")
		t.Logf(s.String())
	}
	{
		var s StringList
		(&s).InitRS('#', "x1#y2#z3")
		t.Logf(s.String())
	}
}

func TestContains(t *testing.T) {
	{
		var s StringList
		s.InitS("x1;y2;z3")
		m := &M{t}
		m.LogTestSucceedOrFailB(true  == s.BContainsS("z3"))
		m.LogTestSucceedOrFailB(false == s.BContainsS("x3"))
		m.LogTestSucceedOrFailB(true  == s.BContainsS("y2"))
		m.LogTestSucceedOrFailB(false == s.BContainsS("y3"))
		m.LogTestSucceedOrFailB(true  == s.BContainsS("x1"))
	}
	{
		var s StringList
		s.InitRS('?', "x1?y2?z3")
		m := &M{t}
		m.LogTestSucceedOrFailB(true  == s.BContainsS("z3"))
		m.LogTestSucceedOrFailB(false == s.BContainsS("x3"))
		m.LogTestSucceedOrFailB(true  == s.BContainsS("y2"))
		m.LogTestSucceedOrFailB(false == s.BContainsS("y3"))
		m.LogTestSucceedOrFailB(true  == s.BContainsS("x1"))
	}
}
