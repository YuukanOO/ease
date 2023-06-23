package parser

import (
	"fmt"
	"go/ast"
	"sync"
)

type (
	Vars  []*Var
	Funcs []*Func

	Func struct {
		*Decl
		lazy    sync.Once
		file    *FileResult
		decl    *ast.FuncDecl
		pkg     *Package
		recv    *Var
		params  Vars
		returns Vars
	}
)

func newFunc(at *FileResult, decl *ast.FuncDecl) *Func {
	return &Func{
		Decl: newDeclaration(decl.Name, decl.Doc),
		file: at,
		pkg:  at.pkg,
		decl: decl,
	}
}

func (f *Func) Recv() *Var {
	f.parse()
	return f.recv
}

func (f *Func) Params() Vars {
	f.parse()
	return f.params
}

func (f *Func) Returns() Vars {
	f.parse()
	return f.returns
}

func (f *Func) Package() *Package { return f.pkg }
func (f *Func) String() string    { return fullyQualifiedName(f.pkg, f.name) }

func (f *Func) parse() {
	f.lazy.Do(func() {
		// Process receiver field
		if f.decl.Recv != nil {
			f.recv = f.file.parseField(f.decl.Recv.List[0])
		}

		// Process function parameters
		if f.decl.Type.Params.List != nil {
			f.params = make(Vars, len(f.decl.Type.Params.List))

			for i, field := range f.decl.Type.Params.List {
				f.params[i] = f.file.parseField(field)
			}
		}

		// Process function results
		if f.decl.Type.Results != nil {
			f.returns = make(Vars, len(f.decl.Type.Results.List))

			for i, field := range f.decl.Type.Results.List {
				f.returns[i] = f.file.parseField(field)
			}
		}
	})
}

// Checks wether or not this function returns an error.
func (v Vars) HasError() bool {
	for _, v := range v {
		if v.Type().IsError() {
			return true
		}
	}

	return false
}

type ResolveResult struct {
	// TODO: add the list of packages needed to import those funcs maybe?
	ordered []*Func
	types   map[string]*Func
}

// Resolve the given types by finding which functions are needed to be called to
// actually instantiate them. It will recusrsively resolve the functions params to build
// up the total chain.
func (fns Funcs) Resolve(types ...*Type) (*ResolveResult, error) {
	r := &ResolveResult{
		types: make(map[string]*Func),
	}

	for _, typ := range types {
		_, found := r.types[typ.String()]

		// We already know how to resolve this type.
		if found {
			continue
		}

		for _, f := range fns {
			for _, ret := range f.Returns() {
				if ret.Type() == typ {
					resolved, err := r.resolveFn(fns, f)

					if err != nil {
						return nil, err
					}

					r.types[typ.String()] = resolved
					// FIXME: refacto needed to remove duplication with resolveFn
				}
			}
		}
	}

	return r, nil
}

func (r *ResolveResult) Funcs() []*Func { return r.ordered }

func (r *ResolveResult) resolveFn(fns Funcs, fn *Func) (*Func, error) {
	for _, p := range fn.Params() {
		_, found := r.types[p.Type().String()]

		if found {
			continue
		}

		for _, f := range fns {
			for _, ret := range f.Returns() {
				if ret.Type() == p.Type() {
					resolved, err := r.resolveFn(fns, f)

					if err != nil {
						return nil, err
					}

					r.types[p.Type().String()] = resolved
					found = true
				}
			}
		}

		if !found {
			return nil, fmt.Errorf("could not find a valid constructor for %s", p.Type().String())
		}
	}

	r.ordered = append(r.ordered, fn)

	return fn, nil
}
