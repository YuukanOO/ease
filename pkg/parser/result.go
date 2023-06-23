package parser

import (
	"go/ast"
	"strings"

	"github.com/YuukanOO/ease/pkg/collection"
	"github.com/YuukanOO/ease/pkg/flag"
)

type (
	Result interface {
		Packages() []*Package
		Types() []*Type
		Funcs() []*Func
	}

	// result of the parsing operation for a multitude of packages.
	result struct {
		pkgs  *collection.Set[*Package]
		types *collection.Set[*Type]
		funcs *collection.Set[*Func]
	}
)

func newResult() *result {
	return &result{
		pkgs:  collection.NewSet[*Package](),
		types: collection.NewSet[*Type](),
		funcs: collection.NewSet[*Func](),
	}
}

func (r *result) Packages() []*Package { return r.pkgs.Items() }
func (r *result) Types() []*Type       { return r.types.Items() }
func (r *result) Funcs() []*Func       { return r.funcs.Items() }

// Register the given function declaration.
func (r *result) RegisterFunc(at *FileResult, decl *ast.FuncDecl) {
	fn := newFunc(at, decl)

	r.funcs.Set(fn.String(), fn)
}

// Register the given type declaration.
func (r *result) RegisterType(at *FileResult, decl *ast.TypeSpec, comment *ast.CommentGroup) {
	typ := newTypeFromDeclaration(at, decl, comment)

	r.types.Set(typ.String(), typ)
}

// Package returns the package with the given path if it exists or creates
// it if it doesn't.
func (r *result) Package(path string) *Package {
	sanitizedPath := strings.Trim(path, "\"")

	return r.pkgs.SetFunc(sanitizedPath, func() *Package {
		return newPackage(sanitizedPath)
	})
}

// Returns the type behind the ident if it exists or creates it if it doesn't.
func (r *result) Type(pkg *Package, ident *ast.Ident) *Type {
	// If that's a builtin identifier, just remove the pkg
	if IsBuiltin(ident.Name) {
		pkg = nil
	}

	fqn := fullyQualifiedName(pkg, ident.Name)

	return r.types.SetFunc(fqn, func() *Type {
		return newType(pkg, ident)
	})
}

type (
	// Maps package name/alias to package struct
	ImportsMap map[string]*Package

	// Result scoped to a specific ast File.
	FileResult struct {
		parent  *result
		pkg     *Package
		imports ImportsMap
	}
)

// Build a FileResult scoped to the given package/ast file.
func (r *result) ParseFile(pkgPath string, file *ast.File) error {
	fileResult := &FileResult{
		parent:  r,
		pkg:     r.Package(pkgPath),
		imports: r.ImportsMap(file.Imports),
	}

	for _, decl := range file.Decls {
		if err := fileResult.visitDeclaration(decl); err != nil {
			return err
		}
	}

	return nil
}

// Builds a new mapping between package name/alias and package path from a raw ImportSpec array.
func (r *result) ImportsMap(imports []*ast.ImportSpec) ImportsMap {
	im := make(ImportsMap, len(imports))

	for _, i := range imports {
		pkg := r.Package(i.Path.Value)
		if i.Name != nil {
			im[i.Name.Name] = pkg
		} else {
			im[pkg.name] = pkg
		}
	}

	return im
}

func (r *FileResult) visitDeclaration(decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				// Only handle types declarations for now
				r.parent.RegisterType(r, s, d.Doc)
			}
		}
	case *ast.FuncDecl:
		r.parent.RegisterFunc(r, d)
	}

	return nil
}

// Parse a single field and returns a Var. It is defined on a scoped FileResult
// object because the import mapping is required to correctly resolve a type.
func (r *FileResult) parseField(field *ast.Field) *Var {
	var name *ast.Ident

	if len(field.Names) > 0 {
		name = field.Names[0]
	}

	v := &Var{
		Decl: newDeclaration(name, field.Doc, field.Comment),
	}

	v.underlying, v.kind = r.parseType(field.Type, r.pkg, VarKindUnknown)

	return v
}

func (r *FileResult) parseType(expr ast.Expr, pkg *Package, kind VarKind) (*Type, VarKind) {
	switch t := expr.(type) {
	case *ast.Ident:
		var (
			typ      = r.parent.Type(pkg, t)
			nextKind = kind
		)

		if !typ.IsBuiltin() {
			nextKind |= VarKindIdent
		} else {
			nextKind |= VarKindBuiltin
		}

		return typ, nextKind
	case *ast.StarExpr:
		nextKind := kind

		if flag.IsSet(kind, VarKindSlice) {
			nextKind |= VarKindSliceOfPointer
		} else {
			nextKind |= VarKindPointer
		}

		return r.parseType(t.X, pkg, nextKind)
	case *ast.ArrayType:
		return r.parseType(t.Elt, pkg, kind|VarKindSlice)
	case *ast.MapType:
		return r.parseType(t.Value, pkg, kind|VarKindMap) // FIXME: handle type of key too maybe
	case *ast.SelectorExpr:
		return r.parseType(t.Sel, r.imports[t.X.(*ast.Ident).Name], kind)
	}

	return nil, VarKindUnknown
}
