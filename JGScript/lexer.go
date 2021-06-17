package JGScript

import (
	"fmt"
	"os"
)

type Lexer struct {
	fn   *string
	code *string
	c    rune
	pos  *Position
}

func NewLexer(fn, code *string) *Lexer {
	return &Lexer{fn, code, 0, NewPosition(fn, code)}
}

func (p *Lexer) NextToken() *Token {
	p.c = p.pos.Get()
	p.skipAll()
	start := p.pos.GetP()
	if p.c != 0 {
        if p.c == '#' || (p.c == '0' && p.pos.GetN() == 'x') {
            return p.lexHex()
        }
        if p.c == '@' {
            return p.lexOct()
        }
        if p.c == '$' {
            return p.lexBin()
        }
		if CIsAlpha(p.c) {
			return p.lexWord()
		}
		if CIsDigit(p.c) {
			return p.lexNumber()
		}
		if p.c == '"' || p.c == '\'' {
			return p.lexString()
		}
		switch p.c {
		case '=':
			if p.pos.GetN() == '=' {
				return p.makeToken("_eq", "==")
			}
			return p.makeToken("EQ", "=")

		case '+':
			if p.pos.GetN() == '+' {
				return p.makeToken("_inc", "++")
			}
			if p.pos.GetN() == '=' {
				return p.makeToken("_add_eq", "+=")
			}
			return p.makeToken("_add", "+")
		case '-':
			if p.pos.GetN() == '>' {
				return p.makeToken("_right_arrow", "->")
			}
			if p.pos.GetN() == '-' {
				return p.makeToken("_dec", "--")
			}
			if p.pos.GetN() == '=' {
				return p.makeToken("_sub_eq", "-=")
			}
			return p.makeToken("_sub", "-")

		case '*':
			if CIsAlpha(p.pos.GetN()) || p.pos.GetN() == '*' {
				return p.lexWord()
			}
			if p.pos.GetN() == '=' {
				return p.makeToken("_mul_eq", "*=")
			}
			return p.makeToken("_mul", "*")
		case '/':
			if p.pos.GetN() == '=' {
				return p.makeToken("_div_eq", "/=")
			}
			return p.makeToken("_div", "/")

		case '%':
			if p.pos.GetN() == '=' {
				return p.makeToken("_mod_eq", "%=")
			}
			return p.makeToken("_mod", "%")

		case '^':
			if p.pos.GetN() == '=' {
				return p.makeToken("_pow_eq", "^=")
			}
			return p.makeToken("_pow", "^")

		case '<':
			if p.pos.GetN() == '-' {
				return p.makeToken("_left_arrow", "<-")
			}
            if p.pos.GetN() == '<' {
                return p.makeToken("_shift_left", "<<")
            }
			if p.pos.GetN() == '=' {
				return p.makeToken("_lte", "<=")
			}
			return p.makeToken("_lt", "<")
		case '>':
			if p.pos.GetN() == '=' {
				return p.makeToken("_gte", ">=")
			}
            if p.pos.GetN() == '>' {
                return p.makeToken("_shift_right", ">>")
            }
			return p.makeToken("_gt", ">")

		case '&':
			for p.c == '&' {
				p.c = p.pos.Next()
			}
			return &Token{"_and", "&", start, p.pos.GetP()}
		case '|':
			for p.c == '|' {
				p.c = p.pos.Next()
			}
			return &Token{"_or", "|", start, p.pos.GetP()}
		case '!':
			if p.pos.GetN() == '=' {
				return p.makeToken("_not_eq", "!=")
			}
			return p.makeToken("_not", "!")

		case '(':
			return p.makeToken("LP", "(")
		case ')':
			return p.makeToken("RP", ")")

		case '[':
			return p.makeToken("LS", "[")
		case ']':
			return p.makeToken("RS", "]")

		case '{':
			return p.makeToken("LB", "{")
		case '}':
			return p.makeToken("RB", "}")

		case '.':
			if p.pos.GetN() == '.' && p.pos.GetNN() == '.' {
				return p.makeToken("DDD", "...")
			}
			if CIsDigit(p.pos.GetN()) {
				return p.lexNumber()
			}
			return p.makeToken("DOT", ".")
		case ',':
			return p.makeToken("COMMA", ",")
		case ';':
			return p.makeToken("SEMI", ";")
		case ':':
			return p.makeToken("COLON", ":")
        case 13:
		    fallthrough
        case 10:
			return p.makeToken("LN", "\n")
		}
		fmt.Println("File: " + *p.fn + ", Line: " + fmt.Sprint(p.pos.row) + ", Column: " + fmt.Sprint(p.pos.col) + " => Unexpected character: " + string(p.c))
		os.Exit(1)
	}
	return &Token{"EOF", "", start, p.pos.GetP()}
}

