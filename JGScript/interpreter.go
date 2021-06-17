package JGScript

import (
	"strconv"
	"strings"
)

func AST_Type(t string, types []string) bool {
    for _, a := range types {
        if a == t {
            return true
        }
    }
    return false
}

func ExecAst(a *AST, c *Context) *Result {
	res := MkRes()
	switch a.T {
    case "import":
        return import_file(a, c)
	case "body":
		name := Type(&Str{a.V})
		c.set("__name__", &name)
		null := Type(&Null{})
		ret := &null
        max_body := len(a.C) - 1
		for i, l := range a.C {
			if l.T == "break" {
				res.b = true
                res.r = false
				return res.suss(&Null{}, c)
			}
			if l.T == "ret" {
				t := res.reg(ExecAst(l.C[0], c))
				if res.has() {
					return res
				}
				res.b = true
				res.r = true
				return res.sussP(t, c)
			}
			if l.T == "body_ret" {
				t := res.reg(ExecAst(l.C[0], c))
				if res.has() {
					return res
				}
				ret = t
			} else {
				res.reg(ExecAst(l, c))
				if res.has() {
					return res
				}
				if res.b || res.r {
					//res.b = false
					break
				}
                if i == max_body && res.t != nil && !AST_Type(l.T, []string{
                    "if",
                    "for",
                    "switch",
                    "try" }) {
                    ret = res.t
                }
				res.t = nil
			}
		}
		return res.sussP(ret, c)
	case "null":
		return res.suss(&Null{}, c)
	case "hex":
		n, _ := strconv.ParseInt(a.V, 16, 64)
        if len(a.V) <= 2 {
            return res.suss(&Byte{byte(n)}, c)
        }
		return res.suss(&Int{n}, c)
	case "oct":
		n, _ := strconv.ParseInt(a.V, 8, 64)
		return res.suss(&Int{n}, c)
	case "bin":
		n, _ := strconv.ParseInt(a.V, 2, 64)
        if len(a.V) <= 8 {
            return res.suss(&Byte{byte(n)}, c)
        }
		return res.suss(&Int{n}, c)
	case "int":
		n, _ := strconv.ParseInt(a.V, 10, 64)
		return res.suss(&Int{n}, c)
	case "float":
		n, _ := strconv.ParseFloat(a.V, 10)
		return res.suss(&Flt{n}, c)
	case "string":
		return res.suss(&Str{a.V}, c)
	case "bool":
		return res.suss(&Bool{a.V == "true"}, c)
	case "array":
		fallthrough
	case "list":
		types := &[]*Type{}
		for _, ast := range a.C {
			t := res.reg(ExecAst(ast, c))
			if res.has() {
				return res
			}
			*types = append(*types, t)
		}
		return res.suss(&List{types}, c)
    case "range":
        s := res.reg(ExecAst(a.C[0], c))
        if res.has() {
            return res
        }
        if (*s).GetT() != "int" {
            return res.fail(MakeError("RangeError", "Range start must be int", a.C[0].S, a.C[0].E, c))
        }
        start := (*s).(*Int).val
        e := res.reg(ExecAst(a.C[1], c))
        if res.has() {
            return res
        }
        if (*e).GetT() != "int" {
            return res.fail(MakeError("RangeError", "Range end must be int", a.C[0].S, a.C[0].E, c))
        }
        end := (*e).(*Int).val
        types := &[]*Type{}
        if end > start {
            for i := start; i < end; i++ {
                t := Type(&Int{i})
                *types = append(*types, &t)
            }
        } else {
            for i := start; i > end; i-- {
                t := Type(&Int{i})
                *types = append(*types, &t)
            }
        }
		return res.suss(&List{types}, c)
	case "obj":
		ctx := MakeContext(c, a.S)
		types := []string{}
		if a.V != "" {
			types = strings.Split(a.V, " ")
		}
		o := &Obj{ctx, types}
		to := Type(o)
		ctx.set("this", &to)
		for _, ast := range a.C {
			res.reg(ExecAst(ast, ctx))
			if res.has() {
				return res
			}
		}
		return res.suss(o, c)
	case "fun":
		return res.suss(&Fun{(a.C[0]), (a.C[1]), false}, c)
	case "afun":
		return res.suss(&Fun{(a.C[0]), (a.C[1]), true}, c)
	case "access":
		t := c.get(a.V)
		if t == nil {
			return res.fail(MakeRTError("Access to undefined: "+a.V, a.S, a.E, c))
		}
		res.n = a.V
		/*if (*t).GetT() == "obj" {
			o := (*t).(*Obj)
			o.vals.parent = c
		}*/
		return res.sussP(t, c)
	case "member":
		t := res.reg(ExecAst(a.C[0], c))
		if res.has() {
			return res
		}
        return (*t).GetMember(a.C[1], c)
	case "assign":
		res.reg(execAssign(a, c, c))
		if res.has() {
			return res
		}
		res.c = c
		return res
	case "_index":
		res.reg(execOperator("_index_get", a.C[0], a.C[1], c))
		if res.has() {
			return res
		}
		res.c = c
		return res
	case "_call":
		res.reg(execFun(a.C[0], a.C[1], c, true))
		if res.has() {
			return res
		}
		res.c = c
		return res
	case "bin_operator":
		args := []*AST{a.C[1]}
		res.reg(execOperator(a.V, a.C[0], &AST{T: "args", C: args, S: a.C[1].S, E: a.C[1].E}, c))
		if res.has() {
			return res
		}
		res.c = c
		return res
	case "unary_operator":
		res.reg(execUnaryOperator(a.V, a.C[0], c))
		if res.has() {
			return res
		}
		res.c = c
		return res
	case "try":
		ctx := MakeContext(c, a.S)
		res.reg(ExecAst(a.C[0], ctx))
		if !res.has() {
			res.c = c
			return res
		}
		etype, edes := Type(&Str{res.e.t}), Type(&Str{res.e.v})
		res.e = nil
		vals := []*Type{&etype, &edes}
		ctx = MakeContext(c, a.S)
		for i, na := range a.C[1].C {
			var val *Type
			if i < len(vals) {
				val = (vals)[i]
			} else {
				valt := Type(&Null{})
				val = &valt
			}
			execAssignVal(na, val, ctx)
		}
		res.reg(ExecAst(a.C[2], ctx))
		if res.has() {
			return res
		}
		res.c = c
		return res
	case "if":
		t := res.reg(ExecAst(a.C[0], c))
		ctx := MakeContext(c, a.S)
		if !res.has() && (*t).getBool() {
			res.reg(ExecAst(a.C[1], ctx))
			if res.has() {
				return res
			}
			res.c = c
			return res
		}
		res.e = nil
		if len(a.C) > 2 {
			res.reg(ExecAst(a.C[2], c))
			if res.has() {
				return res
			}
			res.c = c
			return res
		}
		return res.suss(&Null{}, c)
	case "switch":
		t := res.reg(ExecAst(a.C[0], c))
		if res.has() {
			return res
		}
		ctx := MakeContext(c, a.S)
		exec := false
		rets := []*Type{}
		for _, cas := range a.C[1].C {
			if !exec && cas.T == "case" {
				ccond := res.reg(ExecAst(cas.C[0], ctx))
				if res.has() {
					return res
				}
				if (*t).equals(ccond) {
					exec = true
				}
			}
			if exec || cas.T == "default" {
				if cas.T == "default" {
					ret := res.reg(ExecAst(cas.C[0], ctx))
					if res.has() {
						return res
					}
					rets = append(rets, ret)
					if res.b || res.r {
						res.b = false
						break
					}
				} else if len(cas.C) > 1 {
					ret := res.reg(ExecAst(cas.C[1], ctx))
					if res.has() {
						return res
					}
					rets = append(rets, ret)
					if res.b || res.r {
						res.b = false
						break
					}
				}
			}
		}
		return res.suss(&List{&rets}, c)
	case "for":
		ctx := MakeContext(c, a.S)
		cond := a.C[0].C
		body := a.C[1]
		rets := &[]*Type{}
		switch len(cond) {
		case 0:
			for {
				res.reg(execForBody(rets, body, ctx))
				if res.has() {
					return res
				}
				if res.b || res.r {
					res.b = false
					break
				}
			}
		case 1:
			t := res.reg(ExecAst(cond[0], ctx))
			for !res.has() && (*t).getBool() {
				res.reg(execForBody(rets, body, ctx))
				if res.has() {
					return res
				}
				if res.b || res.r {
					res.b = false
					break
				}
				t = res.reg(ExecAst(cond[0], ctx))
				if res.has() {
					return res
				}
			}
		case 2:
			t := res.reg(ExecAst(cond[1], ctx))
			switch (*t).GetT() {
			case "string":
				l := (*t).(*Str)
				for key, val := range l.val {
					t := Type(&Int{int64(key)})
					v := Type(&Str{string(val)})
					tl := Type(&List{&[]*Type{&t, &v}})
					execAssignVal(cond[0], &tl, ctx)
					res.reg(execForBody(rets, body, ctx))
					if res.has() {
						return res
					}
					if res.b || res.r {
						res.b = false
						break
					}
				}
				return res.suss(&List{rets}, c)
			case "list":
				l := (*t).(*List)
				for key, val := range *l.vals {
					t := Type(&Int{int64(key)})
					tl := Type(&List{&[]*Type{&t, val}})
					execAssignVal(cond[0], &tl, ctx)
					res.reg(execForBody(rets, body, ctx))
					if res.has() {
						return res
					}
					if res.b || res.r {
						res.b = false
						break
					}
				}
				return res.suss(&List{rets}, c)
			case "obj":
				o := (*t).(*Obj)
				for key, val := range *o.vals.st {
					if key != "this" {
						t := Type(&Str{key})
						tl := Type(&List{&[]*Type{&t, val}})
						execAssignVal(cond[0], &tl, ctx)
						res.reg(execForBody(rets, body, ctx))
						if res.has() {
							return res
						}
						if res.b || res.r {
							res.b = false
							break
						}
					}
				}
				return res.suss(&List{rets}, c)
			}
			res.fail(MakeSError("For in can just be apply to types list and obj.", cond[1].S, cond[1].E, ctx))
		case 3:
			res.reg(ExecAst(cond[0], ctx))
			if res.has() {
				return res
			}
			t := res.reg(ExecAst(cond[1], ctx))
			for !res.has() && (*t).getBool() {
				res.reg(execForBody(rets, body, ctx))
				if res.has() {
					return res
				}
				if res.b || res.r {
					res.b = false
					break
				}
				res.reg(ExecAst(cond[2], ctx))
				if res.has() {
					return res
				}
				t = res.reg(ExecAst(cond[1], ctx))
			}
		}
		return res.suss(&List{rets}, c)
	}
	return res.fail(MakeSError("Invalid syntax", a.S, a.E, c))
}

