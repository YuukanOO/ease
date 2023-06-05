package parser

import (
	"go/ast"
	"strings"
)

type Package struct {
	*ast.Package

	Name string
	Path string
	Dir  string
}

func newPackage(mod *Module, dir string, pkg *ast.Package) *Package {
	return &Package{
		Package: pkg,
		Name:    pkg.Name,
		Dir:     dir,
		Path:    strings.ReplaceAll(strings.ReplaceAll(dir, mod.Dir, mod.Path), "\\", "/"), // For the path, we always need forward slashes
	}
}
