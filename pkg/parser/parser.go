package parser

import (
	"fmt"
	"go/parser"
	"go/token"
)

func Parse(dir string) error {
	mod, err := findModule(dir, 0)

	if err != nil {
		return err
	}

	fmt.Println(mod)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments|parser.AllErrors)

	if err != nil {
		panic(err)
	}

	for _, p := range pkgs {
		pkg := mod.CreatePackage(dir, p)

		fmt.Println(pkg)
	}

	return nil
}
