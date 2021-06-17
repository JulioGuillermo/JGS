package JGScript

import (
	"fmt"
	"os"
)

func parserError(t *Token) {
	fmt.Println("File: " + *t.s.Fn + ", Line: " + fmt.Sprint(t.s.Row) + ", Column: " + fmt.Sprint(t.s.Col) + "\nInvalidSyntaxError: " + t.v)
	fmt.Println(TextArrow(t.s, t.e))
	os.Exit(2)
}

type Parser struct {
	lexer *Lexer
	tok   *Token
}

func MakeParser(l *Lexer) *Parser {
	return &Parser{l, l.NextToken()}
}

func (p *Parser) next() {
	p.tok = p.lexer.NextToken()
}

func (p *Parser) nextC(t string) {
	if p.tok.t != t {
		parserError(p.tok)
	}
	p.tok = p.lexer.NextToken()
}

func (p *Parser) skipEnd() {
	for p.tok.t == "SEMI" || p.tok.t == "COMMA" || p.tok.t == "LN" {
		p.next()
	}
}

func (p *Parser) tokType(s ...string) bool {
	for _, t := range s {
		if p.tok.t == t {
			return true
		}
	}
	return false
}

func (p *Parser) Parse() *AST {
	s := p.tok.s
	e := p.tok.e
	var asts []*AST
	p.skipEnd()
	for p.tok.t != "EOF" {
		asts = append(asts, p.parseExpr())
		p.skipEnd()
	}
	if len(asts) > 0 {
		s = asts[0].S
		e = asts[len(asts)-1].E
	}
	return &AST{"body", "", asts, s, e}
}

func (p *Parser) parseExpr() *AST {
	return p.parseAssign()
}

func (p *Parser) parseAssign() *AST {
	left := p.parseArray()
	for p.tok.t == "EQ" {
		p.next()
		right := p.parseArray()
		left = &AST{T: "assign", C: []*AST{left, right}, S: left.S, E: right.E}
	}
	return left
}

func (p *Parser) parseArray() *AST {
	base := p.parseAnd()
	if p.tok.t == "COMMA" {
		p.next()
		p.skipEnd()
		args := []*AST{base, p.parseAnd()}
		for p.tok.t == "COMMA" {
			p.next()
			p.skipEnd()
			args = append(args, p.parseAnd())
		}
		return &AST{T: "array", C: args, S: base.S, E: args[len(args)-1].E}
	}
	return base
}

func (p *Parser) parseArgs() *AST {
	s := p.tok.s
	args := []*AST{}
	p.skipEnd()
	if p.tok.t != "RP" {
		arg := p.parseArgArray()
		if arg.T == "array" {
			args = arg.C
		} else {
			args = append(args, arg)
		}
	}
	return &AST{T: "args", C: args, S: s, E: p.tok.e}
}

func (p *Parser) parseArgArray() *AST {
	p.skipEnd()
	base := p.parseArgAssign()
	if p.tok.t == "COMMA" {
		p.skipEnd()
		args := []*AST{base, p.parseArgAssign()}
		for p.tok.t == "COMMA" {
			p.skipEnd()
			args = append(args, p.parseArgAssign())
		}
		p.skipEnd()
		return &AST{T: "array", C: args, S: base.S, E: args[len(args)-1].E}
	}
	return base
}

func (p *Parser) parseArgAssign() *AST {
	left := p.parseAnd()
	for p.tok.t == "EQ" {
		p.next()
		right := p.parseAnd()
		left = &AST{T: "assign", C: []*AST{left, right}, S: left.S, E: right.E}
	}
	return left
}

func (p *Parser) parseAnd() *AST {
	left := p.parseOr()
	for p.tokType("_and") {
		op := p.tok.t
		p.next()
		right := p.parseOr()
		left = &AST{"bin_operator", op, []*AST{left, right}, left.S, right.E}
	}
	return left
}

func (p *Parser) parseOr() *AST {
	left := p.parseCmp()
	for p.tokType("_or") {
		op := p.tok.t
		p.next()
		right := p.parseCmp()
		left = &AST{"bin_operator", op, []*AST{left, right}, left.S, right.E}
	}
	return left
}

func (p *Parser) parseCmp() *AST {
	left := p.parseRange()
	for p.tokType("_eq", "_not_eq", "_lt", "_lte", "_gt", "_gte") {
		op := p.tok.t
		p.next()
		right := p.parseRange()
		left = &AST{"bin_operator", op, []*AST{left, right}, left.S, right.E}
	}
	return left
}

