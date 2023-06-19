package parser

import (
	"go/ast"
	"strings"
)

// Base type for all declarations.
type Decl struct {
	doc        string
	directives map[string]*Directive
}

func declFromComments(comments ...*ast.CommentGroup) *Decl {
	decl := &Decl{
		directives: make(map[string]*Directive),
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