func (p *Lexer) skipWhiteSpace() {
	for p.c == ' ' || p.c == '\t' {
		p.c = p.pos.Next()
	}
}

func (p *Lexer) skipLineComment() {
	if p.c == '/' && p.pos.GetN() == '/' {
		p.pos.Next()
		p.c = p.pos.Next()
		for p.c != '\n' && p.c != 0 {
			p.c = p.pos.Next()
		}
		p.c = p.pos.Next()
	}
}

func (p *Lexer) skipMultiLineComment() {
	if p.c == '/' && p.pos.GetN() == '*' {
		p.pos.Next()
		p.c = p.pos.Next()
		for p.c != '*' || p.pos.GetN() != '/' {
			if p.pos.index >= len(*p.code) {
				break
			}
			p.c = p.pos.Next()
		}
		p.pos.Next()
		p.c = p.pos.Next()
	}
}

func (p *Lexer) skipAll() {
	for p.c != 0 && ((p.c == ' ' || p.c == '\t') || (p.c == '/' && (p.pos.GetN() == '/' || p.pos.GetN() == '*'))) {
		p.skipMultiLineComment()
		p.skipLineComment()
		p.skipWhiteSpace()
	}
}

func (p *Lexer) makeToken(t, v string) *Token {
	start := p.pos.GetP()
	for i := 0; i < len(v); i++ {
		p.pos.Next()
	}
	p.c = p.pos.Get()
	return &Token{t, v, start, p.pos.GetP()}
}

func (p *Lexer) lexWord() *Token {
	w := ""
	start := p.pos.GetP()
	for p.c == '*' {
		w += string('*')
		p.c = p.pos.Next()
	}
	for CIsAlphaDigit(p.c) {
		w += string(p.c)
		p.c = p.pos.Next()
	}
	return &Token{"WORD", w, start, p.pos.GetP()}
}

func (p *Lexer) lexNumber() *Token {
	n := ""
	d := false
    e := false
	start := p.pos.GetP()
	for CIsDigit(p.c) || p.c == '.' || p.c == '_' {
		if p.c != '_' {
			if p.c == '.' {
				if d {
					return &Token{"FLOAT", n, start, p.pos.GetP()}
				}
				d = true
			}
			n += string(p.c)
		}
		p.c = p.pos.Next()
	}
    if !e && (p.c == 'e' || p.c == 'E') && (p.pos.GetN() == '-' || p.pos.GetN() == '+') && CIsDigit(p.pos.GetNN()) {
        e = true
        n += "e"
        p.c = p.pos.Next()
        n += string(p.c)
        p.c = p.pos.Next()
	    for CIsDigit(p.c) || p.c == '_' {
		    if p.c != '_' {
			    n += string(p.c)
		    }
		    p.c = p.pos.Next()
	    }
    }
	if d || e {
	    return &Token{"FLOAT", n, start, p.pos.GetP()}
	}
	return &Token{"INT", n, start, p.pos.GetP()}
}

func (p *Lexer) lexHex() *Token {
    n := ""
    start := p.pos.GetP()
    if p.c == '0' && p.pos.GetN() == 'x' {
        p.pos.Next()
        p.c = p.pos.Next()
    } else if p.c == '#' {
        p.c = p.pos.Next()
    }
    for CIsHDigit(p.c) || p.c == '_' {
        if p.c != '_' {
            n += string(p.c)
        }
        p.c = p.pos.Next()
    }
    return &Token{"HEX", n, start, p.pos.GetP()}
}

func (p *Lexer) lexOct() *Token {
    n := ""
    start := p.pos.GetP()
    p.c = p.pos.Next()
    for CIsODigit(p.c) || p.c == '_' {
        if p.c != '_' {
            n += string(p.c)
        }
        p.c = p.pos.Next()
    }
    return &Token{"OCT", n, start, p.pos.GetP()}
}

func (p *Lexer) lexBin() *Token {
    n := ""
    start := p.pos.GetP()
    p.c = p.pos.Next()
    for p.c == '0' || p.c == '1' || p.c == '_' {
        if p.c != '_' {
            n += string(p.c)
        }
        p.c = p.pos.Next()
    }
    return &Token{"BIN", n, start, p.pos.GetP()}
}

func (p *Lexer) lexString() *Token {
	t := p.c
	start := p.pos.GetP()
	p.c = p.pos.Next()
	s := ""
	c := false
	for (p.c != t || c) && p.c != 0 {
		if c {
			if p.c == 'n' {
				s += "\n"
			} else if p.c == 't' {
				s += "\t"
			} else if p.c == 'r' {
				s += "\r"
			} else {
				s += string(p.c)
			}
			c = false
		} else if p.c == '\\' {
			c = true
		} else {
			s += string(p.c)
		}
		p.c = p.pos.Next()
	}
	p.c = p.pos.Next()
	return &Token{"STRING", s, start, p.pos.GetP()}
}
