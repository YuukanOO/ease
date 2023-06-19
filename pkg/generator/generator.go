package generator

type (
	Generator interface {
		Generate() error
	}

	Extension interface {
		Generate() error
	}

	generator struct {
		extensions []Extension
	}
)

func New(extensions ...Extension) Generator {
	return &generator{
		extensions: extensions,
	}
}

func (g *generator) Generate() error {
	for _, extension := range g.extensions {
		if err := extension.Generate(); err != nil {
			return err
		}
	}

	return nil
}