func execForBody(rets *[]*Type, body *AST, ctx *Context) *Result {
	res := MkRes()
	t := res.reg(ExecAst(body, ctx))
	if res.has() {
		return res
	}
	*rets = append(*rets, funCopyTypes(t, ctx))
	return res.sussP(t, ctx)
}

func execAssign(a *AST, c, ac *Context) *Result {
	res := MkRes()
	l := a.C[0]
	if l.T == "access" {
		t := res.reg(ExecAst(a.C[1], ac))
		if res.has() {
			return res
		}
		c.set(l.V, t)
		return res.sussP(t, c)
	}
	if l.T == "assign" {
		res.reg(execAssign(&AST{C: []*AST{l.C[1], a.C[0]}, S: l.C[1].S, E: a.E}, c, ac))
		if res.has() {
			return res
		}
		res.reg(execAssign(l, c, ac))
		if res.has() {
			return res
		}
		res.c = c
		return res
	}
	if l.T == "member" {
		t := res.reg(ExecAst(l.C[0], c))
		if (*t).GetT() != "obj" {
			return res.fail(MakeRTError("Type "+(*t).GetT()+" has not members.", l.S, l.E, c))
		}
		o, _ := (*t).(*Obj)
		return o.SetMemb(l.C[1], a.C[1], c)
	}
	if l.T == "_index" {
		args := &AST{T: "args", C: []*AST{l.C[1], a.C[1]}, S: l.C[1].S, E: a.C[1].E}
		res.reg(execOperator("_index_set", l.C[0], args, c))
		if res.has() {
			return res
		}
		res.c = c
		return res
	}
	if l.T == "list" || l.T == "array" {
		t := res.reg(ExecAst(a.C[1], c))
		if res.has() {
			return res
		}
		vals := &[]*Type{}
		switch (*t).GetT() {
		case "list":
			vals = (*t).(*List).vals
		case "string":
			for _, v := range (*t).(*Str).val {
				val := Type(&Str{string(v)})
				*vals = append(*vals, &val)
			}
		default:
			*vals = append(*vals, t)
		}
		if len(*vals) == 1 {
			val := (*vals)[0]
			for _, na := range l.C {
				execAssignVal(na, val, c)
			}
		} else {
			for i, na := range l.C {
				var val *Type
				if i < len(*vals) {
					val = (*vals)[i]
				} else {
					valt := Type(&Null{})
					val = &valt
				}
				execAssignVal(na, val, c)
			}
		}
		return res.sussP(t, c)
	}
	return res.fail(MakeSError("A "+l.T+" can not be reassign.", a.S, a.E, c))
}

