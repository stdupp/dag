package dag

import "sync"

type Context struct {
	taskMap sync.Map
	Item    interface{} // user data
}

// Dag represents directed acyclic graph
type Dag struct {
	Ctx  *Context
	jobs []*Job
}

// New creates new DAG
func New(ctx *Context) *Dag {
	return &Dag{
		Ctx:  ctx,
		jobs: make([]*Job, 0),
	}
}

func (dag *Dag) lastJob() *Job {
	jobsCount := len(dag.jobs)
	if jobsCount == 0 {
		return nil
	}

	return dag.jobs[jobsCount-1]
}

// Run starts the tasks
// It will block until all functions are done
func (dag *Dag) Run() {
	for _, job := range dag.jobs {
		if job.sequential {
			runSync(dag.Ctx, job)
		} else {
			runAsync(dag.Ctx, job)
		}
	}
}

// RunAsync executes Run on another goroutine
func (dag *Dag) RunAsync(onComplete func()) {
	go func() {
		dag.Run()
		if onComplete != nil {
			onComplete()
		}
	}()
}

// Pipeline executes tasks sequentially
func (dag *Dag) Pipeline(tasks ...DagTaskIFace) *pipelineResult {
	job := &Job{
		tasks:      make([]*Task, len(tasks)),
		sequential: true,
	}

	for i, task := range tasks {
		job.tasks[i] = taskWrap(task)
	}

	dag.jobs = append(dag.jobs, job)

	return &pipelineResult{dag}
}

// Spawns executes tasks concurrently
func (dag *Dag) Spawns(tasks ...DagTaskIFace) *spawnsResult {
	job := &Job{
		tasks:      make([]*Task, len(tasks)),
		sequential: false,
	}

	for i, task := range tasks {
		job.tasks[i] = taskWrap(task)
	}

	dag.jobs = append(dag.jobs, job)

	return &spawnsResult{dag}
}
