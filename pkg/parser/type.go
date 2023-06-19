package parser

import (
	"fmt"
	"go/ast"
	"strings"
	"sync"

	"github.com/YuukanOO/ease/pkg/crypto"
)

const (
	ContextTypeName = "context.Context"
	ErrorTypeName   = "error"
)

type Type struct {
	*Decl
	lazy  sync.Once
	alias string
	file  *FileResult
	pkg   *Package
	decl  *ast.TypeSpec
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

func (t *Type) IsContext() bool   { return t.String() == ContextTypeName }
func (t *Type) IsError() bool     { return t.String() == ErrorTypeName }
func (t *Type) IsBuiltin() bool   { return t.pkg == nil }
func (t *Type) Package() *Package { return t.pkg }
func (t *Type) String() string    { return fullyQualifiedName(t.pkg, t.name) }

// Returns a unique alias representing this type, useful for naming variables.
func (t *Type) Alias() string {
	t.lazy.Do(func() {
		if t.pkg != nil {
			t.alias = fmt.Sprintf("%s_%s", strings.ToLower(t.name), crypto.Prefix(t.String(), aliasSuffixSize))
		}
	})

	return t.alias
}

// Returns a string representing the type declaration prefixed with the package alias.
func (t *Type) Declaration() string {
	if t.pkg == nil {
		return t.Name()
	}

	return fmt.Sprintf("%s.%s", t.pkg.Alias(), t.Name())
}

func fullyQualifiedName(pkg *Package, name string) string {
	if pkg == nil {
		return name
	}

	return fmt.Sprintf("%s.%s", pkg.Path(), name)
}
