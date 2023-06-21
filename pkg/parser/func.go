package parser

import (
	"go/ast"
	"sync"
)

type (
	Vars []*Var

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
