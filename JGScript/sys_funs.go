package JGScript

import (
	"bufio"
	"fmt"
	"io/ioutil"
    "math"
	"math/rand"
	"os"
	"path"
	"plugin"
	"strings"
	"time"
    "encoding/base64"
)

func ExecSysFun(f, arg *AST, c *Context) *Result {
	if f.T == "access" {
        res := MathFuns(f, arg, c)
        if res != nil {
            return res
        }
		res = MkRes()
		ctx := MakeContext(c, f.S)
		ctx.n = f.V
		alen := len(arg.C)
		switch f.V {
		// NATIVE GO
		case "native_go":
			if alen < 3 {
				return res.fail(MakeRTError("Arguments required: 3", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], c))
			if res.has() {
				return res
			}
			if (*t).GetT() != "string" {
				return res.fail(&Error{"NativeCallError", "First argument must be string", arg.C[0].S, arg.C[0].E, ctx})
			}
			fn := (*t).(*Str).val
			t = res.reg(ExecAst(arg.C[1], c))
			if res.has() {
				return res
			}
			if (*t).GetT() != "string" {
				return res.fail(&Error{"NativeCallError", "Second argument must be string", arg.C[1].S, arg.C[1].E, ctx})
			}
			fun := (*t).(*Str).val
			t = res.reg(ExecAst(arg.C[2], c))
			if res.has() {
				return res
			}
			if (*t).GetT() != "obj" {
				return res.fail(&Error{"NativeCallError", "Third argument must be object", arg.C[2].S, arg.C[2].E, ctx})
			}
			nargs := (*t).(*Obj).vals
			dir, _ := path.Split(*c.s.Fn)
			plug, err := plugin.Open(path.Join(dir, fn))
			if err != nil {
				plug, err = plugin.Open(fn)
				if err != nil {
					return res.fail(&Error{"NativeCallError", "Load plugin fail: " + fn, f.S, arg.E, ctx})
				}
			}
			context, err := plug.Lookup("CTX")
			if err != nil {
				return res.fail(&Error{"NativeCallError", "Error loading CTX", arg.S, arg.E, ctx})
			}
            *context.(**Context) = nargs
			result, err := plug.Lookup("RES")
			if err != nil {
				return res.fail(&Error{"NativeCallError", "Error loading RES", arg.S, arg.E, ctx})
			}
            *result.(**Result) = res
			f, err := plug.Lookup(fun)
			if err != nil {
				return res.fail(&Error{"NativeCallError", "Error loading function: " + fun, arg.S, arg.E, ctx})
			}
			f.(func())()
			return res
			// CTRL
		case "panic":
			if alen < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], c))
			if res.has() {
				return res
			}
			d := res.reg(ExecAst(arg.C[1], c))
			if res.has() {
				return res
			}
			res.se = true
			return res.fail(&Error{(*t).str(), (*d).str(), arg.S, arg.E, ctx})
			// IO
		case "print":
			return execPrint(arg, ctx, "\n")
		case "input":
			execPrint(arg, ctx, "")
			reader := bufio.NewReader(os.Stdin)
			in, _ := reader.ReadString('\n')
			input := strings.TrimRight(in, "\r\n")
			return res.suss(&Str{input}, ctx)
        case "eval":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1 string", arg.S, arg.E, ctx))
			}
			pa := arg.C[0]
			pt := res.reg(ExecAst(pa, ctx))
			if res.has() {
				return res
			}
			if (*pt).GetT() != "string" {
				return res.fail(MakeRTError("The argument must be string.", f.S, arg.E, ctx))
			}
			p := (*pt).(*Str).val
            ret := res.reg(EvalCode("RunTimeCode", "EvalRunTimeCode", p, ctx))
            if res.has() {
                return res
            }
            res.r = false
            res.b = false
            return res.sussP(ret, c)
			// Primitives
        case "NaN":
            return res.suss(&Flt{math.NaN()}, ctx)
        case "Inf":
            return res.suss(&Flt{math.Inf(1)}, ctx)
        case "nInf":
            return res.suss(&Flt{math.Inf(-1)}, ctx)
		case "copy":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], c))
			if res.has() {
				return res
			}
			return funCopy(t, c)
        case "byte":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			return res.suss(&Byte{(*t).getByte()}, ctx)
        case "bytes":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			la := arg.C[0]
			l := res.reg(ExecAst(la, ctx))
			if res.has() {
				return res
			}
			if (*l).GetT() == "string" {
				lis := &((*l).(*Str).val)
                a := []*Type{}
                for _, c := range *lis {
                    t := Type(&Byte{byte(c)})
                    a = append(a, &t)
                }
				return res.suss(&List{&a}, ctx)
			}
            if (*l).GetT() == "list" {
				lis := ((*l).(*List).vals)
                a := []*Type{}
                for _, c := range *lis {
                    t := Type(&Byte{byte((*c).getByte())})
                    a = append(a, &t)
                }
				return res.suss(&List{&a}, ctx)
            }
			return res.fail(MakeRTError("Invalid function arguments.", f.S, arg.E, ctx))
		case "int":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			return res.suss(&Int{(*t).getInt()}, ctx)
		case "float":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			return res.suss(&Flt{(*t).getFloat()}, ctx)
		case "string":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			return res.suss(&Str{(*t).str()}, ctx)
		case "string_from_bytes":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
            if (*t).GetT() != "list" {
				return res.fail(MakeRTError("Argument must be el list of bytes", arg.S, arg.E, ctx))
            }
            list := (*t).(*List).vals
            s := ""
            for _, b := range *list {
                if (*b).GetT() != "byte" {
				    return res.fail(MakeRTError("Argument must be el list of bytes", arg.S, arg.E, ctx))
                }
                s += string((*b).(*Byte).val)
            }
			return res.suss(&Str{s}, ctx)
		case "base64_encode":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
            if (*t).GetT() != "list" {
				return res.fail(MakeRTError("Argument must be el list of bytes", arg.S, arg.E, ctx))
            }
            list := (*t).(*List).vals
            arr := []byte {}
            for _, b := range *list {
                if (*b).GetT() != "byte" {
				    return res.fail(MakeRTError("Argument must be el list of bytes", arg.S, arg.E, ctx))
                }
                arr = append(arr, (*b).(*Byte).val)
            }
			return res.suss(&Str{base64.StdEncoding.EncodeToString(arr)}, ctx)
        case "base64_decode":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			la := arg.C[0]
			l := res.reg(ExecAst(la, ctx))
			if res.has() {
				return res
			}
			if (*l).GetT() == "string" {
				lis, err := base64.StdEncoding.DecodeString((*l).(*Str).val)
                if err != nil {
			        return res.fail(MakeError("Base64Error", fmt.Sprint(err), f.S, arg.E, ctx))
                }
                a := []*Type{}
                for _, c := range lis {
                    t := Type(&Byte{c})
                    a = append(a, &t)
                }
				return res.suss(&List{&a}, ctx)
			}
			return res.fail(MakeRTError("Invalid function arguments.", f.S, arg.E, ctx))
		case "char":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			return res.suss(&Str{string(byte((*t).getInt()))}, ctx)
        case "to_hex_str":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			return res.suss(&Str{IntToHexStr((*t).getInt())}, ctx)
        case "to_oct_str":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			return res.suss(&Str{IntToOctStr((*t).getInt())}, ctx)
        case "to_bin_str":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			return res.suss(&Str{IntToBinStr((*t).getInt())}, ctx)
		case "bool":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			return res.suss(&Bool{(*t).getBool()}, ctx)
			// LIST
		case "range":
			i := float64(0)
			s := float64(1)
			e := float64(0)
			switch alen {
			case 1:
				t := res.reg(ExecAst(arg.C[0], ctx))
				if res.has() {
					return res
				}
				if (*t).GetT() == "int" {
				    e = float64((*t).(*Int).val)
				} else if (*t).GetT() == "float" {
				    e = (*t).(*Flt).val
                } else {
				    return res.fail(MakeRTError("Arguments type must be int or float.", arg.C[0].S, arg.C[0].E, ctx))
                }
			case 2:
				t := res.reg(ExecAst(arg.C[0], ctx))
				if res.has() {
					return res
				}
				if (*t).GetT() == "int" {
				    i = float64((*t).(*Int).val)
				} else if (*t).GetT() == "float" {
				    i = (*t).(*Flt).val
                } else {
				    return res.fail(MakeRTError("Arguments type must be int or float.", arg.C[0].S, arg.C[0].E, ctx))
                }
				t = res.reg(ExecAst(arg.C[1], ctx))
				if res.has() {
					return res
				}
				if (*t).GetT() == "int" {
				    e = float64((*t).(*Int).val)
				} else if (*t).GetT() == "float" {
				    e = (*t).(*Flt).val
                } else {
				    return res.fail(MakeRTError("Arguments type must be int or float.", arg.C[1].S, arg.C[1].E, ctx))
                }
			case 3:
				t := res.reg(ExecAst(arg.C[0], ctx))
				if res.has() {
					return res
				}
				if (*t).GetT() == "int" {
				    i = float64((*t).(*Int).val)
				} else if (*t).GetT() == "float" {
				    i = (*t).(*Flt).val
                } else {
				    return res.fail(MakeRTError("Arguments type must be int or float.", arg.C[0].S, arg.C[0].E, ctx))
                }
				t = res.reg(ExecAst(arg.C[1], ctx))
				if res.has() {
					return res
				}
				if (*t).GetT() == "int" {
				    e = float64((*t).(*Int).val)
				} else if (*t).GetT() == "float" {
				    e = (*t).(*Flt).val
                } else {
				    return res.fail(MakeRTError("Arguments type must be int or float.", arg.C[1].S, arg.C[1].E, ctx))
                }
				t = res.reg(ExecAst(arg.C[2], ctx))
				if res.has() {
					return res
				}
				if (*t).GetT() == "int" {
				    s = float64((*t).(*Int).val)
				} else if (*t).GetT() == "float" {
				    s = (*t).(*Flt).val
                } else {
			        return res.fail(MakeRTError("Arguments type must be int or float.", arg.C[2].S, arg.C[2].E, ctx))
                }
			default:
				res.fail(MakeRTError("This function require 1, 2 or 3 arguments.", f.S, arg.E, ctx))
			}
			array := []*Type{}
			if e < i {
				for v := i; v > e; v -= s {
					it := Type(&Flt{v})
					array = append(array, &it)
				}
			} else {
				for v := i; v < e; v += s {
					it := Type(&Flt{v})
					array = append(array, &it)
				}
			}
			return res.suss(&List{&array}, ctx)
		case "zip":
			if alen < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, ctx))
			}
			list := []*[]*Type{}
			size := -1
			for _, na := range arg.C {
				t := res.reg(ExecAst(na, ctx))
				if res.has() {
					return res
				}
				if (*t).GetT() != "list" {
					res.fail(MakeRTError("Arguments type must be list.", f.S, arg.E, ctx))
				}
				l := (*t).(*List).vals
				if size == -1 {
					size = len(*l)
				} else if size != len(*l) {
					res.fail(MakeRTError("List lengths must be the same.", f.S, arg.E, ctx))
				}
				list = append(list, l)
			}
			lis := &[]*Type{}
			for i := 0; i < size; i++ {
				tl := &[]*Type{}
				for _, e := range list {
					*tl = append(*tl, (*e)[i])
				}
				t := Type(&List{tl})
				*lis = append(*lis, &t)
			}
			return res.suss(&List{lis}, ctx)
		case "delete":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			la := arg.C[0]
			l := res.reg(ExecAst(la, ctx))
			if res.has() {
				return res
			}
			if (*l).GetT() == "string" {
				name := ((*l).(*Str).val)
				ret := c.remove(name)
				if ret == nil {
					return res.suss(&Null{}, ctx)
				}
				return res.sussP(ret, ctx)
			}
			return res.fail(MakeRTError("Invalid function arguments.", f.S, arg.E, ctx))
		case "exist":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			la := arg.C[0]
			l := res.reg(ExecAst(la, ctx))
			if res.has() {
				return res
			}
			if (*l).GetT() == "string" {
				name := ((*l).(*Str).val)
				return res.suss(&Bool{c.exist(name)}, c)
			}
			return res.fail(MakeRTError("Invalid function arguments.", f.S, arg.E, ctx))
			// GetType
		case "is_type":
			if alen < 2 {
				return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, ctx))
			}
			oa := arg.C[0]
			ta := arg.C[1]
			o := res.reg(ExecAst(oa, ctx))
			if res.has() {
				return res
			}
			if (*o).GetT() == "obj" {
				lis := &((*o).(*Obj).types)
				t := res.reg(ExecAst(ta, ctx))
				if res.has() {
					return res
				}
				if (*t).GetT() != "string" {
					return res.fail(MakeRTError("Second argument must be string.", f.S, arg.E, ctx))
				}
				ts := (*t).(*Str).val
				if ts == "obj" {
					return res.suss(&Bool{true}, ctx)
				}
				for _, val := range *lis {
					if ts == val {
						return res.suss(&Bool{true}, ctx)
					}
				}
			}
			return res.suss(&Bool{false}, ctx)
		case "type":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			oa := arg.C[0]
			o := res.reg(ExecAst(oa, ctx))
			if res.has() {
				return res
			}
			if (*o).GetT() == "obj" {
				lis := ((*o).(*Obj).types)
				t := "obj"
				if len(lis) > 0 {
					t = (lis)[0]
				}
				return res.suss(&Str{t}, ctx)
			}
			return res.suss(&Str{(*o).GetT()}, ctx)
			// Others
		case "sleep":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			t := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
            time_sec := float64(1)
            if (*t).GetT() == "int" {
                time_sec = float64((*t).(*Int).val)
            } else if (*t).GetT() == "float" {
                time_sec = (*t).(*Flt).val
            } else {
				return res.fail(MakeRTError("Argument must be the number of seconds", arg.S, arg.E, ctx))
            }
            time.Sleep(time.Duration(time_sec * float64(time.Second)))
			return res.suss(&Null{}, ctx)
        case "remove_file":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			pa := arg.C[0]
			pt := res.reg(ExecAst(pa, ctx))
			if res.has() {
				return res
			}
			if (*pt).GetT() != "string" {
				return res.fail(MakeRTError("The argument must be string.", f.S, arg.E, ctx))
			}
			p := (*pt).(*Str).val
			err := os.Remove(p)
			if err != nil {
                return res.fail(MakeRTError("Fail to remove file: " + p, f.S, arg.E, ctx))
			}
            return res.suss(&Null{}, ctx)
        case "read_file":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			pa := arg.C[0]
			pt := res.reg(ExecAst(pa, ctx))
			if res.has() {
				return res
			}
			if (*pt).GetT() != "string" {
				return res.fail(MakeRTError("The argument must be string.", f.S, arg.E, ctx))
			}
			p := (*pt).(*Str).val
			code, err := ioutil.ReadFile(p)
			if err != nil {
                return res.fail(MakeRTError("Fail to read file: " + p, f.S, arg.E, ctx))
			}
            list := []*Type{}
            for _,b := range code {
                t := Type(&Byte{b})
                list = append(list, &t)
            }
            return res.suss(&List{&list}, ctx)
        case "write_file":
			if alen < 3 {
				return res.fail(MakeRTError("Arguments required: 3", arg.S, arg.E, ctx))
			}
			pa := arg.C[0]
			pt := res.reg(ExecAst(pa, ctx))
			if res.has() {
				return res
			}
			if (*pt).GetT() != "string" {
				return res.fail(MakeRTError("The first argument must be string.", f.S, arg.E, ctx))
			}
			p := (*pt).(*Str).val
			pa = arg.C[1]
			pt = res.reg(ExecAst(pa, ctx))
			if res.has() {
				return res
			}
			if (*pt).GetT() != "list" {
				return res.fail(MakeRTError("The second argument must be a list of bytes.", f.S, arg.E, ctx))
			}
			c := (*pt).(*List).vals
			pa = arg.C[2]
			pt = res.reg(ExecAst(pa, ctx))
			if res.has() {
				return res
			}
			if (*pt).GetT() != "int" {
				return res.fail(MakeRTError("The third argument must be int.", f.S, arg.E, ctx))
			}
			pu := (*pt).(*Int).val
            cb := []byte{}
            for _, b := range *c {
                if (*b).GetT() != "byte" {
				    return res.fail(MakeRTError("The second argument must be a list of bytes.", f.S, arg.E, ctx))
                }
                cb = append(cb, (*b).(*Byte).val)
            }
			err := ioutil.WriteFile(p, cb, os.FileMode(pu))
			if err != nil {
                return res.fail(MakeRTError("Fail to write file: " + p, f.S, arg.E, ctx))
			}
            return res.suss(&Null{}, ctx)
		case "get_time":
			t := time.Now()
			nano := Type(&Int{int64(t.Nanosecond())})
			sec := Type(&Int{int64(t.Second())})
			min := Type(&Int{int64(t.Minute())})
			hou := Type(&Int{int64(t.Hour())})
			wday := Type(&Int{int64(t.Weekday())})
			day := Type(&Int{int64(t.Day())})
			mon := Type(&Int{int64(t.Month())})
			yea := Type(&Int{int64(t.Year())})
			ctx := MakeContext(c, arg.S)
			ctx.set("nsec", &nano)
			ctx.set("second", &sec)
			ctx.set("minute", &min)
			ctx.set("hour", &hou)
			ctx.set("weekday", &wday)
			ctx.set("day", &day)
			ctx.set("month", &mon)
			ctx.set("year", &yea)
			return res.suss(&Obj{ctx, []string{"Time"}}, c)
            // MathFuns
		case "random":
			return res.suss(&Flt{rand.Float64()}, ctx)
		case "random_seed":
			if alen < 1 {
				return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
			}
			i := res.reg(ExecAst(arg.C[0], ctx))
			if res.has() {
				return res
			}
			if (*i).GetT() != "int" {
				return res.fail(MakeRTError("Index must be int.", arg.C[0].S, arg.C[0].E, ctx))
			}
			it, _ := (*i).(*Int)
			ind := it.val
			rand.Seed(ind)
			return res.suss(&Null{}, ctx)
        case "max":
			if alen < 1 { return res.fail(MakeRTError("Arguments required: 2 or more numbers, or a list of numbers", arg.S, arg.E, ctx)) }
            init := false
            var max float64
            var num float64
            pos := -1
            if alen == 1 {
                listt := res.reg(ExecAst(arg.C[0], ctx))
                if (*listt).GetT() != "list" {
				    return res.fail(MakeRTError("Arguments required: 2 or more numbers, or a list of numbers", arg.S, arg.E, ctx))
                }
                list := (*listt).(*List)
                for i, tv := range *list.vals {
			        switch (*tv).GetT() {
                    case "int":
                        num = float64((*tv).(*Int).val)
                    case "float":
                        num = (*tv).(*Flt).val
                    default:
				        return res.fail(MakeRTError("Arguments required: 2 or more numbers, or a list of numbers", arg.S, arg.E, ctx))
			        }
                    if !init {
                        max = num
                        pos = i
                        init = true
                    } else {
                        max = math.Max(max, num)
                        pos = i
                    }
                }
            } else {
                for i, av := range arg.C {
                    tv := res.reg(ExecAst(av, ctx))
			        if res.has() {
				        return res
			        }
			        switch (*tv).GetT() {
                    case "int":
                        num = float64((*tv).(*Int).val)
                    case "float":
                        num = (*tv).(*Flt).val
                    default:
				        return res.fail(MakeRTError("Arguments required: 2 or more numbers, or a list", av.S, av.E, ctx))
			        }
                    if !init {
                        max = num
                        pos = i
                        init = true
                    } else {
                        max = math.Max(max, num)
                        pos = i
                    }
                }
            }
            i := Type(&Int{int64(pos)})
            m := Type(&Flt{max})
			return res.suss(&List{&[]*Type{&i, &m}}, ctx)
        case "min":
			if alen < 1 { return res.fail(MakeRTError("Arguments required: 2 or more numbers, or a list of numbers", arg.S, arg.E, ctx)) }
            init := false
            var max float64
            var num float64
            pos := -1
            if alen == 1 {
                listt := res.reg(ExecAst(arg.C[0], ctx))
                if (*listt).GetT() != "list" {
				    return res.fail(MakeRTError("Arguments required: 2 or more numbers, or a list of numbers", arg.S, arg.E, ctx))
                }
                list := (*listt).(*List)
                for i, tv := range *list.vals {
			        switch (*tv).GetT() {
                    case "int":
                        num = float64((*tv).(*Int).val)
                    case "float":
                        num = (*tv).(*Flt).val
                    default:
				        return res.fail(MakeRTError("Arguments required: 2 or more numbers, or a list of numbers", arg.S, arg.E, ctx))
			        }
                    if !init {
                        max = num
                        pos = i
                        init = true
                    } else {
                        max = math.Min(max, num)
                        pos = i
                    }
                }
            } else {
                for i, av := range arg.C {
                    tv := res.reg(ExecAst(av, ctx))
			        if res.has() {
				        return res
			        }
			        switch (*tv).GetT() {
                    case "int":
                        num = float64((*tv).(*Int).val)
                    case "float":
                        num = (*tv).(*Flt).val
                    default:
				        return res.fail(MakeRTError("Arguments required: 2 or more numbers, or a list", av.S, av.E, ctx))
			        }
                    if !init {
                        max = num
                        pos = i
                        init = true
                    } else {
                        max = math.Min(max, num)
                        pos = i
                    }
                }
            }
            i := Type(&Int{int64(pos)})
            m := Type(&Flt{max})
			return res.suss(&List{&[]*Type{&i, &m}}, ctx)
		}
	}
	return nil
}

