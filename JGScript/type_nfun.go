package JGScript

// NativeFuns
type NFun struct { F func(t *Type, c *Context, args *AST) *Result }

func (p *NFun) GetT() string {
    return "nfun"
}

func (p *NFun) str() string {
    return "native_function"
}

func (p *NFun) getBool() bool {
    return true
}

func (p *NFun) getByte() byte {
    return 0
}

func (p *NFun) getInt() int64 {
    return 0
}

func (p *NFun) getFloat() float64 {
    return 0
}

func (p *NFun) equals(t *Type) bool {
    if (*t).GetT() == "nfun" {
        i, _ := (*t).(*NFun)
        return p == i
    }
    return false
}

func (p *NFun) operator(op string, o *Type, c *Context, s, e *BPos) *Result {
	res := MkRes()
	if op == "_add" && (*o).GetT() == "string" {
		return res.suss(&Str{p.str() + (*o).str()}, c)
	}
	if op == "_eq" {
		return res.suss(&Bool{p.equals(o)}, c)
	}
	if op == "_not_eq" {
		return res.suss(&Bool{!p.equals(o)}, c)
	}
	return res.fail(MakeRTError("Invalid operation for "+(*p).str()+" and "+(*o).str(), s, e, c))
}

func (p *NFun) ExecFun(t *Type, args *AST, c *Context) *Result {
	return p.F(t, c, args)
}

func (p *NFun) GetMember(a *AST, c *Context) *Result {
    return MkRes().fail(MakeError("MemberAccessError", "Type " + p.GetT() + " has not member " + a.V, a.S, a.E, c))
}
