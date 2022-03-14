package main

import (
	"context"
	"fmt"

	"github.com/d5/tengo/v2"
)

func main() {
	script := tengo.NewScript([]byte(`
x := 1
f := field(10, r=56)

Teste := type("Teste",
	fields(
		loadCount=f
	),
	methods(
		do = func(this, ...args) {
			this.done = args
			x++
			return this
		},
		inc = func(this, ...args) {
			this.loadCount++
		}
	)
)

obj := Teste(x=5)

z := obj.inc
z()
z()
z()
f.tags.T = 34
T2 := type(obj.__type__.__map__)
out = string(T2.__map__)`))
	script.Add("out", tengo.UndefinedValue)

	compiled, err := script.RunContext(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(compiled.Get("out").Value())
}