func (p *Parser) parseRange() *AST {
	a := p.parseAddSub()
	if p.tokType("COLON") {
        e := p.tok.e
		p.next()
        if !p.tokType("EOF", "DOT", "DDD", "COMMA", "SEMI", "COLON", "RP", "RB", "RS", "EQ", "_and", "_or", "_not", "_eq", "_not_eq", "_lt", "_lte", "_gt", "_gte") {
            b := p.parseAddSub()
		    return &AST{"range", "", []*AST{a, b}, a.S, b.E}
        }
		return &AST{"range_start", "", []*AST{a}, a.S, e}
	}
	return a
}

func (p *Parser) parseAddSub() *AST {
	left := p.parseMulDiv()
	for p.tokType("_add", "_add_eq", "_sub", "_sub_eq") {
		op := p.tok.t
		p.next()
		right := p.parseMulDiv()
		left = &AST{"bin_operator", op, []*AST{left, right}, left.S, right.E}
	}
	return left
}

func (p *Parser) parseMulDiv() *AST {
	left := p.parseMod()
	for p.tokType("_mul", "_mul_eq", "_div", "_div_eq") {
		op := p.tok.t
		p.next()
		right := p.parseMod()
		left = &AST{"bin_operator", op, []*AST{left, right}, left.S, right.E}
	}
	return left
}

func (p *Parser) parseMod() *AST {
	left := p.parsePow()
	for p.tokType("_mod", "_mod_eq") {
		op := p.tok.t
		p.next()
		right := p.parsePow()
		left = &AST{"bin_operator", op, []*AST{left, right}, left.S, right.E}
	}
	return left
}

func (p *Parser) parsePow() *AST {
	left := p.parseShift()
	for p.tokType("_pow", "_pow_eq") {
		op := p.tok.t
		p.next()
		right := p.parseShift()
		left = &AST{"bin_operator", op, []*AST{left, right}, left.S, right.E}
	}
	return left
}

func (p *Parser) parseShift() *AST {
    left := p.parseIncDec()
    for p.tokType("_shift_left", "_shift_right") {
        op := p.tok.t
        p.next()
        right := p.parseIncDec()
        left = &AST {"bin_operator", op, []*AST{left, right}, left.S, right.E}
    }
    return left
}

func (p *Parser) parseIncDec() *AST {
	left := p.parseArrow()
	for p.tokType("_inc", "_dec") {
		op := p.tok.t
		p.next()
		left = &AST{"unary_operator", op, []*AST{left}, left.S, left.E}
	}
	return left
}

func (p *Parser) parseArrow() *AST {
	left := p.parseLPB()
	for p.tokType("_left_arrow", "_right_arrow") {
		op := p.tok.t
		p.next()
		right := p.parseLPB()
		left = &AST{T: "bin_operator", V: op, C: []*AST{left, right}, S: left.S, E: right.E}
	}
	return left
}

func (p *Parser) parseLPB() *AST {
	left := p.parseAtom()
	for p.tokType("LP", "LS", "DOT") {
		if p.tok.t == "LP" {
			left = p.parseCall(left)
		} else if p.tok.t == "DOT" {
            left = p.parseMember(left)
        } else {
			left = p.parseSquert(left)
		}
	}
	return left
}

func (p *Parser) parseCall(a *AST) *AST {
	p.next()
	args_ast := p.parseArgs()
	e := p.tok.e
	p.nextC("RP")
	return &AST{T: "_call", C: []*AST{a, args_ast}, S: a.S, E: e}
}

func (p *Parser) parseSquert(a *AST) *AST {
	p.next()
	args_ast := p.parseArgs()
	e := p.tok.e
	p.nextC("RS")
	return &AST{T: "_index", C: []*AST{a, args_ast}, S: a.S, E: e}
}

func (p *Parser) parseMember(a *AST) *AST {
    p.next()
    m := &AST{T: "access", V: p.tok.v, S: p.tok.s, E: p.tok.e}
    p.nextC("WORD")
    return &AST {T: "member", C: []*AST {a, m}, S: a.S, E: m.E}
}

