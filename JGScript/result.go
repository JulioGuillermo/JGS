package JGScript

type Result struct {
	t  *Type
	e  *Error
	c  *Context
	b  bool
	r  bool
	n  string
	se bool
}

func MkRes() *Result {
	return &Result{nil, nil, nil, false, false, "", false}
}

func (p *Result) reg(r *Result) *Type {
	if r.e != nil {
		if p.e == nil {
			p.e = r.e
			p.se = r.se
		}
	} else {
		p.t = r.t
		p.c = r.c
		p.b = r.b
		p.n = r.n
		p.r = r.r
	}
	return p.t
}

func (p *Result) Reg(r *Result) *Type {
    return p.reg(r)
}

func (p *Result) suss(t Type, c *Context) *Result {
	p.t = &t
	p.c = c
	return p
}

func (p *Result) Suss(t Type, c *Context) *Result {
    return p.suss(t, c)
}

func (p *Result) sussP(t *Type, c *Context) *Result {
	p.t = t
	p.c = c
	return p
}

func (p *Result) SussP(t *Type, c *Context) *Result {
    return p.sussP(t, c)
}

func (p *Result) fail(e *Error) *Result {
	p.e = e
	return p
}

func (p *Result) Fail(e *Error) *Result {
    return p.fail(e)
}

func (p *Result) has() bool {
	return p.e != nil
}

func (p *Result) Has() bool {
    return p.has()
}
