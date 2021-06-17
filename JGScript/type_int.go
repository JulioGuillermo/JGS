package JGScript

import (
	"fmt"
	"math"
)

// INT
type Int struct{ val int64 }

func (p *Int) GetT() string {
	return "int"
}

func (p *Int) str() string {
	return fmt.Sprint(p.val)
}

func (p *Int) operator(op string, t *Type, c *Context, s, e *BPos) *Result {
	res := MkRes()
	if op == "_add" && (*t).GetT() == "string" {
		return res.suss(&Str{p.str() + (*t).str()}, c)
	} else {
		switch op {
		case "_inc":
			p.val++
			return res.suss(p, c)
		case "_dec":
			p.val--
			return res.suss(p, c)
		case "_pos":
			return res.suss(p, c)
		case "_neg":
			p.val = -p.val
			return res.suss(p, c)
		case "_eq":
			if (*t).GetT() != "int" && (*t).GetT() != "float" {
				return res.suss(&Bool{false}, c)
			}
		case "_not_eq":
			if (*t).GetT() != "int" && (*t).GetT() != "float" {
				return res.suss(&Bool{true}, c)
			}
		}
	}
	switch (*t).GetT() {
	case "byte":
		i, _ := (*t).(*Byte)
		switch op {
		case "_add":
			return res.suss(&Int{p.val + int64(i.val)}, c)
		case "_sub":
			return res.suss(&Int{p.val - int64(i.val)}, c)
		case "_mul":
			return res.suss(&Int{p.val * int64(i.val)}, c)
		case "_div":
			return res.suss(&Int{p.val / int64(i.val)}, c)
		case "_mod":
			return res.suss(&Int{p.val % int64(i.val)}, c)
		case "_pow":
			return res.suss(&Int{int64(math.Pow(float64(p.val), float64(i.val)))}, c)
		case "_shift_left":
            if i.val < 0 {
			    return res.suss(&Int{p.val >> int64(-i.val)}, c)
            }
			return res.suss(&Int{p.val << int64(i.val)}, c)
		case "_shift_right":
            if i.val < 0 {
			    return res.suss(&Int{p.val << int64(-i.val)}, c)
            }
			return res.suss(&Int{p.val >> int64(i.val)}, c)
		case "_add_eq":
			p.val += int64(i.val)
			return res.suss(p, c)
		case "_sub_eq":
			p.val -= int64(i.val)
			return res.suss(p, c)
		case "_mul_eq":
			p.val *= int64(i.val)
			return res.suss(p, c)
		case "_div_eq":
			p.val /= int64(i.val)
			return res.suss(p, c)
		case "_mod_eq":
			p.val = p.val % int64(i.val)
			return res.suss(p, c)
		case "_pow_eq":
			p.val = int64(math.Pow(float64(p.val), float64(i.val)))
			return res.suss(p, c)
		case "_eq":
			return res.suss(&Bool{p.val == int64(i.val)}, c)
		case "_not_eq":
			return res.suss(&Bool{p.val != int64(i.val)}, c)
		case "_lt":
			return res.suss(&Bool{p.val < int64(i.val)}, c)
		case "_lte":
			return res.suss(&Bool{p.val <= int64(i.val)}, c)
		case "_gt":
			return res.suss(&Bool{p.val > int64(i.val)}, c)
		case "_gte":
			return res.suss(&Bool{p.val >= int64(i.val)}, c)
        case "_and":
            return res.suss(&Int{p.val & int64(i.val)}, c)
        case "_or":
            return res.suss(&Int{p.val | int64(i.val)}, c)
		}
	case "int":
		i, _ := (*t).(*Int)
		switch op {
		case "_add":
			return res.suss(&Int{p.val + i.val}, c)
		case "_sub":
			return res.suss(&Int{p.val - i.val}, c)
		case "_mul":
			return res.suss(&Int{p.val * i.val}, c)
		case "_div":
			return res.suss(&Int{p.val / i.val}, c)
		case "_mod":
			return res.suss(&Int{p.val % i.val}, c)
		case "_pow":
			return res.suss(&Int{int64(math.Pow(float64(p.val), float64(i.val)))}, c)
		case "_shift_left":
            if i.val < 0 {
			    return res.suss(&Int{p.val >> -i.val}, c)
            }
			return res.suss(&Int{p.val << i.val}, c)
		case "_shift_right":
            if i.val < 0 {
			    return res.suss(&Int{p.val << -i.val}, c)
            }
			return res.suss(&Int{p.val >> i.val}, c)
		case "_add_eq":
			p.val += i.val
			return res.suss(p, c)
		case "_sub_eq":
			p.val -= i.val
			return res.suss(p, c)
		case "_mul_eq":
			p.val *= i.val
			return res.suss(p, c)
		case "_div_eq":
			p.val /= i.val
			return res.suss(p, c)
		case "_mod_eq":
			p.val = p.val % i.val
			return res.suss(p, c)
		case "_pow_eq":
			p.val = int64(math.Pow(float64(p.val), float64(i.val)))
			return res.suss(p, c)
		case "_eq":
			return res.suss(&Bool{p.val == i.val}, c)
		case "_not_eq":
			return res.suss(&Bool{p.val != i.val}, c)
		case "_lt":
			return res.suss(&Bool{p.val < i.val}, c)
		case "_lte":
			return res.suss(&Bool{p.val <= i.val}, c)
		case "_gt":
			return res.suss(&Bool{p.val > i.val}, c)
		case "_gte":
			return res.suss(&Bool{p.val >= i.val}, c)
        case "_and":
            return res.suss(&Int{p.val & i.val}, c)
        case "_or":
            return res.suss(&Int{p.val | i.val}, c)
		}
	case "float":
		i, _ := (*t).(*Flt)
		switch op {
		case "_add":
			return res.suss(&Flt{float64(p.val) + i.val}, c)
		case "_sub":
			return res.suss(&Flt{float64(p.val) - i.val}, c)
		case "_mul":
			return res.suss(&Flt{float64(p.val) * i.val}, c)
		case "_div":
			return res.suss(&Flt{float64(p.val) / i.val}, c)
		case "_mod":
			return res.suss(&Int{p.val % int64(i.val)}, c)
		case "_pow":
			return res.suss(&Flt{math.Pow(float64(p.val), i.val)}, c)
		case "_shift_left":
            if i.val < 0 {
			    return res.suss(&Int{p.val >> int64(-i.val)}, c)
            }
			return res.suss(&Int{p.val << int64(i.val)}, c)
		case "_shift_right":
            if i.val < 0 {
			    return res.suss(&Int{p.val << int64(-i.val)}, c)
            }
			return res.suss(&Int{p.val >> int64(i.val)}, c)
		case "_add_eq":
			p.val += int64(i.val)
			return res.suss(p, c)
		case "_sub_eq":
			p.val -= int64(i.val)
			return res.suss(p, c)
		case "_mul_eq":
			p.val *= int64(i.val)
			return res.suss(p, c)
		case "_div_eq":
			p.val /= int64(i.val)
			return res.suss(p, c)
		case "_mod_eq":
			p.val = p.val % int64(i.val)
			return res.suss(p, c)
		case "_pow_eq":
			p.val = int64(math.Pow(float64(p.val), i.val))
			return res.suss(p, c)
		case "_eq":
			return res.suss(&Bool{float64(p.val) == i.val}, c)
		case "_not_eq":
			return res.suss(&Bool{float64(p.val) != i.val}, c)
		case "_lt":
			return res.suss(&Bool{float64(p.val) < i.val}, c)
		case "_lte":
			return res.suss(&Bool{float64(p.val) <= i.val}, c)
		case "_gt":
			return res.suss(&Bool{float64(p.val) > i.val}, c)
		case "_gte":
			return res.suss(&Bool{float64(p.val) >= i.val}, c)
        case "_and":
            return res.suss(&Int{p.val & int64(i.val)}, c)
        case "_or":
            return res.suss(&Int{p.val | int64(i.val)}, c)
		}
	}
	return res.fail(MakeRTError("Invalid operation for "+(*p).str()+" and "+(*t).str(), s, e, c))
}

func (p *Int) getBool() bool {
	return p.val != 0
}

func (p *Int) getByte() byte {
    return byte(p.val)
}

func (p *Int) getInt() int64 {
	return p.val
}

func (p *Int) GetInt() *int64 {
    return &p.val
}

func (p *Int) getFloat() float64 {
	return float64(p.val)
}

func (p *Int) equals(t *Type) bool {
	if (*t).GetT() == "byte" {
		i, _ := (*t).(*Byte)
		return p.val == int64(i.val)
	}
	if (*t).GetT() == "int" {
		i, _ := (*t).(*Int)
		return p.val == i.val
	}
	if (*t).GetT() == "float" {
		i, _ := (*t).(*Flt)
		return float64(p.val) == i.val
	}
	return false
}

func (p *Int) GetMember(a *AST, c *Context) *Result {
    return MkRes().fail(MakeError("MemberAccessError", "Type " + p.GetT() + " has not member " + a.V, a.S, a.E, c))
}
