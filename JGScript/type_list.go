package JGScript

import (
	"fmt"
)

// LIST
type List struct{ vals *[]*Type }

func (p *List) GetT() string {
	return "list"
}

func (p *List) str() string {
	s := "[ "
	max := len(*p.vals) - 1
	for i, e := range *p.vals {
		s += (*e).str()
		if i < max {
			s += ", "
		}
	}
	s += " ]"
	return s
}

func (p *List) GetIndex(a *AST, c *Context) *Result {
	res := MkRes()
    if len(a.C) == 0 {
        t := Type(p)
        return funCopy(&t, c)
    }
    atype := a.C[0].T
    if atype == "range" || atype == "range_start" || atype == "range_end" || atype == "range_empty" {
	    max := int64(len(*p.vals)) + 1
        f := int64(0)
        t := max - 1
        if atype == "range" || atype == "range_start" {
            from := res.reg(ExecAst(a.C[0].C[0], c))
            if res.has() {
                return res
            }
            if (*from).GetT() != "int" {
                return res.fail(MakeError("ListIndexError", "Invalid list index start type: " + (*from).GetT(), a.C[0].C[0].S, a.C[0].C[0].E, c))
            }
            f = (*from).(*Int).val
	        if f < 0 {
		        f += max
	        }
	        if f >= max {
		        return res.fail(MakeError("ListIndexError", "Index start "+(*from).str()+" out of range "+fmt.Sprint(max), a.C[0].C[0].S, a.C[0].C[0].E, c))
	        }
        }
        if atype == "range" || atype == "range_end" {
            index := 1
            if atype == "range_end" {
                index = 0
            }
            to := res.reg(ExecAst(a.C[0].C[index], c))
            if res.has() {
                return res
            }
            if (*to).GetT() != "int" {
                return res.fail(MakeError("ListIndexError", "Invalid list index end type: " + (*to).GetT(), a.C[0].C[index].S, a.C[0].C[index].E, c))
            }
            t = (*to).(*Int).val
	        if t < 0 {
		        t += max
	        }
	        if t >= max {
		        return res.fail(MakeError("ListIndexError", "Index end "+(*to).str()+" out of range "+fmt.Sprint(max), a.C[0].C[index].S, a.C[0].C[index].E, c))
	        }
        }
        types := &[]*Type{}
        if t > f {
            for i := f; i < t; i++ {
                *types = append(*types, (*p.vals)[i])
            }
        } else {
            for i := f; i > t; i-- {
                *types = append(*types, (*p.vals)[i])
            }
        }
        if len(a.C) > 1 {
            otypes := &[]*Type {}
            a.C = a.C[1:]
            for _, e := range *types {
                if (*e).GetT() != "list" {
		            return res.fail(MakeError("ListIndexError", "Non list type", a.C[1].S, a.C[1].E, c))
                }
                l := (*e).(*List)
                *otypes = append(*otypes, res.reg(l.GetIndex(a, c)))
            }
		    return res.suss(&List{otypes}, c)
        }
		return res.suss(&List{types}, c)
    }
	t := res.reg(ExecAst(a.C[0], c))
	if res.has() {
		return res
	}
	if (*t).GetT() == "int" {
	    it, _ := (*t).(*Int)
	    i := it.val
	    max := int64(len(*p.vals))
	    if i < 0 {
		    i += max
	    }
	    if i >= max {
		    return res.fail(MakeError("ListIndexError", "Index "+(*t).str()+" out of range "+fmt.Sprint(max), a.S, a.E, c))
	    }
        if len(a.C) > 1 {
            ot := (*p.vals)[i]
            if (*ot).GetT() == "list" {
                a.C = a.C[1:]
                return (*ot).(*List).GetIndex(a, c)
            }
        }
	    return res.sussP((*p.vals)[i], c)
    }
    return res.fail(MakeError("ListIndexError", "Invalid list index type: " + (*t).GetT(), a.C[0].S, a.C[0].E, c))
}

func (p *List) SetIndex(a, v *AST, c *Context) *Result {
    res := MkRes()
	val := res.reg(ExecAst(v, c))
	if res.has() {
		return res
	}
    return p.SetIndexVal(a, val, c)
}

