package dag

// Job - Each job consists of one or more tasks
// Each Job can runs tasks in order(Sequential) or unordered
type Job struct {
	tasks      []*Task
	sequential bool
	onComplete func()
}

type DagTaskIFace interface {
	Name() string
	Process(*Context)
}

type Task struct {
	n string
	f func(*Context)
}

func (t *Task) Name() string {
	return t.n
}
func (t *Task) Process(ctx *Context) {
	t.f(ctx)
}

func taskWrap(dti DagTaskIFace) *Task {
	name := dti.Name()
	f := func(ctx *Context) {
		if name == "" {
			dti.Process(ctx)
			return
		}
		ch := make(chan struct{}, 0)
		v, loaded := ctx.taskMap.LoadOrStore(name, ch)
		if loaded {
			ch = v.(chan struct{})
			<-ch
		} else {
			dti.Process(ctx)
			close(ch)
		}
	}

	return &Task{name, f}
}

// Combine wraps tasks as a single task
func Combine(tasks ...DagTaskIFace) DagTaskIFace {
	f := func(ctx *Context) {
		for _, task := range tasks {
			taskWrap(task).f(ctx)
		}

	}
	return &Task{"", f}
}
