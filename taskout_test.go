package taskout

import (
	"context"
	"math"
	"testing"
	"time"
)

func TestInterval(t *testing.T) {
	timesPrinted := 0
	prinitLimit := 3

	taskManager := NewTaskManager()

	intervalId := taskManager.SetInterval(func(ctx context.Context) {
		timesPrinted += 1
	}, time.Second*2)

	time.Sleep(time.Second * 7)

	taskManager.Cancel(intervalId, nil)

	if timesPrinted != prinitLimit {
		t.Fatalf(`SetInterval() expected to execute %d times but executed %d times`, prinitLimit, timesPrinted)
	}
}

func TestIntervalWithExtend(t *testing.T) {
	expectedExecution := 3
	timesExecuted := 0

	taskManager := NewTaskManager()

	intervalId := taskManager.SetInterval(func(ctx context.Context) {
		timesExecuted += 1
	}, time.Duration(time.Second*3))

	taskManager.Extend(intervalId, time.Second*1)

	time.Sleep(time.Second * 3)

	taskManager.Cancel(intervalId, nil)

	time.Sleep(time.Second * 1)

	if expectedExecution != timesExecuted {
		t.Fatalf(`SetInterval() when used with Extend() expected to execute %d times but executed %d times`, expectedExecution, timesExecuted)
	}
}

func TestIntervalWithExecution(t *testing.T) {
	expectedExecution := 3
	timesExecuted := 0

	taskManager := NewTaskManager()

	intervalId := taskManager.SetInterval(func(ctx context.Context) {
		timesExecuted += 1
	}, time.Second*2)

	time.Sleep(time.Second * 5)

	taskManager.Execute(intervalId)

	time.Sleep(time.Second * 3)

	if expectedExecution != timesExecuted {
		t.Fatalf(`SetInterval() when used with Execute() expected to execute %d times but executed %d times`, expectedExecution, timesExecuted)
	}
}

func TestIntervalWithCancellation(t *testing.T) {
	expectedExecution := 0
	timesExecuted := 0

	taskManager := NewTaskManager()

	intervalId := taskManager.SetInterval(func(ctx context.Context) {
		timesExecuted += 1
	}, time.Second*2)

	time.Sleep(time.Second * 1)

	taskManager.Cancel(intervalId, nil)

	time.Sleep(time.Second * 3)

	if expectedExecution != timesExecuted {
		t.Fatalf(`SetInterval() when used with Cancel() expected to execute %d times but executed %d times`, expectedExecution, timesExecuted)
	}
}

func TestTimeout(t *testing.T) {
	timesExecuted := 0
	expectedExecution := 1

	taskManager := NewTaskManager()

	taskManager.SetTimeout(func(ctx context.Context) {
		timesExecuted += 1
	}, time.Second*2)

	time.Sleep(time.Second * 5)

	if timesExecuted != expectedExecution {
		t.Fatalf(`SetTimeout() expected to execute %d times but executed %d times`, expectedExecution, timesExecuted)
	}
}

func TestTimeoutWithExecution(t *testing.T) {
	startTime := time.Now()
	expectedExecutionTimeInSeconds := 1
	executionTimeInSecond := 0

	taskManager := NewTaskManager()

	timeoutId := taskManager.SetTimeout(func(ctx context.Context) {

		executionTimeInSecond = int(math.Floor(time.Since(startTime).Seconds()))
	}, time.Second*3)

	time.Sleep(time.Second * 1)

	taskManager.Execute(timeoutId)

	time.Sleep(time.Second * 1)

	if executionTimeInSecond != expectedExecutionTimeInSeconds {
		t.Fatalf(`SetTimeout() when used with Execute() expected to execute in %d second but executed %d in seconds`, expectedExecutionTimeInSeconds, executionTimeInSecond)
	}

}

func TestTimeoutWithExtend(t *testing.T) {
	startTime := time.Now()
	expectedExecutionTimeInSeconds := 3
	executionTimeInSecond := 0

	taskManager := NewTaskManager()

	timeoutId := taskManager.SetTimeout(func(ctx context.Context) {
		executionTimeInSecond = int(math.Floor(time.Since(startTime).Seconds()))
	}, time.Second*4)

	time.Sleep(time.Second * 1)

	taskManager.Extend(timeoutId, time.Second*2)

	time.Sleep(time.Second * 4)

	if executionTimeInSecond != expectedExecutionTimeInSeconds {
		t.Fatalf(`SetTimeout() when used with Extend() expected to execute in %d second but executed %d in seconds`, expectedExecutionTimeInSeconds, executionTimeInSecond)
	}
}

func TestTimeoutWithCancel(t *testing.T) {
	expectedExecution := 0
	timesExecuted := 0

	taskManager := NewTaskManager()

	timeoutId := taskManager.SetTimeout(func(ctx context.Context) {
		timesExecuted += 1
	}, time.Second*2)

	time.Sleep(time.Second * 1)

	taskManager.Cancel(timeoutId, nil)

	time.Sleep(time.Second * 3)

	if expectedExecution != timesExecuted {
		t.Fatalf(`SetTimeout() when used with Extend() expected to execute in %d second but executed %d in seconds`, expectedExecution, timesExecuted)
	}
}
