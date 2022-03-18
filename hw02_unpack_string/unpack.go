package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type parseState int

const (
	processing parseState = iota
	prevIsSym
	prevIsSlash
)

func Unpack(str string) (string, error) {
	var state parseState
	var prevSym, curSym string
	var res strings.Builder

	for _, s := range str {
		curSym = string(s)
		num, err := strconv.Atoi(curSym)

		switch state {
		case processing:
			if err == nil {
				return "", ErrInvalidString
			}

			switch curSym {
			case "\\":
				state = prevIsSlash
			default:
				prevSym = curSym
				state = prevIsSym
			}
		case prevIsSym:
			if err == nil {
				res.WriteString(strings.Repeat(prevSym, num))
				state = processing
				continue
			}

			res.WriteString(prevSym)

			switch curSym {
			case "\\":
				state = prevIsSlash
			default:
				prevSym = curSym
				state = prevIsSym
			}
		case prevIsSlash:
			prevSym = curSym
			state = prevIsSym
		}
	}

	if state == prevIsSym {
		res.WriteString(prevSym)
	}

	return res.String(), nil
}
