package main

import (
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		panic("please provide a command")
	}

	execute(args[1])
}
