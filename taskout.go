package taskout

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TaskManager is responsible for managing tasks.
// It supports scheduling, canceling, and extending tasks.
type TaskManager struct {
	mu    sync.Mutex
	tasks map[TaskID]*Task
}

// Task represents a scheduled task, either one-shot or recurring.
type Task struct {
	ctx       context.Context
	cancel    context.CancelFunc
	extension chan time.Duration
	execute   chan bool
	interval  time.Duration
	oneShot   bool
	ticker    *time.Ticker
}

// TaskID is a unique identifier for a task.
type TaskID string

// NewTaskManager creates and initializes a TaskManager.
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[TaskID]*Task),
	}
}

// deleteTask removes a task from the manager.
func (tm *TaskManager) deleteTask(id TaskID) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	delete(tm.tasks, id)
}

// addTask adds a task to the manager.
func (tm *TaskManager) addTask(id TaskID, task *Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.tasks[id] = task
}

// run manages the lifecycle of a task, executing it based on its schedule.
func (tm *TaskManager) run(id TaskID, duration time.Duration, task func(ctx context.Context)) {
	tm.mu.Lock()
	t := tm.tasks[id]
	tm.mu.Unlock()

	if t.oneShot {
		for {
			select {
			case <-t.ctx.Done():
				if t.ctx.Err() == context.DeadlineExceeded {
					task(t.ctx)
					tm.deleteTask(id)
					return
				}

				if t.ctx.Err() == context.Canceled {
					tm.deleteTask(id)
					return
				}

			case <-t.execute:
				task(t.ctx)
				tm.deleteTask(id)
				return

			case ext := <-t.extension:
				ctx, cancel := context.WithTimeout(context.Background(), ext)
				t.ctx = ctx
				t.cancel = cancel
			}
		}
	} else {
		for {
			select {
			case <-t.ctx.Done():
				// For SetTimeout based tasks or one-time tasks
				if t.ctx.Err() == context.Canceled {
					tm.deleteTask(id)
					return
				}

			case <-t.execute:
				task(t.ctx)
				tm.deleteTask(id)
				return

			case ext := <-t.extension:
				t.ticker.Reset(ext)

			case <-t.ticker.C:
				task(t.ctx)
			}
		}
	}
}

// SetTimeout schedules a one-shot task to be executed after a specified duration.
func (tm *TaskManager) SetTimeout(task func(ctx context.Context), duration time.Duration) TaskID {
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	t := &Task{
		ctx:       ctx,
		cancel:    cancel,
		extension: make(chan time.Duration, 1),
		execute:   make(chan bool, 1),
		oneShot:   true,
	}

	taskId, _ := generateId()

	tm.addTask(TaskID(taskId), t)

	go tm.run(TaskID(taskId), duration, task)

	return TaskID(taskId)
}

// SetInterval schedules a recurring task with a specified interval.
func (tm *TaskManager) SetInterval(task func(context.Context), interval time.Duration) TaskID {
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(interval)

	t := &Task{
		ctx:       ctx,
		cancel:    cancel,
		extension: make(chan time.Duration, 1),
		execute:   make(chan bool),
		interval:  interval,
		ticker:    ticker,
		oneShot:   false,
	}

	taskId, _ := generateId()

	tm.addTask(TaskID(taskId), t)

	go tm.run(TaskID(taskId), interval, task)

	return TaskID(taskId)
}

// Cancel cancels a task and optionally executes a cleanup function.
func (tm *TaskManager) Cancel(id TaskID, onDelete func()) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if t, exists := tm.tasks[id]; exists {
		t.cancel()
		if onDelete != nil {
			onDelete()
		}

		delete(tm.tasks, id)
	}
}

// Extend adds additional duration to a task.
func (tm *TaskManager) Extend(id TaskID, duration time.Duration) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if task, exists := tm.tasks[id]; exists {
		select {
		case task.extension <- duration:
			return nil
		default:
			fmt.Println("Failed to extend duration, channel busy: ", id)
		}
	} else {
		return fmt.Errorf("invalid task id")
	}

	return nil
}

// Execute triggers a task immediately and remove from tasks.
func (tm *TaskManager) Execute(id TaskID) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if task, exists := tm.tasks[id]; exists {
		select {
		case task.execute <- true:
			return nil
		default:
			return nil
		}
	}

	return fmt.Errorf("invalid task id")
}
