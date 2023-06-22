package parser

import (
	"go/token"

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

	result := newResult()

	// And process each package files
	for _, pkg := range pkgs {
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