func execPrint(a *AST, c *Context, end string) *Result {
	res := MkRes()
	sep := " "
	max := len(a.C)
	ps := 0
	for _, e := range a.C {
		if e.T == "assign" && e.C[0].T == "access" && e.C[0].V == "sep" {
			max--
			t := res.reg(ExecAst(e.C[1], c))
			if res.has() {
				return res
			}
			sep = (*t).str()
		} else if e.T == "assign" && e.C[0].T == "access" && e.C[0].V == "end" {
			max--
			t := res.reg(ExecAst(e.C[1], c))
			if res.has() {
				return res
			}
			end = (*t).str()
        }
	}
	for _, e := range a.C {
		if !(e.T == "assign" && e.C[0].T == "access" && (e.C[0].V == "sep" || e.C[0].V == "end")) {
			t := res.reg(ExecAst(e, c))
			if res.has() {
				return res
			}
            if (*t).GetT() == "list" {
                printList(t)
            } else {
			    fmt.Print((*t).str())
            }
			ps++
			if ps < max {
				fmt.Print(sep)
			}
		}
	}
	fmt.Print(end)
	return res.suss(&Null{}, c)
}

func printList(t *Type) {
    p := (*t).(*List)
    fmt.Print("[ ")
	max := len(*p.vals) - 1
	for i, e := range *p.vals {
        if (*e).GetT() == "list" {
            printList(e)
        } else {
            fmt.Print((*e).str())
        }
		if i < max {
            fmt.Print(", ")
		}
	}
    fmt.Print(" ]")
}

