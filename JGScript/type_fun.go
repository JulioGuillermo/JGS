package JGScript

// FUN
type Fun struct {
	args, body *AST
	async      bool
}

func (p *Fun) GetT() string {
	return "fun"
}

func (p *Fun) str() string {
	return "function"
}

func (p *Fun) ExecFun(args *AST, name string, c, ac *Context) *Result {
	res := MkRes()
	ctx := MakeContext(c, args.S)
	ctx.n = name
	res.reg(execFunArgs(p.args, args, ctx, ac))
	if res.has() {
		return res
	}
	if p.async {
		go ExecAst(p.body, ctx)
		return res.suss(&Null{}, c)
	}
	ret := res.reg(ExecAst(p.body, ctx))
	if res.has() {
		return res
	}
	return res.sussP(ret, c)
}

func (p *Fun) operator(op string, o *Type, c *Context, s, e *BPos) *Result {
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

func (p *Fun) getBool() bool {
	return true
}

func (p *Fun) getByte() byte {
    return 0
}

func (p *Fun) getInt() int64 {
	return 0
}

func (p *Fun) getFloat() float64 {
	return 0
}

func (p *Fun) equals(t *Type) bool {
	if (*t).GetT() == "fun" {
		i, _ := (*t).(*Fun)
		return p == i
	}
	return false
}

func (p *Fun) GetMember(a *AST, c *Context) *Result {
    return MkRes().fail(MakeError("MemberAccessError", "Type " + p.GetT() + " has not member " + a.V, a.S, a.E, c))
}
