package parser

import (
	"go/ast"
	"strings"
	"sync"

	"github.com/YuukanOO/ease/pkg/flag"
)

type (
	Result interface {
		Packages() map[string]*Package
		Types() map[string]*Type
		Funcs() map[string]*Func
	}

	// result of the parsing operation for a multitude of packages.
	result struct {
		mu    sync.Mutex
		pkgs  map[string]*Package
		types map[string]*Type
		funcs map[string]*Func
	}
)

func newResult() *result {
	return &result{
		pkgs:  make(map[string]*Package),
		types: make(map[string]*Type),
		funcs: make(map[string]*Func),
	}
}

func (r *result) Packages() map[string]*Package { return r.pkgs }
func (r *result) Types() map[string]*Type       { return r.types }
func (r *result) Funcs() map[string]*Func       { return r.funcs }

// Register the given function declaration.
func (r *result) RegisterFunc(at *FileResult, decl *ast.FuncDecl) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fn := newFunc(at, decl)

	r.funcs[fn.String()] = fn
}

// Register the given type declaration.
func (r *result) RegisterType(at *FileResult, decl *ast.TypeSpec, comment *ast.CommentGroup) {
	r.mu.Lock()
	defer r.mu.Unlock()

	typ := newTypeFromDeclaration(at, decl, comment)

	r.types[typ.String()] = typ
}

// Package returns the package with the given path if it exists or creates
// it if it doesn't.
func (r *result) Package(path string) *Package {
	r.mu.Lock()
	defer r.mu.Unlock()

	sanitizedPath := strings.Trim(path, "\"")
	pkg, found := r.pkgs[sanitizedPath]

	if found {
		return pkg
	}

	pkg = newPackage(sanitizedPath)
	r.pkgs[sanitizedPath] = pkg

	return pkg
}

// Returns the type behind the ident if it exists or creates it if it doesn't.
func (r *result) Type(pkg *Package, ident *ast.Ident) *Type {
	r.mu.Lock()
	defer r.mu.Unlock()

	// If that's a builtin identifier, just remove the pkg
	if IsBuiltin(ident.Name) {
		pkg = nil
	}

	fqn := fullyQualifiedName(pkg, ident.Name)
	typ, found := r.types[fqn]

	if found {
		return typ
	}

	typ = newType(pkg, ident)

	r.types[fqn] = typ

	return typ
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
