package main

import (
	"os"

	"github.com/YuukanOO/ease/pkg/parser"
)

func main() {
	dir, _ := os.Getwd()

	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	parser.Parse(dir)
}
