package JGScript

// NULL
type Null struct{}

func (p *Null) GetT() string {
	return "null"
}

func (p *Null) str() string {
	return "null"
}

func (p *Null) operator(op string, o *Type, c *Context, s, e *BPos) *Result {
	res := MkRes()
	if op == "_add" && (*o).GetT() == "string" {
		return res.suss(&Str{p.str() + (*o).str()}, c)
	}
	switch op {
	case "_eq":
		return res.suss(&Bool{(*o).GetT() == "null"}, c)
	case "_not_eq":
		return res.suss(&Bool{(*o).GetT() != "null"}, c)
	}
	return res.fail(MakeRTError("Invalid operation for "+(*p).str()+" and "+(*o).str(), s, e, c))
}

func (p *Null) getBool() bool {
	return false
}

func (p *Null) getByte() byte {
    return 0
}

func (p *Null) getInt() int64 {
	return 0
}

func (p *Null) getFloat() float64 {
	return 0
}

func (p *Null) equals(t *Type) bool {
	return (*t).GetT() == "null"
}

func (p *Null) GetMember(a *AST, c *Context) *Result {
    return MkRes().fail(MakeError("MemberAccessError", "Type " + p.GetT() + " has not member " + a.V, a.S, a.E, c))
}
