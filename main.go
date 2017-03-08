package main

import (
	"fmt"
	"strings"
)

func main() {
	r := strings.NewReader(`(a b 123 1.8e3 "a b")`)
	l := NewLexer(r)

	fmt.Println(l.Tokens())
}