func (p *List) SetIndexVal(a *AST, val *Type, c *Context) *Result {
	res := MkRes()
    atype := a.C[0].T
    if atype == "range" || atype == "range_start" || atype == "range_end" || atype == "range_empty" {
	    max := int64(len(*p.vals)) + 1
        f := int64(0)
        t := max - 1
        if atype == "range" || atype == "range_start" {
            from := res.reg(ExecAst(a.C[0].C[0], c))
            if res.has() {
                return res
            }
            if (*from).GetT() != "int" {
                return res.fail(MakeError("ListIndexError", "Invalid list index start type: " + (*from).GetT(), a.C[0].C[0].S, a.C[0].C[0].E, c))
            }
            f = (*from).(*Int).val
	        if f < 0 {
		        f += max
	        }
	        if f >= max {
		        return res.fail(MakeError("ListIndexError", "Index start "+(*from).str()+" out of range "+fmt.Sprint(max), a.C[0].C[0].S, a.C[0].C[0].E, c))
	        }
        }
        if atype == "range" || atype == "range_end" {
            index := 1
            if atype == "range_end" {
                index = 0
            }
            to := res.reg(ExecAst(a.C[0].C[index], c))
            if res.has() {
                return res
            }
            if (*to).GetT() != "int" {
                return res.fail(MakeError("ListIndexError", "Invalid list index end type: " + (*to).GetT(), a.C[0].C[index].S, a.C[0].C[index].E, c))
            }
            t = (*to).(*Int).val
	        if t < 0 {
		        t += max
	        }
	        if t >= max {
		        return res.fail(MakeError("ListIndexError", "Index end "+(*to).str()+" out of range "+fmt.Sprint(max), a.C[0].C[index].S, a.C[0].C[index].E, c))
	        }
        }
        if len(a.C) > 1 {
            types := &[]*Type{}
            if t > f {
                for i := f; i < t; i++ {
                    *types = append(*types, (*p.vals)[i])
                }
            } else {
                for i := f; i > t; i-- {
                    *types = append(*types, (*p.vals)[i])
                }
            }
            a.C = a.C[1:]
            if (*val).GetT() == "list" {
                val_list := *(*val).(*List).vals
                max := len(val_list)
                for i, e := range *types {
                    if (*e).GetT() != "list" {
		                return res.fail(MakeError("ListIndexError", "Non list type", a.C[1].S, a.C[1].E, c))
                    }
                    l := (*e).(*List)
                    if i < max {
                        l.SetIndexVal(a, val_list[i], c)
                    } else {
                        nt := Type(&Null{})
                        l.SetIndexVal(a, &nt, c)
                    }
                }
            } else {
                for _, e := range *types {
                    if (*e).GetT() != "list" {
		                return res.fail(MakeError("ListIndexError", "Non list type", a.C[1].S, a.C[1].E, c))
                    }
                    l := (*e).(*List)
                    nt := Type(&Null{})
                    l.SetIndexVal(a, &nt, c)
                }
            }
        } else {
            if (*val).GetT() == "list" {
                val_list := *(*val).(*List).vals
                max := len(val_list)
                index := 0
                if t > f {
                    for i := f; i < t; i++ {
                        if index < max {
                            (*p.vals)[i] = val_list[index]
                            index++
                        } else {
                            nt := Type(&Null{})
                            (*p.vals)[i] = &nt
                        }
                    }
                } else {
                    for i := f; i > t; i-- {
                        if index < max {
                            (*p.vals)[i] = val_list[index]
                            index++
                        } else {
                            nt := Type(&Null{})
                            (*p.vals)[i] = &nt
                        }
                    }
                }
            } else {
                if t > f {
                    for i := f; i < t; i++ {
                        (*p.vals)[i] = val
                    }
                } else {
                    for i := f; i > t; i-- {
                        (*p.vals)[i] = val
                    }
                }
            }
        }
		return res.suss(p, c)
    }
	t := res.reg(ExecAst(a.C[0], c))
	if res.has() {
		return res
	}
	if (*t).GetT() != "int" {
		return res.fail(MakeRTError("List index most be integer", a.C[0].S, a.C[0].E, c))
	}
	it, _ := (*t).(*Int)
	i := it.val
	max := int64(len(*p.vals))
	if i < 0 {
		i += max
	}
	if i >= max {
		return res.fail(MakeRTError("Index "+(*t).str()+" out of range "+fmt.Sprint(max), a.S, a.E, c))
	}
	(*p.vals)[i] = funCopyTypes(val, c)
	return res.suss(p, c)
}

