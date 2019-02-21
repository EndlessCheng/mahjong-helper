package main

import (
	"os/exec"
	"os"
	"runtime"
)

var clearFuncMap map[string]func() //create a map for storing clear funcs

func init() {
	clearFuncMap = make(map[string]func()) //Initialize it
	clearFuncMap["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearFuncMap["darwin"] = clearFuncMap["linux"]
	clearFuncMap["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clearConsole() {
	clearFunc, ok := clearFuncMap[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		clearFunc() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