func execAssignVal(a *AST, v *Type, c *Context) *Result {
	res := MkRes()
	switch a.T {
	case "access":
		c.set(a.V, v)
		return res.sussP(v, c)
	case "member":
		t := res.reg(ExecAst(a.C[0], c))
		if (*t).GetT() != "obj" {
			return res.fail(MakeRTError("Type "+(*t).GetT()+" has not members.", a.S, a.E, c))
		}
		o, _ := (*t).(*Obj)
		res.reg(execAssignVal(a.C[1], v, o.vals))
		if res.has() {
			return res
		}
		return res.sussP(v, c)
	case "_index":
		ctx := MakeContext(c, a.S)
		ctx.set("JGSTempValue", v)
		val := &AST{T: "access", V: "JGSTempValue", S: a.S, E: a.E}
		args := &AST{T: "args", C: []*AST{a.C[1], val}, S: a.C[1].S, E: val.E}
		res.reg(execOperator("_index_set", a.C[0], args, ctx))
		if res.has() {
			return res
		}
		return res.sussP(v, c)
	case "list":
		fallthrough
	case "array":
		vals := &[]*Type{}
		switch (*v).GetT() {
		case "list":
			vals = (*v).(*List).vals
		case "string":
			for _, e := range (*v).(*Str).val {
				val := Type(&Str{string(e)})
				*vals = append(*vals, &val)
			}
		default:
			*vals = append(*vals, v)
		}
		if len(*vals) == 1 {
			val := (*vals)[0]
			for _, na := range a.C {
				execAssignVal(na, val, c)
			}
		} else {
			for i, na := range a.C {
				var val *Type
				if i < len(*vals) {
					val = (*vals)[i]
				} else {
					valt := Type(&Null{})
					val = &valt
				}
				execAssignVal(na, val, c)
			}
		}
		return res.sussP(v, c)
	}
	return res.fail(MakeSError("Invalid operation.", a.S, a.E, c))
}

