package parser

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
	name       string
	kind       VarKind
	underlying *Type
}

func (v *Var) Name() string { return v.name }
func (v *Var) Type() *Type  { return v.underlying }
