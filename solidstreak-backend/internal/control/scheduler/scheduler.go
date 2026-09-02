package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"

	st "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/schedulertask"
)

type Scheduler struct {
	MaxTaskHandlers        int
	TaskWaitingDuration    time.Duration
	TaskBatchSizePerSource int
	LockDuration           time.Duration
	Res                    common.Resources
	TaskSources            []st.TaskSource
}

func (s Scheduler) Run(ctx context.Context, doneCh chan struct{}) {
	defer func() { doneCh <- struct{}{} }()

	s.Res.Logger.Info("scheduler started")

	handlers := make(map[string]struct{}, s.MaxTaskHandlers)
	handlerDoneCh := make(chan string, s.MaxTaskHandlers)
	handlerID := ""

	noTaskTicker := time.NewTicker(s.TaskWaitingDuration)
	haveTaskTicker := time.NewTicker(time.Millisecond)
	ticker := haveTaskTicker

	var tasks []st.Task
	var err error

loop:
	for {
		select {
		// scheduler stopping
		case <-ctx.Done():
			break loop

		// handler finished
		case handlerID = <-handlerDoneCh:
			delete(handlers, handlerID)

		// scheduler tick
		case <-ticker.C:
			if len(handlers) >= s.MaxTaskHandlers {
				s.Res.Logger.Warn("scheduler has no task handler available, waiting for a handler to finish")
				handlerID = <-handlerDoneCh
				delete(handlers, handlerID)
			}
			handlerID = uuid.NewString()[:8]
			if tasks, err = s.getTasksFromSources(handlerID); err != nil {
				s.Res.Logger.Error("scheduler failed to fetch tasks from sources", "error", err)
			}
			if len(tasks) > 0 {
				handlers[handlerID] = struct{}{}
				go TaskRunner{
					ID:  handlerID,
					Res: s.Res,
				}.Run(ctx, handlerDoneCh, tasks)
				ticker = haveTaskTicker
			} else {
				s.Res.Logger.Debug("scheduler found no tasks to process, waiting for next tick")
				ticker = noTaskTicker
			}
		}
	}

	// wait for all handlers to finish
	for len(handlers) > 0 {
		s.Res.Logger.Info("scheduler waiting for task handlers to finish")
		handlerID := <-handlerDoneCh
		delete(handlers, handlerID)
	}

	s.Res.Logger.Info("scheduler execution completed")
}

func (s Scheduler) getTasksFromSources(handlerID string) ([]st.Task, error) {
	var (
		resultTasks []st.Task
		err         error
		source      st.TaskSource
		tasks       []st.Task
	)

	for _, source = range s.TaskSources {
		if tasks, err = source.FetchSchedulerTasksWithLocking(s.TaskBatchSizePerSource, s.LockDuration, handlerID); err != nil {
			return nil, err
		}
		resultTasks = append(resultTasks, tasks...)
	}

	return resultTasks, nil
}
