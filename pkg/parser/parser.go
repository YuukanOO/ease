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
		Visit(*File) error
	}

	File struct {
		Ast *ast.File
		Pkg *packages.Package
	}

	parserImpl struct {
		extensions []Extension
	}
)

// New creates a new Parser.
func New(extensions ...Extension) Parser {
	return &parserImpl{extensions: extensions}
}

func (p *parserImpl) Parse(packageNames ...string) error {
	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Fset: fset,
		Mode:/*packages.NeedFiles | packages.NeedDeps | packages.NeedImports |*/ packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedModule | packages.NeedName,
		// ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		// 	// Maybe filter files here to prevent some files to be parsed
		// 	return parser.ParseFile(fset, filename, src, parser.ParseComments|parser.AllErrors)
		// },
	}, packageNames...)

	if err != nil {
		return err
	}

	// And process each package files
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			f := &File{
				Ast: file,
				Pkg: pkg,
			}

			for _, ext := range p.extensions {
				if err := ext.Visit(f); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
