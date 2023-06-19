package parser

import (
	"fmt"
	"go/ast"
)

const (
	ContextTypeName = "context.Context"
	ErrorTypeName   = "error"
)

type Type struct {
	*Decl
	pkg  *Package
	spec *ast.TypeSpec
	name string
}

func newType(pkg *Package, ident *ast.Ident, comments ...*ast.CommentGroup) *Type {
	typ := &Type{
		pkg:  pkg,
		name: ident.Name,
	}

	return typ.DefinedBy(ident.Obj, comments...)
}

func (t *Type) Name() string    { return t.name }
func (t *Type) IsContext() bool { return t.String() == ContextTypeName }
func (t *Type) IsError() bool   { return t.String() == ErrorTypeName }
func (t *Type) IsBuiltin() bool { return t.pkg == nil }
func (t *Type) String() string  { return fullyQualifiedName(t.pkg, t.name) }

// If not already set, keep the type object declaration to parse it later.
func (t *Type) DefinedBy(obj *ast.Object, comments ...*ast.CommentGroup) *Type {
	if t.spec == nil && obj != nil && obj.Kind == ast.Typ {
		t.spec = obj.Decl.(*ast.TypeSpec)
		t.Decl = declFromComments(append(comments, t.spec.Doc, t.spec.Comment)...)
	}

	return t
}

func fullyQualifiedName(pkg *Package, name string) string {
	if pkg == nil {
		return name
	}

	return fmt.Sprintf("%s.%s", pkg.Path(), name)
}
