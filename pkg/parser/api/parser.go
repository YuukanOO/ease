package api

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"

	"github.com/YuukanOO/ease/pkg/parser"
)

var (
	ErrInvalidPath = errors.New("invalid API path")
)

type (
	apiParser struct {
		Schema *API
	}

	API struct {
		Endpoints []*Endpoint
	}

	Endpoint struct {
		Method Method
		Path   string
	}

	Method string
)

const (
	MethodOptions Method = "OPTIONS"
	MethodGet     Method = "GET"
	MethodPost    Method = "POST"
	MethodPatch   Method = "PATCH"
	MethodPut     Method = "PUT"
	MethodDelete  Method = "DELETE"
	MethodInvalid Method = ""

	api = "api"
)

func New() parser.Extension {
	return &apiParser{
		Schema: &API{},
	}
}

func (p *apiParser) Visit(file *parser.File) error {
	for _, decl := range file.Ast.Decls {
		decl, ok := decl.(*ast.FuncDecl)

		if !ok || !decl.Name.IsExported() {
			continue
		}

		for _, directive := range parser.ParseDirectives(decl.Doc, api) {
			switch directive.Name {
			case api:
				endpoint, err := parseEndpoint(directive, decl)

				obj := file.Pkg.Types.Scope().Lookup("TodoService").Type().(*types.Named)
				sig := obj.Method(0).Type().(*types.Signature)
				p1 := sig.Params().At(1)

				fmt.Println(p1)

				if err != nil {
					return err
				}

				p.Schema.Endpoints = append(p.Schema.Endpoints, endpoint)
			}
		}
	}

	return nil
}

func parseEndpoint(directive *parser.Directive, decl *ast.FuncDecl) (*Endpoint, error) {
	endpoint := &Endpoint{}

	for name, value := range directive.Params {
		switch name {
		case "method":
			endpoint.Method = methodFromRawValue(value)
		case "path":
			endpoint.Path = value
		}
	}

	// Default to GET if not specified.
	if endpoint.Method == MethodInvalid {
		endpoint.Method = MethodGet
	}

	if endpoint.Path == "" {
		return nil, ErrInvalidPath
	}

	return endpoint, nil
}

func methodFromRawValue(value string) Method {
	switch value {
	case "OPTIONS":
		return MethodOptions
	case "GET":
		return MethodGet
	case "POST":
		return MethodPost
	case "PATCH":
		return MethodPatch
	case "PUT":
		return MethodPut
	case "DELETE":
		return MethodDelete
	default:
		return MethodInvalid
	}
}
