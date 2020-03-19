package main

import (
	"fmt"

	"github.com/stdupp/dag"
)

type Task struct {
	N string
	F func(*dag.Context)
}

func (t *Task) Name() string {
	return t.N
}
func (t *Task) Process(ctx *dag.Context) {
	t.F(ctx)
}

func f(name string, i int) func(*dag.Context) {
	return func(ctx *dag.Context) {
		v := ctx.Data.([]int)
		v[i] = 1
		println(name)
	}
}

func main() {
	t1 := &Task{"1", f("f1", 0)}
	t2 := &Task{"2", f("f2", 1)}
	t3 := &Task{"3", f("f3", 2)}
	t4 := &Task{"4", f("f4", 3)}
	t5 := &Task{"5", f("f5", 4)}
	t6 := &Task{"6", f("f6", 5)}
	t7 := &Task{"7", f("f7", 6)}

	fb := func(ctx *dag.Context) {
		d := dag.New()
		d.Spawns(t5, t1).Join().Pipeline(t3)
		d.Run(ctx)
	}
	tb := &Task{"", fb}

	fa := func(ctx *dag.Context) {
		d := dag.New()
		d.Pipeline(t1).Then().Spawns(t2, tb).Join().Pipeline(t4)
		d.Run(ctx)
	}
	ta := &Task{"", fa}

	ctx := &dag.Context{Data: make([]int, 7)}
	d := dag.New()
	d.Spawns(ta, dag.Combine(t5, t6)).
		Join().
		Pipeline(t7)
	d.Run(ctx)

	fmt.Println(ctx.Data)
}
