package main

import (
	"errors"
	"os"

	"github.com/YuukanOO/ease/pkg/generator"
	"github.com/YuukanOO/ease/pkg/generator/ginhttp"
	"github.com/YuukanOO/ease/pkg/parser"
	"github.com/YuukanOO/ease/pkg/parser/api"
)

var ErrNoPackagesGiven = errors.New("missing packages names")

func main() {
	if len(os.Args) < 2 {
		panic(ErrNoPackagesGiven)
	}

	pkgsToAnalyze := os.Args[1:]

	apiParser := api.New()

	p := parser.New(apiParser)

	if err := p.Parse(pkgsToAnalyze...); err != nil {
		panic(err)
	}

	if err := generator.New(ginhttp.New(apiParser.Schema())).Generate(); err != nil {
		panic(err)
	}
}
