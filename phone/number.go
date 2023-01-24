package phone

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

const numberExpr = `(\d\d)-?(\d\d)-?(\d\d)`

var numberRe = regexp.MustCompile(numberExpr)

func NumberExpr() string {
	x := &strings.Builder{}
	for _, r := range numberExpr {
		if r != '(' && r != ')' {
			x.WriteRune(r)
		}
	}
	return x.String()
}

type Number int

func MakeNumber(s string) (Number, error) {
	b := make([]int, 0, 3)
	if !numberRe.MatchString(s) {
		return 0, errors.New("bad number format")
	}
	p := numberRe.FindStringSubmatch(s)
	for _, x := range p[1:] {
		n, err := strconv.Atoi(x)
		if err != nil {
			return 0, err
		}
		b = append(b, n)
	}
	return Number(b[0]*1e4 + b[1]*1e2 + b[2]), nil
}

func RandomNumber() Number {
	return Number(rand.Intn(1e6))
}

func (n Number) String() string {
	m := int(n)
	b := [...]int{
		m / 1e4,
		m / 1e2 % 1e2,
		m % 1e2,
	}
	return fmt.Sprintf("%02d-%02d-%02d", b[0], b[1], b[2])
}
