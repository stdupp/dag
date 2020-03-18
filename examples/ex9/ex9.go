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

// A workflow like this
//                    +-----+
//                    |     |
//             +----->+  2  +------+
//             |      |     |      |
// +-----+     |      +-----+      |      +-----+
// |     |     |                   |      |     |
// |  1  +-----+                   +----->+  4  +-----------------+
// |     |     |                   |      |     |                 |
// +-----+     |                   |      +-----+                 v
//             |      +-----+      |                           +--+--+
//             |      |     |      |                           |     |
//             +----->+  3  +------+                           |  7  |
//                    |     |                                  |     |
//                    +--+--+                                  +--+--+
//                       ^                                        ^
// +-----+               |                +-----+                 |
// |     |               |                |     |                 |
// |  5  +---------------+--------------->+  6  +-----------------+
// |     |                                |     |
// +-----+                                +-----+
//
var t1, t2, t3, t4, t5, t6, t7, ta, tb *Task

func main() {
	t1 = &Task{"1", f1}
	t2 = &Task{"2", f2}
	t3 = &Task{"3", f3}
	t4 = &Task{"4", f4}
	t5 = &Task{"5", f5}
	t6 = &Task{"6", f6}
	t7 = &Task{"7", f7}

	ta = &Task{"", fa}
	tb = &Task{"", fb}

	ctx := &dag.Context{Item: make([]int, 7)}
	d := dag.New()
	d.Spawns(ta, dag.Combine(t5, t6)).
		Join().
		Pipeline(t7)
	d.Run(ctx)

	fmt.Println(ctx.Item)
}

func f1(ctx *dag.Context) {
	v := ctx.Item.([]int)
	v[0] = 1
	println("f1")
}
func f2(ctx *dag.Context) {
	v := ctx.Item.([]int)
	v[1] = 1
	println("f2")
}
func f3(ctx *dag.Context) {
	v := ctx.Item.([]int)
	v[2] = 1
	println("f3")
}
func f4(ctx *dag.Context) {
	v := ctx.Item.([]int)
	v[3] = 1
	println("f4")
}
func f5(ctx *dag.Context) {
	v := ctx.Item.([]int)
	v[4] = 1
	println("f5")
}
func f6(ctx *dag.Context) {
	v := ctx.Item.([]int)
	v[5] = 1
	println("f6")
}
func f7(ctx *dag.Context) {
	v := ctx.Item.([]int)
	v[6] = 1
	println("f7")
}

func fa(ctx *dag.Context) {
	d := dag.New()
	d.Pipeline(t1).Then().Spawns(t2, tb).Join().Pipeline(t4)
	d.Run(ctx)
}

func fb(ctx *dag.Context) {
	d := dag.New()
	d.Spawns(t5, t1).Join().Pipeline(t3)
	d.Run(ctx)
}
