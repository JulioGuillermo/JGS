package JGScript

// BOOL
type Bool struct{ val bool }

func (p *Bool) GetT() string {
	return "bool"
}

func (p *Bool) str() string {
	if p.val {
		return "true"
	}
	return "false"
}

func (p *Bool) operator(op string, o *Type, c *Context, s, e *BPos) *Result {
	res := MkRes()
	if op == "_add" && (*o).GetT() == "string" {
		return res.suss(&Str{p.str() + (*o).str()}, c)
	}
	if op == "_not" {
		return res.suss(&Bool{!p.val}, c)
	}
	if (*o).GetT() == "bool" {
		b := (*o).(*Bool).val
		switch op {
		case "_and":
			return res.suss(&Bool{p.val && b}, c)
		case "_or":
			return res.suss(&Bool{p.val || b}, c)
		case "_eq":
			return res.suss(&Bool{p.val == b}, c)
		case "_not_eq":
			return res.suss(&Bool{p.val != b}, c)
		}
	} else {
		switch op {
		case "_eq":
			return res.suss(&Bool{false}, c)
		case "_not_eq":
			return res.suss(&Bool{true}, c)
		}
	}
	return res.fail(MakeRTError("Invalid operation for "+(*p).str()+" and "+(*o).str(), s, e, c))
}

func (p *Bool) getBool() bool {
	return p.val
}

func (p *Bool) GetBool() *bool {
	return &p.val
}

func (p *Bool) getByte() byte {
    if p.val {
        return 1
    }
    return 0
}

func (p *Bool) getInt() int64 {
	if p.val {
		return 1
	}
	return 0
}

func (p *Bool) getFloat() float64 {
	if p.val {
		return 1
	}
	return 0
}

func (p *Bool) equals(t *Type) bool {
	if (*t).GetT() == "bool" {
		i, _ := (*t).(*Bool)
		return p.val == i.val
	}
	return false
}

func (p *Bool) GetMember(a *AST, c *Context) *Result {
    return MkRes().fail(MakeError("MemberAccessError", "Type " + p.GetT() + " has not member " + a.V, a.S, a.E, c))
}
