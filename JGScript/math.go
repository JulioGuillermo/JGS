package JGScript

import "math"

func strEq(s string, ss []string) bool {
    for _, e := range ss {
        if e == s {
            return true
        }
    }
    return false
}

func MathFuns(f, arg *AST, c *Context) *Result {
    fn := f.V
    res := MkRes()
    if strEq(fn, []string{
        "acos",
        "acosh",
        "asin",
        "asinh",
        "atan",
        "atanh",
        "cos",
        "cosh",
        "sin",
        "sinh",
        "tan",
        "tanh",
        "sqrt",
        "ceil",
        "floor",
        "pow10",
        "round",
        "roundToEven",
        "log",
        "log2",
        "log10",
        "abs",
        "exp",
        "exp2",
        "expm1",
        "trunc",
        "cbrt"}) {
		if len(arg.C) < 1 {
			return res.fail(MakeRTError("Arguments required: 1", arg.S, arg.E, c))
		}
		n := res.reg(ExecAst(arg.C[0], c))
		if res.has() {
			return res
		}
        var num float64
		switch (*n).GetT() {
        case "int":
            num = float64((*n).(*Int).val)
        case "float":
            num = (*n).(*Flt).val
        default:
			return res.fail(MakeRTError("Arguments must be numbers.", arg.C[0].S, arg.C[0].E, c))
		}
        return res.suss(&Flt{arg1mathfuns(fn, num)}, c)
    } else if strEq(fn, []string{
        "atan2",
        "copysign",
        "pow",
        "mod",
        "dim"}) {
		if len(arg.C) < 2 {
			return res.fail(MakeRTError("Arguments required: 2", arg.S, arg.E, c))
		}
		n := res.reg(ExecAst(arg.C[0], c))
		if res.has() {
			return res
		}
        var num float64
		switch (*n).GetT() {
        case "int":
            num = float64((*n).(*Int).val)
        case "float":
            num = (*n).(*Flt).val
        default:
			return res.fail(MakeRTError("Arguments must be numbers.", arg.C[0].S, arg.C[0].E, c))
		}
		n2 := res.reg(ExecAst(arg.C[1], c))
		if res.has() {
			return res
		}
        var num2 float64
		switch (*n2).GetT() {
        case "int":
            num2 = float64((*n2).(*Int).val)
        case "float":
            num2 = (*n2).(*Flt).val
        default:
			return res.fail(MakeRTError("Arguments must be numbers.", arg.C[1].S, arg.C[1].E, c))
		}
        return res.suss(&Flt{arg2mathfuns(fn, num, num2)}, c)
    }
    return nil
}

func arg1mathfuns(f string, n float64) float64 {
    switch f {
    case "acos":
        return math.Acos(n)
    case "acosh":
        return math.Acosh(n)
    case "asin":
        return math.Asin(n)
    case "asinh":
        return math.Asinh(n)
    case "atan":
        return math.Atan(n)
    case "cos":
        return math.Cos(n)
    case "cosh":
        return math.Cosh(n)
    case "sin":
        return math.Sin(n)
    case "sinh":
        return math.Sinh(n)
    case "tan":
        return math.Tan(n)
    case "tanh":
        return math.Tanh(n)
    case "sqrt":
        return math.Sqrt(n)
    case "ceil":
        return math.Ceil(n)
    case "floor":
        return math.Floor(n)
    case "pow10":
        return math.Pow10(int(n))
    case "round":
        return math.Round(n)
    case "roundToEven":
        return math.RoundToEven(n)
    case "log":
        return math.Log(n)
    case "log2":
        return math.Log2(n)
    case "log10":
        return math.Log10(n)
    case "abs":
        return math.Abs(n)
    case "exp":
        return math.Exp(n)
    case "exp2":
        return math.Exp2(n)
    case "expm1":
        return math.Expm1(n)
    case "trunc":
        return math.Trunc(n)
    case "cbrt":
        return math.Cbrt(n)
    }
    return 0
}

func arg2mathfuns(f string, x, y float64) float64 {
    switch f {
    case "atan2":
        return math.Atan2(x, y)
    case "copysign":
        return math.Copysign(x, y)
    case "pow":
        return math.Pow(x, y)
    case "mod":
        return math.Mod(x, y)
    case "dim":
        return math.Dim(x, y)
    }
    return 0
}
