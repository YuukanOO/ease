package parser

import "github.com/YuukanOO/ease/pkg/flag"

type VarKind uint

const VarKindUnknown VarKind = 0

const (
	VarKindIdent VarKind = 1 << iota
	VarKindBuiltin
	VarKindPointer
	VarKindSlice
	VarKindSliceOfPointer
	VarKindMap
)

type Var struct {
	*Decl
	kind       VarKind
	underlying *Type
}

func (v *Var) Type() *Type     { return v.underlying }
func (v *Var) IsPointer() bool { return flag.IsSet(v.kind, VarKindPointer) }
