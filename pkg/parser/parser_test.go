package parser_test

import (
	"testing"

	"github.com/YuukanOO/ease/pkg/parser"
)

func TestParser(t *testing.T) {
	t.Run("should parse a package and build an easy representation of declarations", func(t *testing.T) {
		p := parser.New()
		_, err := p.Parse("github.com/YuukanOO/ease/pkg/parser/testdata")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