func (p *Parser) parseAtom() *AST {
	switch p.tok.t {
	case "_left_arrow":
		s := p.tok.s
		p.next()
		exp := p.parseExpr()
		return &AST{T: "body_ret", C: []*AST{exp}, S: s, E: exp.E}
	case "LP":
		return p.parseAtomicParen()
	case "LS":
		return p.parseAtomicSquert()
	case "LB":
		ast := p.parseBody()
		ast.T = "obj"
		return ast
	case "DDD":
		s := p.tok.s
		p.next()
		return &AST{T: "DDD", S: s, E: p.tok.e}
    case "COLON":
        s := p.tok.s
        e := p.tok.e
        p.next()
        if !p.tokType("EOF", "DOT", "DDD", "COMMA", "SEMI", "COLON", "RP", "RB", "RS", "EQ", "_and", "_or", "_not", "_eq", "_not_eq", "_lt", "_lte", "_gt", "_gte") {
            b := p.parseAddSub()
		    return &AST{"range_end", "", []*AST{b}, s, b.E}
        }
		return &AST{"range_empty", "", []*AST{}, s, e}
	case "_not":
		s := p.tok.s
		p.next()
		a := []*AST{p.parseExpr()}
		return &AST{T: "unary_operator", V: "_not", C: a, S: s, E: p.tok.e}
	case "_inc":
		fallthrough
	case "_dec":
		fallthrough
	case "_add":
		s := p.tok.s
		p.next()
		a := []*AST{p.parseAddSub()}
		return &AST{T: "unary_operator", V: "_pos", C: a, S: s, E: p.tok.e}
	case "_sub":
		s := p.tok.s
		p.next()
		a := []*AST{p.parseAddSub()}
		return &AST{T: "unary_operator", V: "_neg", C: a, S: s, E: p.tok.e}
	case "INT":
		a := &AST{T: "int", V: p.tok.v, S: p.tok.s, E: p.tok.e}
		p.next()
		return a
	case "HEX":
		a := &AST{T: "hex", V: p.tok.v, S: p.tok.s, E: p.tok.e}
		p.next()
		return a
	case "OCT":
		a := &AST{T: "oct", V: p.tok.v, S: p.tok.s, E: p.tok.e}
		p.next()
		return a
	case "BIN":
		a := &AST{T: "bin", V: p.tok.v, S: p.tok.s, E: p.tok.e}
		p.next()
		return a
	case "FLOAT":
		a := &AST{T: "float", V: p.tok.v, S: p.tok.s, E: p.tok.e}
		p.next()
		return a
	case "STRING":
		a := &AST{T: "string", V: p.tok.v, S: p.tok.s, E: p.tok.e}
		p.next()
		return a
	case "WORD":
		return p.parseWord()
	}
	parserError(p.tok)
	return nil
}

func (p *Parser) parseAtomicParen() *AST {
	s := p.tok.s
	p.next()
	if p.tok.t == "RP" {
		args_ast := &AST{T: "args", C: []*AST{}, S: s, E: p.tok.e}
		p.next()
		at := "fun"
		if p.tok.v == "async" {
			at = "afun"
			p.next()
		}
		body := p.parseBody()
		return &AST{T: at, C: []*AST{args_ast, body}, S: s, E: body.E}
	}
	args_ast := p.parseArgs()
	e := p.tok.e
	p.nextC("RP")
	p.skipEnd()
	if p.tok.t == "LB" || p.tok.v == "async" {
		at := "fun"
		if p.tok.v == "async" {
			at = "afun"
			p.next()
		}
		body := p.parseBody()
		return &AST{T: at, C: []*AST{args_ast, body}, S: s, E: body.E}
	}
	if len(args_ast.C) == 1 {
		return args_ast.C[0]
	}
	return &AST{T: "list", C: args_ast.C, S: s, E: e}
}

func (p *Parser) parseAtomicSquert() *AST {
	s := p.tok.s
	p.next()
	args := []*AST{}
	p.skipEnd()
	if p.tok.t != "RS" {
		args = p.parseArgs().C
	}
	e := p.tok.e
	p.nextC("RS")
	return &AST{T: "list", C: args, S: s, E: e}
}

