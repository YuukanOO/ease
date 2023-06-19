package generator

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/YuukanOO/ease/pkg/parser"
)

const defaultPermissions = 0644

type (
	Context interface {
		parser.Result

		EmitTemplate(string, *template.Template, any) error
		EmitFile(string, []byte) error
	}

	context struct {
		parser.Result
		dir string
	}
)

func newContext(dir string, result parser.Result) Context {
	return &context{
		dir:    dir,
		Result: result,
	}
}

func (c *context) EmitTemplate(path string, tmpl *template.Template, data any) error {
	p, err := c.mkdirAll(path)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, defaultPermissions)

	if err != nil {
		return err
	}

	defer file.Close()

	return tmpl.Execute(file, data)
}

func (c *context) EmitFile(path string, data []byte) error {
	p, err := c.mkdirAll(path)

	if err != nil {
		return err
	}

	return os.WriteFile(p, data, defaultPermissions)
}

// Creates all directories and resolve the given path before returning it.
func (c *context) mkdirAll(path string) (string, error) {
	fullpath := filepath.Join(c.dir, path)

	if err := os.MkdirAll(filepath.Dir(fullpath), defaultPermissions); err != nil {
		return "", err
	}

	return fullpath, nil
}
