package main

import (
	"errors"
	"fmt"
	"os"

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

	schema := apiParser.Schema()

	fmt.Println(schema)
}
