package parser

import (
	"go/ast"
	"go/token"
	"strings"
)

// Base type for all declarations.
type Decl struct {
	exported   bool
	name       string
	doc        string
	directives map[string]*Directive
}

func newDeclaration(ident *ast.Ident, comments ...*ast.CommentGroup) *Decl {
	decl := &Decl{
		directives: make(map[string]*Directive),
	}

	if ident != nil {
		decl.exported = token.IsExported(ident.Name)
		decl.name = ident.Name
	}

	var trimmed string

	for _, group := range comments {
		if group == nil {
			continue
		}

		for _, line := range group.List {
			trimmed = strings.Trim(line.Text, "/ ")
			directive := tryParseDirective(trimmed)

			if directive != nil {
				decl.directives[directive.Name] = directive
			} else {
				decl.doc += trimmed + "\n"
			}
		}
	}

	return decl
}

func (d *Decl) IsExported() bool { return d.exported }
func (d *Decl) Name() string     { return d.name }
func (d *Decl) Doc() string      { return d.doc }

func (d *Decl) Directive(name string) (*Directive, bool) {
	directive, found := d.directives[name]
	return directive, found
}
