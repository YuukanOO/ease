package parser

import (
	"go/ast"
	"go/token"
	"strings"
	"sync"
)

// Base type for all declarations.
type Decl struct {
	lazy       sync.Once
	comments   []*ast.CommentGroup
	name       string
	doc        string
	directives map[string]*Directive
}

func newDeclaration(ident *ast.Ident, comments ...*ast.CommentGroup) *Decl {
	decl := &Decl{
		comments: comments,
	}

	if ident != nil {
		decl.name = ident.Name
	}

	return decl
}

func (d *Decl) IsExported() bool { return token.IsExported(d.name) }
func (d *Decl) Name() string     { return d.name }

func (d *Decl) Doc() string {
	d.parse()
	return d.doc
}

// Returns the directive with the given name if it exists.
func (d *Decl) Directive(name string) (*Directive, bool) {
	d.parse()
	directive, found := d.directives[name]
	return directive, found
}

func (d *Decl) parse() {
	d.lazy.Do(func() {
		d.directives = make(map[string]*Directive)

		var trimmed string

		for _, group := range d.comments {
			if group == nil {
				continue
			}

			for _, line := range group.List {
				trimmed = strings.Trim(line.Text, "/ ")
				directive := tryParseDirective(trimmed)

				if directive != nil {
					d.directives[directive.Name] = directive
				} else {
					d.doc += trimmed + "\n"
				}
			}
		}
	})
}
