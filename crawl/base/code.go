package base

import (
	"bytes"
	"fmt"
	"strconv"
)

type Code int

func NewCode(code int) *Code {
	c := Code(code)
	return &c
}

func FromSymbol(symbol string) (*Code, bool) {
	return FromSymbolByte([]byte(symbol))
}

func FromSymbolByte(symbol []byte) (*Code, bool) {
	i := 0
	l := len(symbol)
	for ; i < l; i++ {
		if symbol[i] < '0' || symbol[i] > '9' {
			continue
		}
		break
	}
	if i == l {
		return nil, false
	}

	where := bytes.ToLower(symbol[:i])
	code, _ := strconv.Atoi(string(symbol[i:]))
	if code < 1 {
		return nil, false
	}

	c := NewCode(code)
	if bytes.Equal(where, []byte("sh")) {
		return c, c.InSh()
	}

	if bytes.Equal(where, []byte("sz")) {
		return c, c.InSz()
	}
	return nil, false
}

func (c Code) String() string {
	if c < 1 {
		return ""
	}
	if c.InSh() {
		return fmt.Sprintf("sh%06d", c)
	}
	return fmt.Sprintf("sz%06d", c)
}

func (c Code) InSecondBoardMarket() bool {
	return c/1000 == 300
}

func (c Code) InSh() bool {
	return c > 1000
}

func (c Code) InSz() bool {
	return c > 0
}

func (c Code) InShA() bool {
	return c/10000 == 60
}

func (c Code) InShB() bool {
	return c/10000 == 90
}

func (c Code) InSzA() bool {
	return c != 0 && c/1000 == 0
}

func (c Code) InSzB() bool {
	return c/1000 == 200
}

func (c Code) InSmeBoardMarket() bool {
	return c/1000 == 2
}

func IsChinaShareCode(symbol string) bool {
	_, ok := FromSymbol(symbol)
	return ok
}
