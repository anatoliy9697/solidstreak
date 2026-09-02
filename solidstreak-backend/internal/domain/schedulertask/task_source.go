package schedulertask

import "time"

type TaskSource interface {
	FetchSchedulerTasksWithLocking(int, time.Duration, string) ([]Task, error)
}