func execTypeCall(a, call *AST, c *Context) *Result {
	res := MkRes()
	res.reg(ExecAst(a, c))
	if res.has() {
		return res
	}
	res.c = c
	return res
}

func execCall(a, args *AST, c *Context) *Result {
	res := MkRes()
	t := res.reg(ExecAst(a, c))
	if res.has() || t == nil {
		return nil
	}
	if (*t).GetT() == "fun" {
	    f, _ := (*t).(*Fun)
	    ret := res.reg(f.ExecFun(args, res.n, res.c, c))
	    if res.has() {
		    /*if res.se {
			    res.se = false
			    res.e.s = a.S
			    res.e.e = a.E
		    }*/
		    return res
	    }
	    res.b = false
	    res.r = false
	    return res.sussP(ret, c)
    }
    if (*t).GetT() == "nfun" {
        f,_ := (*t).(*NFun)
        if a.T == "member" {
	        t = res.reg(ExecAst(a.C[0], c))
        }
	    ret := res.reg(f.ExecFun(t, args, c))
	    if res.has() {
		    return res
	    }
	    res.b = false
	    res.r = false
	    return res.sussP(ret, c)
    }
	return res.fail(MakeSError("Can not call type: " + (*t).GetT(), a.S, a.E, c))
}

func execFun(a, args *AST, c *Context, sys bool) *Result {
	var res *Result = execCall(a, args, c)
	if res != nil {
		if res.has() {
			return res
		}
		res.c = c
		return res
	}
	if sys {
		res := ExecSysFun(a, args, c)
		if res != nil {
			if res.has() {
				return res
			}
            res.b = false
            res.r = false
			res.c = c
			return res
		}
	}
	return MkRes().fail(MakeRTError("Access to undefined", a.S, args.E, c))
}

