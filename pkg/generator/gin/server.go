package gin

import (
	"text/template"

	_ "embed"

	"github.com/YuukanOO/ease/pkg/collection"
	"github.com/YuukanOO/ease/pkg/generator"
	"github.com/YuukanOO/ease/pkg/parser"
	"github.com/YuukanOO/ease/pkg/parser/api"
)

var (
	//go:embed server.go.tmpl
	serverTemplateContent string
	serverTemplate        = template.Must(template.New("").Parse(serverTemplateContent))
)

type ginGenerator struct {
	schema *api.API
}

func New(schema *api.API) generator.Extension {
	return &ginGenerator{
		schema: schema,
	}
}

type data struct {
	generator.Context

	Schema       *api.API
	Fields       *collection.Set[*parser.Type]
	Imports      *collection.Set[*parser.Package]
	Dependencies *collection.Set[*parser.Func]
}

func (g *ginGenerator) Generate(ctx generator.Context) error {
	templateData := &data{
		Context:      ctx,
		Schema:       g.schema,
		Fields:       collection.NewSet[*parser.Type](),
		Imports:      collection.NewSet[*parser.Package](),
		Dependencies: collection.NewSet[*parser.Func](),
	}

	// To build the Server struct, we need to find every handler which as a receiver
	for _, endpoint := range g.schema.Endpoints() {
		// Register each package used by params
		for _, param := range endpoint.Handler().Params() {
			pkg := param.Type().Package()

			if pkg != nil {
				templateData.Imports.Set(pkg.Path(), pkg)
			}
		}

		recv := endpoint.Handler().Recv()

		if recv == nil {
			continue
		}

		recvTyp := recv.Type()

		templateData.Fields.Set(recvTyp.String(), recvTyp)
	}

	for _, recv := range templateData.Fields.Items() {
		var ctor *parser.Func

		for _, fn := range ctx.Funcs() {
			for _, ret := range fn.Returns() {
				if ret.Type() == recv {
					ctor = fn
					break
				}
			}

			// TODO: maybe we can try to find every func matching and raise an error if multiple was found instead
			if ctor != nil {
				break
			}
		}

		if ctor == nil {
			// No ctor found, that's an error!
			continue
		}

		templateData.Dependencies.Set(recv.String(), ctor)

		pkg := ctor.Package()

		if pkg == nil {
			continue
		}
		templateData.Imports.Set(pkg.Path(), pkg)
	}

	return ctx.EmitTemplate("server.go", serverTemplate, templateData)
}
