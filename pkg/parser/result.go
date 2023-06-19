package parser

import (
	"go/ast"
	"strings"
	"sync"

	"github.com/YuukanOO/ease/pkg/flag"
)

// Result of the parsing operation for a multitude of packages.
type Result struct {
	mu    sync.Mutex
	pkgs  map[string]*Package
	types map[string]*Type
	funcs map[string]*Func
}

func newResult() *Result {
	return &Result{
		pkgs:  make(map[string]*Package),
		types: make(map[string]*Type),
		funcs: make(map[string]*Func),
	}
}

// Package returns the package with the given path if it exists or creates
// it if it doesn't.
func (r *Result) Package(path string) *Package {
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

// Returns the function behind the ident if it exists or creates it if it doesn't.
func (r *Result) Func(pkg *Package, ident *ast.Ident, comments ...*ast.CommentGroup) *Func {
	r.mu.Lock()
	defer r.mu.Unlock()

	fqn := fullyQualifiedName(pkg, ident.Name)
	fn, found := r.funcs[fqn]

	if found {
		return fn
	}

	fn = &Func{
		Decl: declFromComments(comments...),
		pkg:  pkg,
		name: ident.Name,
	}

	r.funcs[fqn] = fn

	return fn
}

// Returns the type behind the ident if it exists or creates it if it doesn't.
func (r *Result) Type(pkg *Package, ident *ast.Ident, comments ...*ast.CommentGroup) *Type {
	r.mu.Lock()
	defer r.mu.Unlock()

	// If that's a builtin identifier, just remove the pkg
	if IsBuiltin(ident.Name) {
		pkg = nil
	}

	fqn := fullyQualifiedName(pkg, ident.Name)
	typ, found := r.types[fqn]

	if found {
		return typ.DefinedBy(ident.Obj, comments...)
	}

	typ = newType(pkg, ident, comments...)
	r.types[fqn] = typ

	return typ
}

type (
	// Maps package name/alias to package struct
	ImportsMap map[string]*Package

	// Result scoped to a specific ast File.
	FileResult struct {
		parent  *Result
		pkg     *Package
		imports ImportsMap
	}
)

// Build a FileResult scoped to the given package/ast file.
func (r *Result) ParseFile(pkgPath string, file *ast.File) error {
	fileResult := &FileResult{
		parent:  r,
		pkg:     r.Package(pkgPath),
		imports: r.ImportsMap(file.Imports),
	}

	for _, decl := range file.Decls {
		if err := fileResult.parseDeclaration(decl); err != nil {
			return err
		}
	}

	return nil
}

// Builds a new mapping between package name/alias and package path from a raw ImportSpec array.
func (r *Result) ImportsMap(imports []*ast.ImportSpec) ImportsMap {
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

func (r *FileResult) parseDeclaration(decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			if err := r.parseSpec(spec, d.Doc); err != nil {
				return err
			}
		}
	case *ast.FuncDecl:
		return r.parseFunc(d)
	}

	return nil
}

func (r *FileResult) parseSpec(spec ast.Spec, comments ...*ast.CommentGroup) error {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		// Only handle types for now
		r.parent.Type(r.pkg, s.Name, comments...)
	case *ast.ValueSpec, *ast.ImportSpec:
		return nil
	}

	return nil
}

func (r *FileResult) parseFunc(decl *ast.FuncDecl) error {
	fn := r.parent.Func(r.pkg, decl.Name, decl.Doc)

	// Process receiver field
	if decl.Recv != nil {
		fn.recv = r.parseField(decl.Recv.List[0])
	}

	// Process function parameters
	if decl.Type.Params.List != nil {
		fn.params = make([]*Var, len(decl.Type.Params.List))
	}

	for i, field := range decl.Type.Params.List {
		fn.params[i] = r.parseField(field)
	}

	// Process function results
	if decl.Type.Results != nil {
		fn.returns = make([]*Var, len(decl.Type.Results.List))
	}

	for i, field := range decl.Type.Results.List {
		fn.returns[i] = r.parseField(field)
	}

	return nil
}

func (r *FileResult) parseField(field *ast.Field) *Var {
	var name string

	if len(field.Names) > 0 {
		name = field.Names[0].Name
	}

	v := &Var{
		Decl: declFromComments(field.Doc, field.Comment),
		name: name,
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
