package main

import (
	"fmt"
	"os"
)

var cwd string

func init() {
	var err error

	cwd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintln("get current working directory: os:", err.Error()))
	}
}