func (p *List) operator(op string, o *Type, c *Context, s, e *BPos) *Result {
	res := MkRes()
	if op == "_add" && (*o).GetT() == "string" {
		return res.suss(&Str{p.str() + (*o).str()}, c)
	}
    if op == "_mul" && (*o).GetT() == "int" {
        max := (*o).(*Int).val
        l := &[]*Type {}
        var val *Type = nil
        if len(*p.vals) == 1 {
            val = (*p.vals)[0]
        } else {
            t := Type(p)
            val = &t
        }
        if val == nil {
            t := Type(&Null{})
            val = &t
        }
        for i := int64(0); i < max; i++ {
            *l = append(*l, funCopy(val, c).t)
        }
        return res.suss(&List{l}, c)
    }
	if op == "_left_arrow" {
		*p.vals = append(*p.vals, o)
		return res.sussP(o, c)
	}
	if op == "_eq" {
		return res.suss(&Bool{p.equals(o)}, c)
	}
	if op == "_not_eq" {
		return res.suss(&Bool{!p.equals(o)}, c)
	}
	return res.fail(MakeRTError("Invalid operation for "+(*p).str()+" and "+(*o).str(), s, e, c))
}

func (p *List) getBool() bool {
	return true
}

func (p *List) getByte() byte {
    return 0
}

func (p *List) getInt() int64 {
	return 0
}

func (p *List) getFloat() float64 {
	return 0
}

func (p *List) GetList() *[]*Type {
    return p.vals
}

func (p *List) equals(t *Type) bool {
	if (*t).GetT() == "List" {
		i, _ := (*t).(*List)
		if len(*p.vals) != len(*i.vals) {
			return false
		}
		for k, v := range *p.vals {
			if !(*v).equals((*i.vals)[k]) {
				return false
			}
		}
		return true
	}
	return false
}

