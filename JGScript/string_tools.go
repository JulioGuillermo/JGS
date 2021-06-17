package JGScript

func CIsDigit(c rune) bool {
	digits := "0123456789"
	for _, d := range digits {
		if c == d {
			return true
		}
	}
	return false
}

func CIsODigit(c rune) bool {
	digits := "01234567"
	for _, d := range digits {
		if c == d {
			return true
		}
	}
	return false
}

func CIsHDigit(c rune) bool {
	digits := "0123456789abcdefABCDEF"
	for _, d := range digits {
		if c == d {
			return true
		}
	}
	return false
}

func CIsAlpha(c rune) bool {
	alpha := "ABCDEFGHIJKLMNÑOPQRSTUVWXYZabcdefghijklmnñopqrstuvwxyz_áéíóúÁÉÍÓÚ"
	for _, a := range alpha {
		if c == a {
			return true
		}
	}
	return false
}

func CIsAlphaDigit(c rune) bool {
	return CIsDigit(c) || CIsAlpha(c)
}

func IntToHexStr(n int64) string {
    hexmap := "0123456789ABCDEF"
    s := ""
    for n > 0 {
        i := n % 16
        s = string(hexmap[i]) + s
        n /= 16
    }
    if s == "" {
        s = "00"
    } else {
        for len(s) % 2 != 0 {
            s = "0" + s
        }
    }
    return "#" + s
}

func IntToOctStr(n int64) string {
    octmap := "01234567"
    s := ""
    for n > 0 {
        i := n % 8
        s = string(octmap[i]) + s
        n /= 8
    }
    if s == "" {
        s = "00"
    } else {
        for len(s) % 2 != 0 {
            s = "0" + s
        }
    }
    return "@" + s
}

func IntToBinStr(n int64) string {
    binmap := "01"
    s := ""
    for n > 0 {
        i := n % 2
        s = string(binmap[i]) + s
        n /= 2
    }
    if s == "" {
        s = "00000000"
    } else {
        for len(s) % 8 != 0 {
            s = "0" + s
        }
    }
    return "$" + s
}

func TextArrow(s, e *BPos) string {
	str := ""
	text := *(s.Code)
	max := len(text)
	ls := s.Index
	le := ls
	for i := ls - 1; i >= 0; i-- {
		if text[i] == '\n' {
			ls = i
			break
		}
	}
	for i := le; i < max; i++ {
		if text[i] == '\n' {
			le = i
			break
		}
	}
	lc := e.Row - s.Row + 1
	var sl string
	var cs uint64
	var ce uint64
	for i := 0; i < lc; i++ {
		sl = text[ls:le]
		if i == 0 {
			cs = uint64(s.Col)
		} else {
			cs = 0
		}
		if i == lc-1 {
			ce = uint64(e.Col)
		} else {
			ce = uint64(len(sl)) - 1
		}
		str += sl + "\n"
		for i, c := range sl {
			if uint64(i+1) >= cs && uint64(i+1) < ce {
				str += "^"
			} else {
				if c == '\t' {
					str += "\t"
				} else {
					str += " "
				}
			}
		}
		ls = le
		le = max
		for i := ls + 1; i < le; i++ {
			if text[i] == '\n' {
				le = i
				break
			}
		}
	}
	return str
}
