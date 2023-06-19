package parser

import (
	"fmt"
	"strings"
	"sync"

	"github.com/YuukanOO/ease/pkg/crypto"
)

// Prefix length added to package when importing in the generation step to prevent
// collision for package with the same names (which is likely to happen in Go).
const aliasSuffixSize = 4

// Represents a single package and act as a registry of declarations for easy parsing.
type Package struct {
	lazy  sync.Once
	name  string
	path  string
	alias string
}

func newPackage(path string) *Package {
	return &Package{
		name: path[strings.LastIndex(path, "/")+1:],
		path: path,
	}
}

func (p *Package) Name() string { return p.name }
func (p *Package) Path() string { return p.path }

// Alias for the package. This is used to prevent collision with other packages
// when generating import declarations.
func (p *Package) Alias() string {
	p.lazy.Do(func() {
		p.alias = fmt.Sprintf("%s_%s", p.name, crypto.Prefix(p.path, aliasSuffixSize))
	})

	return p.alias
}
