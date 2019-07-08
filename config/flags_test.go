package config

import "testing"

func Test_parseArgs(t *testing.T) {
	t.Log(ParseArgs([]string{"-v", "aaa", "-bb=", "-ccc=ddd"}))
	t.Log(ParseArgs([]string{}))
	t.Log(ParseArgs([]string{"-"}))
}

func TestFlagKV_Bool(t *testing.T) {
	flags, _ := ParseArgs([]string{"-v", "aaa", "-bb=", "-ccc=ddd"})
	t.Log(flags.Bool("aaa") == false)
	t.Log(flags.Bool("v") == true)
	t.Log(flags.Bool("v", "ccc") == true)
}
