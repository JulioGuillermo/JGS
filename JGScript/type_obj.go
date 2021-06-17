package JGScript

// OBJ
type Obj struct {
	vals  *Context
	types []string
}

func (p *Obj) GetT() string {
	return "obj"
}

func (p *Obj) str() string {
    str := (*p.vals.st)["_to_string"]
    if str != nil {
        if (*str).GetT() == "fun" {
            fun_str := (*str).(*Fun)
            res := fun_str.ExecFun(&AST{C:[]*AST{}}, "_to_string", p.vals, p.vals)
            if !res.has() {
                return (*res.t).str()
            }
        } else {
            return (*str).str()
        }
    }
	s := "{ "
	for i, e := range *p.vals.st {
		if i != "this" {
			s += i + ": " + (*e).str() + ", "
		}
	}
	s += "}"
	return s
}

func (p *Obj) GetMember(a *AST, c *Context) *Result {
	res := MkRes()
    m := (*p.vals.st)[a.V]
    if m != nil {
	    return res.sussP(m, p.vals)
    }
    switch a.V {
	case "get_member":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "obj" {
				lis := (*t).(*Obj).vals
			    ia := arg.C[0]
				i := res.reg(ExecAst(ia, c))
				if res.has() {
					return res
				}
				key := (*i).str()
                ret := (*lis.st)[key]
                if ret == nil {
                    tn := Type(&Null{})
                    ret = &tn
                }
				return res.sussP(ret, lis)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "set_member":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, c))
			}
			if (*t).GetT() == "obj" {
				lis := (*t).(*Obj).vals
			    ia := arg.C[0]
			    va := arg.C[1]
				i := res.reg(ExecAst(ia, c))
				if res.has() {
					return res
				}
				key := (*i).str()
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				(*lis.st)[key] = v
				return res.sussP(v, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "remove_member":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "obj" {
				lis := (*t).(*Obj).vals
			    ia := arg.C[0]
				i := res.reg(ExecAst(ia, c))
				if res.has() {
					return res
				}
				key := (*i).str()
				ret := (*lis.st)[key]
				delete(*lis.st, key)
				return res.sussP(ret, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "has_member":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "obj" {
				lis := (*t).(*Obj).vals
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				if (*v).GetT() == "string" {
					attr := (*v).(*Str)
					return res.suss(&Bool{(*lis.st)[attr.val] != nil}, c)
				}
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "is_object_type":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "obj" {
				lis := (*t).(*Obj).types
			    ta := arg.C[0]
				ot := res.reg(ExecAst(ta, c))
				if res.has() {
					return res
				}
				if (*ot).GetT() != "string" {
					return res.fail(MakeRTError("Type argument must be string.", arg.S, arg.E, c))
				}
				ts := (*ot).(*Str).val
				if ts == "obj" {
					return res.suss(&Bool{true}, c)
				}
				for _, val := range lis {
					if ts == val {
						return res.suss(&Bool{true}, c)
					}
				}
			}
			return res.suss(&Bool{false}, c)
        }}, c)
	case "get_object_type":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if (*t).GetT() == "obj" {
				lis := (*t).(*Obj).types
				ot := "obj"
				if len(lis) > 0 {
					ot = (lis)[0]
				}
				return res.suss(&Str{ot}, c)
			}
			return res.suss(&Str{(*t).GetT()}, c)
        }}, c)
    }
    t := "obj"
    if len(p.types) > 0 {
        t = p.types[0]
    }
    return res.fail(MakeError("ObjMemberError", "Object type " + t + " has not member " + a.V, a.S, a.E, c))
}

func (p *Obj) SetMemb(a, v *AST, c *Context) *Result {
	res := MkRes()
	t := res.reg(execAssign(&AST{C: []*AST{a, v}}, p.vals, c))
	if res.has() {
		return res
	}
	return res.sussP(t, c)
}

func (p *Obj) objType(t string) bool {
	for _, s := range p.types {
		if s == t {
			return true
		}
	}
	return false
}

func (p *Obj) objOperator(n string, args *AST, c *Context) *Result {
	res := MkRes()
	t := p.vals.get(n)
	if t == nil {
		return nil
	}
	if (*t).GetT() != "fun" {
		return res.fail(MakeRTError(n+" is not a function", args.S, args.E, c))
	}
	f, _ := (*t).(*Fun)
	ret := res.reg(f.ExecFun(args, n, p.vals, c))
	if res.has() {
		return res
	}
	return res.sussP(ret, c)
}

func (p *Obj) operator(op string, o *Type, c *Context, s, e *BPos) *Result {
	res := MkRes()
	if op == "_add" && (*o).GetT() == "string" {
		return res.suss(&Str{p.str() + (*o).str()}, c)
	}
	if op == "_right_arrow" {
		if (*o).GetT() == "obj" {
			obj := (*o).(*Obj)
			for _, t := range obj.types {
				if !p.objType(t) {
					p.types = append([]string{t}, p.types...)
				}
			}
			for k, v := range *obj.vals.st {
				if k != "this" {
					(*p.vals.st)[k] = v
				}
			}
		} else if (*o).GetT() == "string" {
			if !p.objType((*o).str()) {
				p.types = append([]string{(*o).str()}, p.types...)
			}
		}
		return res.suss(p, c)
	}
	if op == "_eq" {
		return res.suss(&Bool{p.equals(o)}, c)
	}
	if op == "_not_eq" {
		return res.suss(&Bool{!p.equals(o)}, c)
	}
	ot := "obj"
	if len(p.types) > 0 {
		ot = p.types[0]
	}
	return res.fail(MakeRTError("Undefined operation for object type "+ot, s, e, c))
}

func (p *Obj) getBool() bool {
	return true
}

func (p *Obj) getInt() int64 {
	return 0
}

func (p *Obj) getByte() byte {
    return 0
}

func (p *Obj) getFloat() float64 {
	return 0
}

func (p *Obj) GetObj() *Context {
    return p.vals
}

func (p *Obj) GetMap() *map[string] *Type {
    return p.vals.st
}

func (p *Obj) equals(t *Type) bool {
	if (*t).GetT() == "obj" {
		i, _ := (*t).(*Obj)
		return p == i
	}
	return false
}
