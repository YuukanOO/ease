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
	file *FileResult
	pkg  *Package
	decl *ast.TypeSpec
}

func newType(pkg *Package, ident *ast.Ident) *Type {
	return &Type{
		Decl: newDeclaration(ident),
		pkg:  pkg,
	}
}

func newTypeFromDeclaration(at *FileResult, decl *ast.TypeSpec, comment *ast.CommentGroup) *Type {
	return &Type{
		Decl: newDeclaration(decl.Name, decl.Doc, comment),
		file: at,
		pkg:  at.pkg,
		decl: decl,
	}
}

func (t *Type) IsContext() bool { return t.String() == ContextTypeName }
func (t *Type) IsError() bool   { return t.String() == ErrorTypeName }
func (t *Type) IsBuiltin() bool { return t.pkg == nil }
func (t *Type) String() string  { return fullyQualifiedName(t.pkg, t.name) }

func fullyQualifiedName(pkg *Package, name string) string {
	if pkg == nil {
		return name
	}

	return fmt.Sprintf("%s.%s", pkg.Path(), name)
}
