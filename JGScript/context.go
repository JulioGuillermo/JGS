package JGScript

type Context struct {
	parent *Context
	st     *map[string]*Type
	s      *BPos
	n      string
}

func MakeContext(p *Context, s *BPos) *Context {
	st := &map[string]*Type{}
	return &Context{parent: p, st: st, s: s, n: ""}
}

func (p *Context) GetST() *map[string] *Type {
    return p.st
}

func (p *Context) set(n string, v *Type) {
	if n != "_" {
		if n[0] == '*' {
			n = n[1:]
			if p.parent == nil {
				p.set(n, v)
			} else {
				p.parent.set(n, v)
			}
		}
		(*p.st)[n] = funCopyTypes(v, p)
	}
}

func (p *Context) Set(n string, v *Type) {
    p.set(n, v)
}

func (p *Context) get(n string) *Type {
	if n[0] == '*' {
		n = n[1:]
		if p.parent == nil {
			return p.get(n)
		} else {
			return p.parent.get(n)
		}
	}
	var t *Type = (*p.st)[n]
	if t == nil && p.parent != nil {
		return p.parent.get(n)
	}
	return t
}

func (p *Context) Get(n string) *Type {
    return p.get(n)
}

func (p *Context) remove(n string) *Type {
	if n[0] == '*' {
		n = n[1:]
		if p.parent == nil {
			return p.remove(n)
		} else {
			return p.parent.remove(n)
		}
	}
	t := (*p.st)[n]
	if t == nil && p.parent != nil {
		return p.parent.remove(n)
	}
	delete(*p.st, n)
	return t
}

func (p *Context) has(n string) bool {
	if n[0] == '*' {
		n = n[1:]
		if p.parent == nil {
			return p.has(n)
		} else {
			return p.parent.has(n)
		}
	}
	return (*p.st)[n] != nil
}

func (p *Context) exist(n string) bool {
	if n[0] == '*' {
		n = n[1:]
		if p.parent == nil {
			return p.exist(n)
		} else {
			return p.parent.exist(n)
		}
	}
	if (*p.st)[n] != nil {
		return true
	}
	if p.parent != nil {
		return p.parent.exist(n)
	}
	return false
}
