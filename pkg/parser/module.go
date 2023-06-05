package parser

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

const (
	goModFilename = "go.mod"
	goModMaxDepth = 10
)

var ErrMaxDepthReached = fmt.Errorf("max depth of %d reached when looking for a go.mod file", goModMaxDepth)

// Represents a single go module.
type Module struct {
	Path string // Path of the module, used to construct subpackages path
	Dir  string // Physical directory at which the go.mod file was found
}

func (m *Module) CreatePackage(dir string, pkg *ast.Package) *Package {
	return newPackage(m, dir, pkg)
}

// Try to find a go.mod file in the given directory or any of its parents.
func findModule(dir string, currentDepth int) (*Module, error) {
	if currentDepth > goModMaxDepth {
		return nil, ErrMaxDepthReached
	}

	modpath := filepath.Join(dir, goModFilename)

	_, err := os.Stat(modpath)

	if os.IsNotExist(err) {
		return findModule(filepath.Dir(dir), currentDepth+1)
	}

	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(modpath)

	if err != nil {
		return nil, err
	}

	return &Module{
		Path: modfile.ModulePath(data),
		Dir:  dir,
	}, nil
}