func funCopy(t *Type, c *Context) *Result {
	res := MkRes()
	switch (*t).GetT() {
    case "byte":
        e := (*t).(*Byte)
		return res.suss(&Byte{(*e).val}, c)
	case "int":
		e := (*t).(*Int)
		return res.suss(&Int{(*e).val}, c)
	case "float":
		e := (*t).(*Flt)
		return res.suss(&Flt{(*e).val}, c)
	case "string":
		e := (*t).(*Str)
		return res.suss(&Str{(*e).val}, c)
	case "list":
		e := (*t).(*List)
		lis := &[]*Type{}
		for _, e := range *(*e).vals {
			*lis = append(*lis, funCopy(e, c).t)
		}
		return res.suss(&List{lis}, c)
	case "obj":
		e := (*t).(*Obj)
		ctx := MakeContext(e.vals.parent, e.vals.s)
		lis := ctx.st
		for k, e := range *(*e).vals.st {
			(*lis)[k] = funCopyTypes(e, c)
		}
		ts := []string{}
		for _, e := range (*e).types {
			ts = append(ts, e)
		}
		return res.suss(&Obj{ctx, ts}, c)
	}
	return res.sussP(t, c)
}

func funCopyTypes(t *Type, c *Context) *Type {
	s := (*t).GetT()
	if s == "byte" || s == "int" || s == "float" || s == "string" {
		return funCopy(t, c).t
	}
	return t
}

func EvalCode(fn, name, code string, ctx *Context) *Result {
	lexer := NewLexer(&fn, &code)
	parser := MakeParser(lexer)
	ast := parser.Parse()
	if ast.T == "body" {
		ast.V = name
	}
	return ExecAst(ast, ctx)
}
