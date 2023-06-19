package api

import (
	"errors"

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

const apiDirective = "api"

// Builds a new API parser to process files and extract an API schema.
func New() Extension {
	return &apiParser{
		schema: &API{},
	}
}

// Returns the API schema that was build by the parser.
func (p *apiParser) Schema() *API { return p.schema }

func (p *apiParser) Visit(result parser.Result) error {
	for _, fn := range result.Funcs() {
		if !fn.IsExported() {
			continue
		}

		api, hasApiDirective := fn.Directive(apiDirective)

		if !hasApiDirective {
			continue
		}

		endpoint, err := parseEndpoint(api, fn)

		if err != nil {
			return err
		}

		p.schema.endpoints = append(p.schema.endpoints, endpoint)
	}

	return nil
}
