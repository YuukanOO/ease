package parser

import (
	"strings"
)

// Represents a single package and act as a registry of declarations for easy parsing.
type Package struct {
	name string
	path string
}

func newPackage(path string) *Package {
	return &Package{
		name: path[strings.LastIndex(path, "/")+1:],
		path: path,
	}
}

func (p *Package) Name() string { return p.name }
func (p *Package) Path() string { return p.path }
