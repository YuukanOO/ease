package parser

import (
	"go/ast"
	"go/types"
	"strings"
	"sync"
)

type (

	// Resolver to parse and retrieve easily processable types.
	TypeResolver interface {
		Obj(*ast.Object) *Type         // Retrieve a type from its object
		Func(decl *ast.FuncDecl) *Func // Retrieve a function from its declaration
	}

	typeResolver struct {
		mu    sync.Mutex
		types map[string]*Type
	}

	// Maps package name/alias to package struct
	ImportsMap map[string]*Package

	scopedTypeResolver struct {
		parent  *typeResolver
		pkg     *Package // Current package, the one we are currently parsing
		imports ImportsMap
	}
)

func NewTypeResolver() *typeResolver {
	return &typeResolver{
		types: make(map[string]*Type),
	}
}

func (r *typeResolver) Resolve(pkg *Package, ident *ast.Ident) *Type {
	r.mu.Lock()
	defer r.mu.Unlock()

	if IsBuiltin(ident.Name) {
		pkg = nil
	}

	fqn := fullyQualifiedName(pkg, ident.Name)
	typ, found := r.types[fqn]

	if found {
		return typ.DefinedBy(ident.Obj)
	}

	typ = newType(pkg, ident)

	r.types[fqn] = typ

	return typ
}

// Builds a new resolver around a specific file and package.
func (r *typeResolver) Scope(currentPackagePath string, file *ast.File) TypeResolver {
	scoped := &scopedTypeResolver{
		parent:  r,
		pkg:     newPackage(currentPackagePath),
		imports: newImportsMap(file.Imports),
	}

	// Parse every types defined in the file to make sure we can later extract informations.
	if file.Scope != nil {
		for _, obj := range file.Scope.Objects {
			scoped.Obj(obj)
		}
	}

	return scoped
}

func (r *scopedTypeResolver) Obj(obj *ast.Object) *Type {
	if obj.Kind != ast.Typ {
		return nil
	}

	spec := obj.Decl.(*ast.TypeSpec)

	return r.parent.Resolve(r.pkg, spec.Name)
}

func (r *scopedTypeResolver) Func(decl *ast.FuncDecl) *Func {
	fn := &Func{
		name: decl.Name.Name,
		pkg:  r.pkg,
	}

	// Process receiver field
	if decl.Recv != nil {
		fn.recv = r.Field(decl.Recv.List[0])
	}

	// Process function parameters
	if decl.Type.Params.List != nil {
		fn.params = make([]*Var, len(decl.Type.Params.List))
	}

	for i, field := range decl.Type.Params.List {
		fn.params[i] = r.Field(field)
	}

	// Process function results
	if decl.Type.Results != nil {
		fn.returns = make([]*Var, len(decl.Type.Results.List))
	}

	for i, field := range decl.Type.Results.List {
		fn.returns[i] = r.Field(field)
	}

	return fn
}

func (r *scopedTypeResolver) Field(field *ast.Field) *Var {
	var name string

	if len(field.Names) > 0 {
		name = field.Names[0].Name
	}

	v := &Var{
		name: name,
	}

	v.underlying, v.kind = r.parseType(field.Type, r.pkg, VarKindUnknown)

	return v
}

func (r *scopedTypeResolver) parseType(expr ast.Expr, pkg *Package, kind VarKind) (*Type, VarKind) {
	switch t := expr.(type) {
	case *ast.Ident:
		var (
			typ      = r.parent.Resolve(pkg, t)
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

		if IsFlagSet(kind, VarKindSlice) {
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
		return r.parseType(t.Sel, r.imports[t.X.(*ast.Ident).Name], kind|VarKindIdent)
	}

	return nil, VarKindUnknown
}

// Builds a new mapping between package name/alias and package path from a raw ImportSpec.
func newImportsMap(imports []*ast.ImportSpec) ImportsMap {
	im := make(ImportsMap, len(imports))

	for _, i := range imports {
		pkgPath := strings.Trim(i.Path.Value, "\"")
		pkg := newPackage(pkgPath)
		if i.Name != nil {
			im[i.Name.Name] = pkg
		} else {
			im[pkg.name] = pkg
		}
	}

	return im
}

// Checks if the given typename is a builtin one.
func IsBuiltin(typeName string) bool {
	return types.Universe.Lookup(typeName) != nil
}

// Check if the given flag is set in the given value.
func IsFlagSet[T ~uint](value T, flag T) bool {
	return value&flag != 0
}
