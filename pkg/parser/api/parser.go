package api

import (
	"go/ast"

	"github.com/YuukanOO/ease/pkg/parser"
)

type (
	apiParser struct {
		Result *API
	}

	API struct {
		Endpoints []*Endpoint
	}

	Endpoint struct {
		Method Method `ini:"method"`
		Path   string `ini:"path"`
	}

	Method string
)

const (
	MethodGet    Method = "GET"
	MethodPost   Method = "POST"
	MethodPatch  Method = "PATCH"
	MethodPut    Method = "PUT"
	MethodDelete Method = "DELETE"

	annotation = "//ease:api"
)

func New() parser.Extension {
	return &apiParser{}
}

func (p *apiParser) Init() error {
	p.Result = &API{}
	return nil
}

func (p *apiParser) Visit(file *parser.File) error {
	for _, decl := range file.Ast.Decls {
		decl, ok := decl.(*ast.FuncDecl)

		if !ok || !decl.Name.IsExported() {
			continue
		}

		var endpoint Endpoint

		annotationFound, err := parser.ParseAnnotation(annotation, decl.Doc, &endpoint)

		if err != nil {
			return err
		}

		if !annotationFound {
			continue
		}

		p.Result.Endpoints = append(p.Result.Endpoints, &endpoint)

		if decl.Recv == nil {
			continue
		}

		// var typeName string
		// recvExpr := decl.Recv.List[0].Type

		// switch t := recvExpr.(type) {
		// case *ast.StarExpr:
		// 	typeName = t.X.(*ast.Ident).Name
		// case *ast.Ident:
		// 	typeName = t.Name
		// }

		// obj := file.Pkg.Types.Scope().Lookup(typeName)
		// fmt.Println(obj)
	}

	return nil
}