func execFunArgs(def, arg *AST, ctx, ac *Context) *Result {
	res := MkRes()
	d_args := def.C
	n_args := arg.C
	argOrder := true
	list := false
	name := ""
	var last *Type = nil
	listT := []*Type{}
	for i, a := range n_args {
		if a.T == "DDD" {
			if last == nil || (*last).GetT() != "list" {
				return res.fail(MakeSError("Invalid function arguments.", a.S, a.E, ctx))
			}
			for i, e := range listT {
				if e == last {
					listT = append(listT[:i], listT[i+1:]...)
				}
			}
			l := (*last).(*List).vals
			for pos, ele := range *l {
				if !argOrder {
					return res.fail(MakeSError("Invalid function arguments", a.S, a.E, ctx))
				}
				if list {
					listT = append(listT, ele)
					last = ele
				} else if d_args[i+pos].T == "DDD" {
					list = true
					if last != nil {
						listT = append(listT, last)
					}
					listT = append(listT, ele)
					last = ele
				} else {
					if i+pos >= len(d_args) {
						return res.fail(MakeSError("Too many arguments", a.S, a.E, ctx))
					}
					var n string
					if d_args[i+pos].T == "assign" {
						n = d_args[i+pos].C[0].V
					} else if d_args[i].T == "access" {
						n = d_args[i+pos].V
					}
					name = n
					last = ele
					ctx.set(n, ele)
				}
			}
		} else if a.T == "assign" {
			argOrder = false
			res.reg(execAssign(a, ctx, ac))
			if res.has() {
				return res
			}
		} else {
			if !argOrder {
				return res.fail(MakeSError("Invalid function arguments", a.S, a.E, ctx))
			}
			t := res.reg(ExecAst(a, ac))
			if res.has() {
				return res
			}
			if list {
				listT = append(listT, t)
				last = t
            } else if i >= len(d_args) {
				return res.fail(MakeSError("Too many arguments", a.S, a.E, ctx))
			} else if d_args[i].T == "DDD" {
				list = true
				if last != nil {
					listT = append(listT, last)
				}
				listT = append(listT, t)
				last = t
			} else {
				var n string = "a"
				if d_args[i].T == "assign" {
					n = d_args[i].C[0].V
				} else if d_args[i].T == "access" {
					n = d_args[i].V
				}
				name = n
				last = t
				ctx.set(n, t)
			}
		}
	}
	if list {
		if name == "" {
			name = "args"
		}
		t := Type(&List{&listT})
		ctx.set(name, &t)
	}
	for _, a := range d_args {
        if a.T != "DDD" {
            n := ""
		    if a.T == "assign" {
		    	n = a.C[0].V
		    } else if a.T == "access" {
			    n = a.V
		    }
            if (*ctx.st)[n] == nil {
		        if a.T == "assign" {
			        res.reg(execAssign(a, ctx, ac))
			        if res.has() {
				        return res
			        }
		        } else if a.T == "access" {
			        n := Type(&Null{})
			        ctx.set(a.V, &n)
		        }
            }
        }
	}
	return res.suss(&Null{}, ctx)
}

func execOperator(op string, to, a *AST, c *Context) *Result {
	res := MkRes()
	t := res.reg(ExecAst(to, c))
	if res.has() {
		return res
	}
	if (*t).GetT() == "obj" && op != "_right_arrow" {
		o, _ := (*t).(*Obj)
		r := o.objOperator(op, a, c)
		if r != nil {
			if res.has() {
				return res
			}
			r.c = c
			return r
		}
	} else if (*t).GetT() == "list" {
		l, _ := (*t).(*List)
		if op == "_index_get" {
			res.reg(l.GetIndex(a, c))
			if res.has() {
				return res
			}
			res.c = c
			return res
		}
		if op == "_index_set" {
			res.reg(l.SetIndex(a.C[0], a.C[1], c))
			if res.has() {
				return res
			}
			res.c = c
			return res
		}
	} else if (*t).GetT() == "string" {
		l, _ := (*t).(*Str)
		if op == "_index_get" {
			res.reg(l.GetIndex(a, c))
			if res.has() {
				return res
			}
			res.c = c
			return res
		}
		if op == "_index_set" {
			res.reg(l.SetIndex(a.C[0], a.C[1], c))
			if res.has() {
				return res
			}
			res.c = c
			return res
		}
	}
	o := res.reg(ExecAst(a.C[0], c))
	if res.has() {
		return res
	}
	res.reg((*t).operator(op, o, c, to.S, a.E))
	if res.has() {
		return res
	}
	res.c = c
	return res
}

func execUnaryOperator(op string, to *AST, c *Context) *Result {
	res := MkRes()
	t := res.reg(ExecAst(to, c))
	if res.has() {
		return res
	}
	a := Type(&Null{})
	res.reg((*t).operator(op, &a, c, to.S, to.E))
	if res.has() {
		return res
	}
	res.c = c
	return res
}
