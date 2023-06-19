package generator

import "github.com/YuukanOO/ease/pkg/parser"

type (
	Generator interface {
		// Generate needed files from given parser results.
		Generate(parser.Result) error
	}

	Extension interface {
		Generate(Context) error
	}

	generator struct {
		dir        string
		extensions []Extension
	}
)

// Builds a new generator with output dir and given extensions.
func New(dir string, extensions ...Extension) Generator {
	return &generator{
		dir:        dir,
		extensions: extensions,
	}
}

func (g *generator) Generate(result parser.Result) error {
	ctx := newContext(g.dir, result)

	for _, extension := range g.extensions {
		if err := extension.Generate(ctx); err != nil {
			return err
		}
	}

	return nil
}
