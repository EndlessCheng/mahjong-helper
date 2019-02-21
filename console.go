package main

import (
	"os/exec"
	"os"
	"runtime"
)

var clearFunc map[string]func() //create a map for storing clear funcs

func init() {
	clearFunc = make(map[string]func()) //Initialize it
	clearFunc["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearFunc["darwin"] = clearFunc["linux"]
	clearFunc["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clearConsole() {
	value, ok := clearFunc[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
