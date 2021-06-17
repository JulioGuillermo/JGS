package JGScript

import "fmt"

type Error struct {
	t string
	v string
	s *BPos
	e *BPos
	c *Context
}

func MakeError(t, v string, s, e *BPos, c *Context) *Error {
	return &Error{t, v, s, e, c}
}

func MakeRTError(v string, s, e *BPos, c *Context) *Error {
	return &Error{"RunTimeError", v, s, e, c}
}

func MakeSError(v string, s, e *BPos, c *Context) *Error {
	return &Error{"InvalidSyntaxError", v, s, e, c}
}

func trackback(c *Context) string {
	s := ""
	if c.parent != nil {
		s = trackback(c.parent)
	}
	if c.n != "" {
		s += "\nFile: " + *c.s.Fn + ", line: " + fmt.Sprint(c.s.Row) + ", column: " + fmt.Sprint(c.s.Col) + ", at \"" + c.n + "\""
	}
	return s
}

func ToStr(e *Error) string {
	s := trackback(e.c)
	s += "\n\nFile: " + *e.s.Fn + ", Line: " + fmt.Sprint(e.s.Row) + ", Column: " + fmt.Sprint(e.s.Col)
	s += "\n" + e.t + ": " + e.v + "\n"
	s += TextArrow(e.s, e.e)
	return s
}
