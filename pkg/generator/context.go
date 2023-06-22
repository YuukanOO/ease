package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"github.com/YuukanOO/ease/pkg/collection"
	"github.com/YuukanOO/ease/pkg/crypto"
	"github.com/YuukanOO/ease/pkg/parser"
)

const (
	defaultPermissions     = 0644
	defaultDirPermissions  = 0755
	identifierPrefixLength = 6
)

type (
	Context interface {
		parser.Result

		// Template helpers

		// TODO: merge them both since they are basically the same
		FuncDeclaration(*parser.Func) string // Generates a function declaration
		TypeDeclaration(*parser.Type) string // Generates a field declaration
		Identifier(string, string) string    // Generates a unique identifier for the second string, the first one is used as a prefix, this is useful to avoid name conflicts

		// Generation helpers

		EmitTemplate(string, *template.Template, any) error // Emit a template at the given relative path
		EmitFile(string, []byte) error                      // Emit a file at the given relative path
	}

	context struct {
		parser.Result

		identifiers *collection.Set[string]
		dir         string
	}
)

func newContext(dir string, result parser.Result) Context {
	return &context{
		dir:         dir,
		identifiers: collection.NewSet[string](),
		Result:      result,
	}
}

func (c *context) Identifier(prefix string, key string) string {
	return c.identifiers.SetLazy(key, func() string {
		return fmt.Sprintf("%s_%s", prefix, crypto.Prefix(key, identifierPrefixLength))
	})
}

func (c *context) FuncDeclaration(fn *parser.Func) string {
	if fn.Package() == nil {
		return fn.Name()
	}

	return fmt.Sprintf("%s.%s",
		c.Identifier(fn.Package().Name(), fn.Package().Path()),
		fn.Name(),
	)
}

func (c *context) TypeDeclaration(typ *parser.Type) string {
	if typ.Package() == nil {
		return typ.Name()
	}

	return fmt.Sprintf("%s.%s",
		c.Identifier(typ.Package().Name(), typ.Package().Path()),
		typ.Name(),
	)
}

func (c *context) EmitTemplate(path string, tmpl *template.Template, data any) error {
	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return c.EmitFile(path, buf.Bytes())
}

func (c *context) EmitFile(path string, data []byte) error {
	p, err := c.mkdirAll(path)

	if err != nil {
		return err
	}

	// Format the source file
	data, err = format.Source(data)

	if err != nil {
		return err
	}

	return os.WriteFile(p, data, defaultPermissions)
}

// Creates all directories and resolve the given path before returning it.
func (c *context) mkdirAll(path string) (string, error) {
	fullpath := filepath.Join(c.dir, path)

	if err := os.MkdirAll(filepath.Dir(fullpath), defaultDirPermissions); err != nil {
		return "", err
	}

	return fullpath, nil
}
