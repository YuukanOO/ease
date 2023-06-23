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
	Imports      *collection.Set[*parser.Package]
	Dependencies []*parser.Func
}

func (g *ginGenerator) Generate(ctx generator.Context) error {
	fields := collection.NewSet[*parser.Type]()
	templateData := &data{
		Context: ctx,
		Schema:  g.schema,
		Imports: collection.NewSet[*parser.Package](),
	}

	// To build the Server struct, we need to find every handler which as a receiver
	for _, endpoint := range g.schema.Endpoints() {
		// If its a raw endpoint, we already import the net/http package so don't need to loop
		// through the params
		if !endpoint.IsRaw() {
			pkg := endpoint.Handler().Package()

			if pkg != nil {
				templateData.Imports.Set(pkg.Path(), pkg)
			}

			// Register each package used by params
			for _, param := range endpoint.Handler().Params() {
				pkg = param.Type().Package()

				if pkg != nil {
					templateData.Imports.Set(pkg.Path(), pkg)
				}
			}
		}

		recv := endpoint.Handler().Recv()

		if recv == nil {
			continue
		}

		recvTyp := recv.Type()

		fields.Set(recvTyp.String(), recvTyp)
	}

	resolved, err := ctx.Funcs().Resolve(fields.Items()...)

	if err != nil {
		return err
	}

	templateData.Dependencies = resolved.Funcs()

	// Add packages needed by dependencies
	for _, fn := range templateData.Dependencies {
		pkg := fn.Package()

		if pkg == nil {
			continue
		}

		templateData.Imports.Set(pkg.Path(), pkg)
	}

	return ctx.EmitTemplate("server.go", serverTemplate, templateData)
}
