package parser

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

type (
	// Parser used to process packages names and extract information from them.
	Parser interface {
		// Parse given package names.
		Parse(packageNames ...string) error
	}

	Extension interface {
		// Visit the given file and look for extension related informations.
		Visit(TypeResolver, *ast.File) error
	}

	parserImpl struct {
		extensions []Extension
		resolver   *typeResolver
	}
)

// New creates a new Parser.
func New(extensions ...Extension) Parser {
	return &parserImpl{
		extensions: extensions,
		resolver:   NewTypeResolver(),
	}
}

func (p *parserImpl) Parse(packageNames ...string) error {
	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Fset: fset,
		Mode:/*packages.NeedFiles |  packages.NeedDeps | packages.NeedImports | */ packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedModule | packages.NeedName,
	}, packageNames...)

	if err != nil {
		return err
	}

	// And process each package files
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			scopedResolver := p.resolver.Scope(pkg.PkgPath, file)

			for _, ext := range p.extensions {
				if err := ext.Visit(scopedResolver, file); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
