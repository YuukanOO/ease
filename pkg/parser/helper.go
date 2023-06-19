package parser

import "go/types"

// Checks if the given typename is a builtin one.
func IsBuiltin(typeName string) bool {
	return types.Universe.Lookup(typeName) != nil
}
