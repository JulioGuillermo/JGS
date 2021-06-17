package JGScript

// Type
type Type interface {
	GetT() string
	str() string
	getBool() bool
	getByte() byte
	getInt() int64
	getFloat() float64
	equals(t *Type) bool
	operator(op string, t *Type, c *Context, s, e *BPos) *Result
    GetMember(a *AST, c *Context) *Result
}
