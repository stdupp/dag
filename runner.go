package dag

import "sync"

func runAsync(ctx *Context, job *Job) {
	wg := &sync.WaitGroup{}

	wg.Add(len(job.tasks))
	for _, task := range job.tasks {
		go func(t *Task) {
			t.f(ctx)
			wg.Done()
		}(task)
	}
	wg.Wait()

	if job.onComplete != nil {
		job.onComplete()
	}
}

func runSync(ctx *Context, job *Job) {
	for _, task := range job.tasks {
		task.f(ctx)
	}
	if job.onComplete != nil {
		job.onComplete()
	}
}
