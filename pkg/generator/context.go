package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
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
	// Represents a declaration inside a package.
	ScopedDecl interface {
		Package() *parser.Package
		Name() string
	}

	Context interface {
		parser.Result

		// Template helpers

		Declaration(ScopedDecl) string    // Generates a declaration from a type or a func
		Identifier(string, string) string // Generates a unique identifier for the second string, the first one is used as a prefix, this is useful to avoid name conflicts

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
	return c.identifiers.SetFunc(key, func() string {
		return fmt.Sprintf("%s_%s", strings.ReplaceAll(prefix, "-", "_"), crypto.Prefix(key, identifierPrefixLength))
	})
}

func (c *context) Declaration(decl ScopedDecl) string {
	if decl.Package() == nil {
		return decl.Name()
	}

	return fmt.Sprintf("%s.%s",
		c.Identifier(decl.Package().Name(), decl.Package().Path()),
		decl.Name(),
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
