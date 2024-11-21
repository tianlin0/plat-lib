package task

const (
	JobStatusTodo   JobStatus = "todo"
	JobStatusDone   JobStatus = "done"
	JobStatusFailed JobStatus = "failed"
	jobStatusZombie JobStatus = "zombie"
)

type JobStatus string

type Task interface {
	// Push param payload, return id, error
	Push([]byte) (int64, error)
	// Commit commit a job to a status under manual commit mode, and will not work under auto commit mode
	Commit(int64, int64, JobStatus) error
	// RePush a job by id to redo task
	RePush(int64) error
}