func (p *Parser) parseWord() *AST {
	s := p.tok.s
	e := p.tok.e
	v := p.tok.v
	p.nextC("WORD")
	switch v {
	case "null":
		return &AST{T: "null", S: s, E: e}
	case "true":
		fallthrough
	case "false":
		return &AST{T: "bool", V: v, S: s, E: e}
	case "ret":
        for p.tok.t == "LN" {
            p.next()
        }
		return &AST{T: "ret", C: []*AST{p.parseExpr()}, S: s, E: e}
	case "if":
		return p.parseIf()
	case "switch":
		return p.parseSwitch()
	case "for":
		return p.parseFor()
	case "break":
		return &AST{T: "break", S: s, E: e}
	case "try":
		return p.parseTry()
		/*case "class":
		  return p.parseClass()*/
    case "import":
        pa := p.parseExpr()
        if p.tok.t == "WORD" && p.tok.v == "as" {
            p.next()
            aa := p.parseExpr()
            return &AST {T: "import", C: []*AST{pa, aa}, S: s, E: aa.E}
        }
        return &AST {T: "import", C: []*AST{pa}, S: s, E: pa.E}
	}
	return &AST{T: "access", V: v, S: s, E: e}
}

func (p *Parser) parseTry() *AST {
	s := p.tok.s
	body := p.parseBody()
	p.skipEnd()
	if p.tok.t != "WORD" || p.tok.v != "catch" {
		parserError(p.tok)
	}
	p.next()
	p.skipEnd()
	args := p.parseArgs()
	p.skipEnd()
	catch := p.parseBody()
	args_ast := []*AST{body, args, catch}
	return &AST{T: "try", C: args_ast, S: s, E: catch.E}
}

func (p *Parser) parseIf() *AST {
	s := p.tok.s
	p.skipEnd()
	cond := p.parseExpr()
	p.skipEnd()
	args := []*AST{cond, p.parseBody()}
	p.skipEnd()
	if p.tok.t == "WORD" {
		if p.tok.v == "else" {
			p.next()
			p.skipEnd()
			args = append(args, p.parseBody())
		} else if p.tok.v == "elif" {
			p.next()
			p.skipEnd()
			args = append(args, p.parseIf())
		}
	}
	return &AST{T: "if", C: args, S: s, E: p.tok.e}
}

func (p *Parser) parseSwitch() *AST {
	s := p.tok.s
	p.skipEnd()
	cond := p.parseExpr()
	p.skipEnd()
	p.nextC("LB")
	args := []*AST{}
	p.skipEnd()
	for p.tok.t == "WORD" && p.tok.v == "case" {
		p.next()
		p.skipEnd()
		cas := p.parseExpr()
		p.skipEnd()
		if p.tok.t != "WORD" || (p.tok.v != "case" && p.tok.v != "default") {
			bod := p.parseBody()
			args = append(args, &AST{T: "case", C: []*AST{cas, bod}, S: cas.S, E: bod.E})
		} else {
			args = append(args, &AST{T: "case", C: []*AST{cas}, S: cas.S, E: cas.E})
		}
		p.skipEnd()
	}
	p.skipEnd()
	if p.tok.t == "WORD" && p.tok.v == "default" {
		p.next()
		bod := p.parseBody()
		args = append(args, &AST{T: "default", C: []*AST{bod}, S: bod.S, E: bod.E})
	}
	e := p.tok.e
	p.nextC("RB")
	return &AST{T: "switch", C: []*AST{cond, &AST{C: args}}, S: s, E: e}
}

func (p *Parser) parseFor() *AST {
	s := p.tok.s
	args := []*AST{}
	p.skipEnd()
	if p.tok.t != "LB" {
		args = append(args, p.parseExpr())
		p.skipEnd()
		if p.tok.t != "LB" {
			if p.tok.t == "WORD" && p.tok.v == "in" {
				p.next()
				p.skipEnd()
				args = append(args, p.parseExpr())
			} else {
				args = append(args, p.parseExpr())
				p.skipEnd()
				args = append(args, p.parseExpr())
			}
		}
	}
	body := p.parseBody()
	return &AST{T: "for", C: []*AST{&AST{C: args}, body}, S: s, E: body.E}
}

func (p *Parser) parseBody() *AST {
	s := p.tok.s
	e := p.tok.e
	var asts []*AST
	p.skipEnd()
	if p.tok.t == "LB" {
		p.next()
		p.skipEnd()
		for p.tok.t != "RB" {
			if p.tok.t == "EOF" {
				parserError(p.tok)
			}
			asts = append(asts, p.parseExpr())
			p.skipEnd()
		}
		e = p.tok.e
		p.nextC("RB")
	} else {
		asts = append(asts, p.parseExpr())
		e = p.tok.e
		p.skipEnd()
	}
	return &AST{T: "body", V: "", C: asts, S: s, E: e}
}
