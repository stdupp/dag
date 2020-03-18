# DAG Executor
This repository is inspired by [mostafa-asg/dag](https://github.com/mostafa-asg/dag).

DAG executor has three main concept:
1. pipeline executes the functions sequentially and in order.
2. spawns executes the functions concurrently, so there is no ordering guarantee.
3. same name task in one context only exec once.

## Example
```Go
//
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
//

// A workflow like this can be abstracted with the following code:

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
		v := ctx.Item.([]int)
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

	ctx := &dag.Context{Item: make([]int, 7)}
	d := dag.New()
	d.Spawns(ta, dag.Combine(t5, t6)).
		Join().
		Pipeline(t7)
	d.Run(ctx)

	fmt.Println(ctx.Item)
}

```

