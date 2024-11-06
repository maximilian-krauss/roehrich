package main

import (
	"github.com/maximilian-krauss/roehrich/cmd"
	"os"
)

func main() {
	cmd.Execute(os.Args[1:])
}
