package api

import (
	"errors"
	"go/ast"

	"github.com/YuukanOO/ease/pkg/parser"
)

var (
	ErrInvalidPath   = errors.New("invalid API path")
	ErrInvalidMethod = errors.New("invalid HTTP method")
)

type (
	Extension interface {
		parser.Extension
		Schema() *API
	}

	apiParser struct {
		schema *API
	}
)

const (
	apiDirective         = "api"
	methodDirectiveParam = "method"
	pathDirectiveParam   = "path"
)

// Builds a new API parser to process files and extract an API schema.
func New() Extension {
	return &apiParser{
		schema: &API{},
	}
}

// Returns the API schema that was build by the parser.
func (p *apiParser) Schema() *API { return p.schema }

func (p *apiParser) Visit(resolver parser.TypeResolver, file *ast.File) error {
	for _, decl := range file.Decls {
		decl, isFunc := decl.(*ast.FuncDecl)

		if !isFunc || !decl.Name.IsExported() {
			continue
		}

		for _, directive := range parser.ParseDirectives(decl.Doc, apiDirective) {
			switch directive.Name {
			case apiDirective:
				endpoint, err := parseEndpoint(directive, resolver, decl)

				if err != nil {
					return err
				}

				p.schema.endpoints = append(p.schema.endpoints, endpoint)
			}
		}
	}

	return nil
}
