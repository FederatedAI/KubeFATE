package api

import (
	"fmt"
	"runtime/debug"
)

// todo: Definition of all returned information

const (
	SUCCESS     = "success"
	USEREXISTED = "User is existed."
)

func recovery() {
	if r := recover(); r != nil {
		fmt.Println("Recovered", r)
		debug.PrintStack()
	}
}
