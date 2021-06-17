package JGScript

import (
	"encoding/gob"
	"io/ioutil"
    "path"
    "os"
)

func import_file(arg *AST, ctx *Context) *Result {
    res := MkRes()
	if len(arg.C) < 1 {
		return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, ctx))
	}
	pa := arg.C[0]
	pt := res.reg(ExecAst(pa, ctx))
	if res.has() {
		return res
	}
	if (*pt).GetT() != "string" {
		return res.fail(MakeRTError("The argument must be an string with the path to file.", arg.S, arg.E, ctx))
	}
	p := (*pt).(*Str).val
    var name string
    if len(arg.C) > 1 {
        nt := res.reg(ExecAst(arg.C[1], ctx))
        if res.has() {
            return res
        }
	    if (*nt).GetT() != "string" {
		    return res.fail(MakeRTError("The second argument must be an string.", arg.S, arg.E, ctx))
	    }
	    name = (*nt).(*Str).val
    } else {
	    _, name = path.Split(p)
    }
	dir, _ := path.Split(*ctx.s.Fn)
    codePath := path.Join(dir, p) + ".jgc"
	src, err := os.Open(codePath)
	if err != nil {
        codePath = p + ".jgc"
		src, err = os.Open(codePath)
		if err != nil {
            codePath = path.Join(os.Getenv("HOME"), "JGscript", "libs", p) + ".jgc"
			src, err = os.Open(codePath)
		}
	}
	if err == nil {
		defer src.Close()
	    decoder := gob.NewDecoder(src)
		var ast AST
		decoder.Decode(&ast)
        var context *Context
        if name == "*" {
            context = ctx
        } else {
			context = MakeContext(ctx, ast.S)
			mt := Type(&Obj{context, []string{"module"}})
			ctx.set(name, &mt)
			if ast.T == "body" {
				ast.V = name
			}
        }
		res.reg(ExecAst(&ast, context))
		if res.has() {
			return res
		}
		return res.suss(&Null{}, ctx)
	}
    codePath = path.Join(dir, p) + ".jgs"
	code, err := ioutil.ReadFile(codePath)
	if err != nil {
        codePath = p + ".jgs"
		code, err = ioutil.ReadFile(codePath)
		if err != nil {
            codePath = path.Join(os.Getenv("HOME"), "JGscript", "libs", p) + ".jgs"
			code, err = ioutil.ReadFile(codePath)
			if err != nil {
				return res.fail(MakeRTError("Fail to import "+p, arg.S, arg.E, ctx))
			}
		}
	}
	codestr := string(code)
	lexer := NewLexer(&codePath, &codestr)
	parser := MakeParser(lexer)
	ast := parser.Parse()
    var context *Context
    if name == "*" {
        context = ctx
    } else {
		context = MakeContext(ctx, ast.S)
		mt := Type(&Obj{context, []string{"module"}})
		ctx.set(name, &mt)
		if ast.T == "body" {
			ast.V = name
		}
    }
	res.reg(ExecAst(ast, context))
	if res.has() {
		return res
	}
	return res.suss(&Null{}, ctx)
}
