package JGScript

import (
	"fmt"
	"strconv"
    "strings"
)

// STRING
type Str struct{ val string }

func (p *Str) GetT() string {
	return "string"
}

func (p *Str) str() string {
	return p.val
}

func (p *Str) operator(op string, o *Type, c *Context, s, e *BPos) *Result {
	res := MkRes()
	if op == "_add" {
		return res.suss(&Str{p.val + (*o).str()}, c)
	}
	if op == "_add_eq" {
		p.val += (*o).str()
		return res.suss(p, c)
	}
	switch op {
	case "_eq":
		if (*o).GetT() != "string" {
			res.suss(&Bool{false}, c)
		}
		return res.suss(&Bool{p.val == (*o).str()}, c)
	case "_not_eq":
		if (*o).GetT() != "string" {
			res.suss(&Bool{true}, c)
		}
		return res.suss(&Bool{p.val != (*o).str()}, c)
	case "_lt":
		if (*o).GetT() == "string" {
			return res.suss(&Bool{p.val < (*o).str()}, c)
		}
	case "_lte":
		if (*o).GetT() == "string" {
			return res.suss(&Bool{p.val <= (*o).str()}, c)
		}
	case "_gt":
		if (*o).GetT() == "string" {
			return res.suss(&Bool{p.val > (*o).str()}, c)
		}
	case "_gte":
		if (*o).GetT() == "string" {
			return res.suss(&Bool{p.val >= (*o).str()}, c)
		}
	}
	return res.fail(MakeRTError("Invalid operation for "+(*p).str()+" and "+(*o).str(), s, e, c))
}

func (p *Str) getBool() bool {
	return true
}

func (p *Str) getByte() byte {
    n, e := strconv.ParseInt(p.val, 10, 64)
    if e != nil {
        return 0
    }
    return byte(n)
}

func (p *Str) getInt() int64 {
    base := 10
    str := p.val
    if strings.HasPrefix(p.val, "0x") {
        base = 16
        str = p.val[2:]
    } else if strings.HasPrefix(p.val, "#") {
        base = 16
        str = p.val[1:]
    } else if strings.HasPrefix(p.val, "@") {
        base = 8
        str = p.val[1:]
    } else if strings.HasPrefix(p.val, "$") {
        base = 2
        str = p.val[1:]
    }
	n, e := strconv.ParseInt(str, base, 64)
	if e != nil {
		return 0
	}
	return n
}

func (p *Str) getFloat() float64 {
	n, e := strconv.ParseFloat(p.val, 10)
	if e != nil {
		return 0
	}
	return n
}

func (p *Str) GetString() *string {
    return &p.val
}

func (p *Str) equals(t *Type) bool {
	if (*t).GetT() == "string" {
		i, _ := (*t).(*Str)
		return p.val == i.val
	}
	return false
}

func (p *Str) GetIndex(a *AST, c *Context) *Result {
	res := MkRes()
	t := res.reg(ExecAst(a.C[0], c))
	if res.has() {
		return res
	}
	if (*t).GetT() != "int" {
		return res.fail(MakeRTError("Index most be integer", a.C[0].S, a.C[0].E, c))
	}
	it, _ := (*t).(*Int)
	i := it.val
	max := int64(len(p.val))
	if i < 0 {
		i += max
	}
	if i > max {
		return res.fail(MakeRTError("Index "+(*t).str()+" out of range "+fmt.Sprint(max), a.S, a.E, c))
	}
	return res.suss(&Str{string(p.val[i])}, c)
}

func (p *Str) SetIndex(a, v *AST, c *Context) *Result {
	res := MkRes()
	t := res.reg(ExecAst(a.C[0], c))
	if res.has() {
		return res
	}
	if (*t).GetT() != "int" {
		return res.fail(MakeRTError("Index most be integer", a.C[0].S, a.C[0].E, c))
	}
	it, _ := (*t).(*Int)
	i := it.val
	max := int64(len(p.val))
	if i < 0 {
		i += max
	}
	if i > max {
		return res.fail(MakeRTError("Index "+(*t).str()+" out of range "+fmt.Sprint(max), a.S, a.E, c))
	}
	val := res.reg(ExecAst(v, c))
	if res.has() {
		return res
	}
	str := p.val[:i]
	if (*val).GetT() == "string" {
		ostr := (*val).(*Str)
		str += ostr.val
	} else {
		str += string(byte((*val).getInt()))
	}
	str += p.val[i+1:]
	p.val = str
	pt := Type(p)
	return res.sussP(&pt, c)
}

func (p *Str) GetMember(a *AST, c *Context) *Result {
    res := MkRes()
    switch a.V {
    case "len":
        return res.suss(&Int{int64(len(p.val))}, c)
	case "has":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "string" {
				lis := &((*t).(*Str).val)
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				return res.suss(&Bool{strings.Contains(*lis, (*v).str())}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "find":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "string" {
				lis := &((*t).(*Str).val)
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				i := strings.Index(*lis, (*v).str())
				return res.suss(&Int{int64(i)}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "find_last":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "string" {
				lis := &((*t).(*Str).val)
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				i := strings.LastIndex(*lis, (*v).str())
				return res.suss(&Int{int64(i)}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "count":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "string" {
				lis := (*t).(*Str).val
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				i := strings.Count(lis, (*v).str())
				return res.suss(&Int{int64(i)}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "replace":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, c))
			}
			if (*t).GetT() == "string" {
				lis := &((*t).(*Str).val)
			    va := arg.C[0]
			    na := arg.C[1]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				n := res.reg(ExecAst(na, c))
				if res.has() {
					return res
				}
				i := strings.Replace(*lis, (*v).str(), (*n).str(), 1)
				return res.suss(&Str{i}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "replace_all":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, c))
			}
			if (*t).GetT() == "string" {
				lis := &((*t).(*Str).val)
			    va := arg.C[0]
			    na := arg.C[1]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				n := res.reg(ExecAst(na, c))
				if res.has() {
					return res
				}
				i := strings.ReplaceAll(*lis, (*v).str(), (*n).str())
				return res.suss(&Str{i}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "trim_space":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if (*t).GetT() == "string" {
				return res.suss(&Str{strings.TrimSpace((*t).(*Str).val)}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "upper":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if (*t).GetT() == "string" {
				return res.suss(&Str{strings.ToUpper((*t).(*Str).val)}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "lower":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if (*t).GetT() == "string" {
				return res.suss(&Str{strings.ToLower((*t).(*Str).val)}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "split":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "string" {
				lis := (*t).(*Str).val
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				i := strings.Split(lis, (*v).str())
				st := &[]*Type{}
				for _, v := range i {
					ts := Type(&Str{v})
					*st = append(*st, &ts)
				}
				return res.suss(&List{st}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "fields":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if (*t).GetT() == "string" {
				lis := (*t).(*Str).val
				i := strings.Fields(lis)
				st := &[]*Type{}
				for _, v := range i {
					ts := Type(&Str{v})
					*st = append(*st, &ts)
				}
				return res.suss(&List{st}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
    }
    return res.fail(MakeError("MemberAccessError", "Type " + p.GetT() + " has not member " + a.V, a.S, a.E, c))
}
