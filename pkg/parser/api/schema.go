package api

import (
	"fmt"
	"strings"

	"github.com/YuukanOO/ease/pkg/parser"
)

const (
	FromPath  ParamFrom = iota // Params is extracted from the route path (:id for example)
	FromQuery                  // Params is extracted from the query string
	FromBody                   // Params is extracted from the request body

	MethodOptions Method = "OPTIONS"
	MethodGet     Method = "GET"
	MethodPost    Method = "POST"
	MethodPatch   Method = "PATCH"
	MethodPut     Method = "PUT"
	MethodDelete  Method = "DELETE"
)

const (
	methodDirectiveParam = "method"
	pathDirectiveParam   = "path"
)

type (
	API struct {
		title       string
		description string
		endpoints   []*Endpoint
	}

	// Represents a single endpoint parsed from the API directive and function declaration.
	Endpoint struct {
		handler *parser.Func // Endpoint handler function
		method  Method
		path    string
		params  []*Param
		returns *parser.Var
	}

	Param struct {
		name string // Parameter name
		src  ParamFrom
		decl *parser.Var
	}

	Method    string // HTTP Method
	ParamFrom uint   // Where the parameter is coming from
)

func (s *API) Title() string          { return s.title }
func (s *API) Description() string    { return s.description }
func (s *API) Endpoints() []*Endpoint { return s.endpoints }

func (e *Endpoint) String() string        { return fmt.Sprintf("%s %s", e.method, e.path) }
func (e *Endpoint) Handler() *parser.Func { return e.handler }
func (e *Endpoint) Method() Method        { return e.method }
func (e *Endpoint) Path() string          { return e.path }
func (e *Endpoint) Params() []*Param      { return e.params }
func (e *Endpoint) Returns() *parser.Var  { return e.returns }

func parseEndpoint(directive *parser.Directive, handler *parser.Func) (*Endpoint, error) {
	endpoint := &Endpoint{}

	for name, value := range directive.Params {
		switch name {
		case methodDirectiveParam:
			m, err := parseMethod(value)

			if err != nil {
				return nil, err
			}

			endpoint.method = m
		case pathDirectiveParam:
			endpoint.path = value
		}
	}

	if endpoint.path == "" {
		return nil, ErrInvalidPath
	}

	endpoint.handler = handler

	for _, param := range endpoint.handler.Params() {
		// Context param is a specific one and should not be treated as a request parameter
		if param.Type().IsContext() {
			continue
		}

		endpointParam := &Param{
			name: param.Name(),
			decl: param,
		}

		// Determine the origin of a parameter by checking if its name match a path parameter
		// FIXME: maybe we could use a better way to check it
		if strings.Contains(endpoint.path, ":"+param.Name()) {
			endpointParam.src = FromPath
		} else if endpoint.method == MethodGet {
			endpointParam.src = FromQuery
		} else {
			endpointParam.src = FromBody
		}

		endpoint.params = append(endpoint.params, endpointParam)
	}

	for _, ret := range endpoint.handler.Returns() {
		if ret.Type().IsError() {
			continue
		}

		// The first non builtin return type will be considered as the response type
		endpoint.returns = ret
		break
	}

	return endpoint, nil
}

func parseMethod(value string) (Method, error) {
	switch Method(value) {
	case MethodOptions,
		MethodGet,
		MethodPost,
		MethodPatch,
		MethodPut,
		MethodDelete:
		return Method(value), nil
	case "": // Default to GET if not specified.
		return MethodGet, nil
	default:
		return "", ErrInvalidMethod
	}
}