func (p *List) GetMember(a *AST, c *Context) *Result {
    res := MkRes()
    switch a.V {
    case "len":
        return res.suss(&Int{int64(len(*p.vals))}, c)
	case "get":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
			    ia := arg.C[0]
				lis := ((*t).(*List).vals)
				i := res.reg(ExecAst(ia, c))
				if res.has() {
					return res
				}
				if (*i).GetT() != "int" {
					return res.fail(MakeRTError("Index must be int.", ia.S, ia.E, c))
				}
				it, _ := (*i).(*Int)
				ind := it.val
				if ind < 0 {
					ind += int64(len(*lis))
				}
				if ind >= int64(len(*lis)) {
					return res.fail(MakeRTError("Index out of range.", ia.S, ia.E, c))
				}
				return res.sussP((*lis)[ind], c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "set":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
			    ia := arg.C[0]
			    va := arg.C[1]
				lis := ((*t).(*List).vals)
				i := res.reg(ExecAst(ia, c))
				if res.has() {
					return res
				}
				if (*i).GetT() != "int" {
					return res.fail(MakeRTError("Index must be int.", ia.S, ia.E, c))
				}
				it, _ := (*i).(*Int)
				ind := it.val
				if ind < 0 {
					ind += int64(len(*lis))
				}
				if ind >= int64(len(*lis)) {
					return res.fail(MakeRTError("Index out of range.", ia.S, ia.E, c))
				}
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
                }
                (*lis)[ind] = v
				return res.sussP(t, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
    case "push":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
                for _, e := range arg.C {
				    v := res.reg(ExecAst(e, c))
				    if res.has() {
					    return res
				    }
				    *lis = append(*lis, funCopyTypes(v, c))
                }
				return res.sussP(t, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "extend":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
                for _, va := range arg.C {
				    v := res.reg(ExecAst(va, c))
				    if res.has() {
					    return res
				    }
				    if (*v).GetT() != "list" {
			            return res.fail(MakeRTError("Invalid function arguments.", va.S, va.E, c))
                    }
					l2 := (*v).(*List).vals
					for _, e := range *l2 {
						*lis = append(*lis, funCopyTypes(e, c))
					}
                }
				return res.sussP(t, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "insert":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
			    ia := arg.C[0]
			    va := arg.C[1]
				lis := ((*t).(*List).vals)
				i := res.reg(ExecAst(ia, c))
				if res.has() {
					return res
				}
				if (*i).GetT() != "int" {
					return res.fail(MakeRTError("Index must be int.", ia.S, ia.E, c))
				}
				it, _ := (*i).(*Int)
				ind := it.val
				if ind < 0 {
					ind += int64(len(*lis))
				}
				if ind >= int64(len(*lis)) {
					return res.fail(MakeRTError("Index out of range.", ia.S, ia.E, c))
				}
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				*lis = append((*lis)[:ind], append([]*Type{v}, (*lis)[ind:]...)...)
				return res.sussP(t, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
    case "pop":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
				max := int64(len(*lis))
                ind := max - 1
                if len(arg.C) > 0 {
                    ia := arg.C[0]
				    i := res.reg(ExecAst(ia, c))
				    if res.has() {
					    return res
				    }
				    if (*i).GetT() != "int" {
					    return res.fail(MakeRTError("Index must be int.", ia.S, ia.E, c))
				    }
				    it, _ := (*i).(*Int)
				    ind = it.val
				    if ind < 0 {
					    ind += max
				    }
				    if ind >= max {
					    return res.fail(MakeRTError("Index out of range.", ia.S, ia.E, c))
				    }
                }
				ret := (*lis)[ind]
				*lis = append((*lis)[:ind], (*lis)[ind+1:]...)
				return res.sussP(ret, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
    case "remove":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
                va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				ind := -1
				for i, val := range *lis {
					if (*v).equals(val) {
						ind = i
						break
					}
				}
                var ret *Type
				if ind != -1 {
				    ret = (*lis)[ind]
					*lis = append((*lis)[:ind], (*lis)[ind+1:]...)
				} else {
                    rt := Type(&Null{})
                    ret = &rt
                }
				return res.sussP(ret, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
    case "clear":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if (*t).GetT() == "list" {
				lis := (*t).(*List)
                lis.vals = &[]*Type{}
				return res.suss(lis, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "has":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				for _, val := range *lis {
					if (*v).equals(val) {
						return res.suss(&Bool{true}, c)
					}
				}
				return res.suss(&Bool{false}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "find":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				for i, val := range *lis {
					if (*v).equals(val) {
						return res.suss(&Int{int64(i)}, c)
					}
				}
				return res.suss(&Int{-1}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "find_last":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
                for i := len(*lis) - 1; i >= 0; i-- {
					if (*v).equals((*lis)[i]) {
						return res.suss(&Int{int64(i)}, c)
					}
				}
				return res.suss(&Int{-1}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "find_all":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
                l := []*Type{}
				for i, val := range *lis {
					if (*v).equals(val) {
                        it := Type(&Int{int64(i)})
                        l = append(l, &it)
					}
				}
				return res.suss(&List{&l}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "replace":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
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
				for i, val := range *lis {
					if (*v).equals(val) {
						(*lis)[i] = n
                        break
					}
				}
				return res.sussP(t, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "replace_last":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
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
                for i := len(*lis) - 1; i >= 0; i-- {
					if (*v).equals((*lis)[i]) {
						(*lis)[i] = n
                        break
					}
				}
				return res.sussP(t, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "replace_all":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
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
				for i, val := range *lis {
					if (*v).equals(val) {
						(*lis)[i] = n
					}
				}
				return res.sussP(t, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
	case "count":
        return res.suss(&NFun{func(t *Type, c *Context, arg *AST) *Result {
			if len(arg.C) < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
			}
			if (*t).GetT() == "list" {
				lis := ((*t).(*List).vals)
			    va := arg.C[0]
				v := res.reg(ExecAst(va, c))
				if res.has() {
					return res
				}
				count := int64(0)
				for _, val := range *lis {
					if (*v).equals(val) {
						count++
					}
				}
				return res.suss(&Int{count}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", arg.S, arg.E, c))
        }}, c)
    }
    return res.fail(MakeError("MemberAccessError", "Type " + p.GetT() + " has not member " + a.V, a.S, a.E, c))
}
