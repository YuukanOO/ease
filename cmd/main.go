package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/YuukanOO/ease/pkg/generator"
	"github.com/YuukanOO/ease/pkg/generator/gin"
	"github.com/YuukanOO/ease/pkg/parser"
	"github.com/YuukanOO/ease/pkg/parser/api"
)

var ErrNoPackagesGiven = errors.New("missing packages names")

func main() {
	if len(os.Args) < 2 {
		panic(ErrNoPackagesGiven)
	}

	pkgsToAnalyze := os.Args[1:]

	// Instantiate specific parsers
	apiParser := api.New()

	result, err := parser.
		New(apiParser).
		Parse(pkgsToAnalyze...)

	if err != nil {
		panic(err)
	}

	// Instantiate specific generators
	ginGenerator := gin.New(apiParser.Schema())

	// Retrieve the working directory to compute the output path
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	if err := generator.
		New(filepath.Join(wd, "generated"), ginGenerator).
		Generate(result); err != nil {
		panic(err)
	}
}
