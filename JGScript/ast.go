package JGScript

type AST struct {
	T string
	V string
	C []*AST
	S *BPos
	E *BPos
}
