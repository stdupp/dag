package dag

import "sync"

type Context struct {
	taskMap sync.Map
	Data    interface{} // user data
}

// Dag represents directed acyclic graph
type Dag struct {
	jobs []*Job
}

// New creates new DAG
func New() *Dag {
	return &Dag{
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
func (dag *Dag) Run(ctx *Context) {
	for _, job := range dag.jobs {
		if job.sequential {
			runSync(ctx, job)
		} else {
			runAsync(ctx, job)
		}
	}
}

// RunAsync executes Run on another goroutine
func (dag *Dag) RunAsync(ctx *Context, onComplete func()) {
	go func() {
		dag.Run(ctx)
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
