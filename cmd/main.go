package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/YuukanOO/ease/pkg/generator/gin"
	"github.com/YuukanOO/ease/pkg/parser/api"
)

var ErrNoPackagesGiven = errors.New("missing packages names")

func main() {
	if len(os.Args) < 2 {
		panic(ErrNoPackagesGiven)
	}

	pkgsToAnalyze := os.Args[1:]

	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	// Default parsers / generators
	apiParser := api.New()
	ginGenerator := gin.New(apiParser.Schema())

	if err := Run(
		WithPackages(pkgsToAnalyze...),
		WithParsers(apiParser),
		WithGenerators(filepath.Join(wd, "generated"), ginGenerator),
	); err != nil {
		panic(err)
	}
}
