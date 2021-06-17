package JGScript

type BPos struct {
	Index int
	Row   int
	Col   int
	Fn    *string
	Code  *string
}

type Position struct {
	index int
	row   int
	col   int
	max   int
	fn    *string
	code  *string
}

func NewPosition(fn, code *string) *Position {
	return &Position{0, 1, 1, len(*code), fn, code}
}

func (p *Position) Next() rune {
	if p.index < p.max {
		if p.index < p.max && (*p.code)[p.index] == '\n' {
			p.row++
			p.col = 1
		} else {
			p.col++
		}
		p.index++
	}
	return p.Get()
}

func (p *Position) Get() rune {
	if p.index < p.max {
		return rune((*p.code)[p.index])
	}
	return 0
}

func (p *Position) GetN() rune {
	if p.index < p.max-1 {
		return rune((*p.code)[p.index+1])
	}
	return 0
}

func (p *Position) GetNN() rune {
	if p.index < p.max-2 {
		return rune((*p.code)[p.index+2])
	}
	return 0
}

func (p *Position) GetP() *BPos {
	return &BPos{p.index, p.row, p.col, p.fn, p.code}
}
