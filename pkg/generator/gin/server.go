package gin

import (
	"text/template"

	_ "embed"

	"github.com/YuukanOO/ease/pkg/generator"
	"github.com/YuukanOO/ease/pkg/parser"
	"github.com/YuukanOO/ease/pkg/parser/api"
)

var (
	//go:embed main.go.tmpl
	mainTemplateContent string
	mainTemplate        = template.Must(template.New("").Parse(mainTemplateContent))
)

type ginGenerator struct {
	schema *api.API
}

func New(schema *api.API) generator.Extension {
	return &ginGenerator{
		schema: schema,
	}
}

func (g *ginGenerator) Generate(ctx generator.Context) error {
	recvs := make(map[string]*parser.Type)
	dependencies := make(map[string]*parser.Func)
	imports := make(map[string]*parser.Package)

	// To build the Server struct, we need to find every handler which as a receiver
	for _, endpoint := range g.schema.Endpoints() {
		recv := endpoint.Handler().Recv()
		if recv == nil {
			continue
		}

		recvTyp := recv.Type()

		recvs[recvTyp.String()] = recvTyp
	}

	for t, recv := range recvs {
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

		dependencies[t] = ctor
		pkg := ctor.Package()
		imports[pkg.Path()] = pkg
	}

	return ctx.EmitTemplate("main.go", mainTemplate, data{
		Schema:  g.schema,
		Imports: imports,
		Fields:  recvs,
	})
}

type data struct {
	Schema  *api.API
	Fields  map[string]*parser.Type
	Imports map[string]*parser.Package
}
