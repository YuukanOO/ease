package parser

import (
	"fmt"
	"go/token"
	"path"

	"golang.org/x/tools/go/packages"
)

type (
	// Parser used to process packages names and extract information from them.
	Parser interface {
		// Parse given package names.
		Parse(packageNames ...string) (Result, error)
	}

	Extension interface {
		Visit(Result) error
	}

	parser struct {
		extensions []Extension
	}
)

// New creates a new Parser.
func New(extensions ...Extension) Parser {
	return &parser{
		extensions: extensions,
	}
}

func (p *parser) Parse(packageNames ...string) (Result, error) {
	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Fset: fset,
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedModule | packages.NeedName,
	}, packageNames...)

	if err != nil {
		return nil, err
	}

	var generatedModulePath string
	mod := findMainModule(pkgs)

	if mod != nil {
		generatedModulePath = path.Join(mod.Path, "generated") // FIXME: no hardcoded value!
	}

	result := newResult()

	// And process each package files
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			fmt.Println(pkg.Errors) // FIXME: replace with logger call
			continue
		}

		// Skip the generated package
		if pkg.PkgPath == generatedModulePath {
			continue
		}

		for _, file := range pkg.Syntax {
			if err = result.ParseFile(pkg.PkgPath, file); err != nil {
				return nil, err
			}
		}
	}

	// And finally, visit each extension
	for _, extension := range p.extensions {
		if err = extension.Visit(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func findMainModule(pkgs []*packages.Package) *packages.Module {
	for _, pkg := range pkgs {
		if pkg.Module.Main {
			return pkg.Module
		}
	}

	return nil
}
