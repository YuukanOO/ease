package parser

type Func struct {
	pkg     *Package
	name    string
	recv    *Var
	params  []*Var
	returns []*Var
}

func (f *Func) Params() []*Var  { return f.params }
func (f *Func) Returns() []*Var { return f.returns }
