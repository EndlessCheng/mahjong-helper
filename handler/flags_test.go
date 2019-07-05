package handler

import "testing"

func Test_parseArgs(t *testing.T) {
	t.Log(parseArgs([]string{"-v", "aaa", "-bb=", "-ccc=ddd"}))
	t.Log(parseArgs([]string{}))
	t.Log(parseArgs([]string{"-"}))
}

func TestFlagKV_Bool(t *testing.T) {
	flags, _ := parseArgs([]string{"-v", "aaa", "-bb=", "-ccc=ddd"})
	t.Log(flags.Bool("aaa") == false)
	t.Log(flags.Bool("v") == true)
	t.Log(flags.Bool("v", "ccc") == true)
}
