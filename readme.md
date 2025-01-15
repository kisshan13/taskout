# taskout

[![Go Reference](https://pkg.go.dev/badge/github.com/kisshan13/taskout.svg)](https://pkg.go.dev/github.com/kisshan13/taskout)

`taskout` is a lightweight task scheduling library for Go, inspired by JavaScript's `setTimeout` and `setInterval` functionalities. It enables developers to easily create one-shot or recurring tasks with flexible scheduling and lifecycle management.

## Features

- `SetTimeout` : Schedule a task to run once after a specified duration.
- `SetInterval` : Schedule a recurring task to execute at regular intervals.
- `Cancel` : Cancel scheduled tasks at any time.
- `Extend` : Extend the duration of a one-shot task or modify a recurring task interval.
- `Execute` :  Trigger a task to execute immediately.

## Installation

```bash
go get github.com/kisshan13/taskout
```

## Usage

### Basic Example : `SetTimeout`

Schedule a one-shot task to execute after a delay:

```go
package main

import (
	"fmt"
	"time"

	"github.com/kisshan13/taskout"
)

func main() {
	tm := taskout.NewTaskManager()

	taskID := tm.SetTimeout(func(ctx context.Context) {
		fmt.Println("One-shot task executed")
	}, 5*time.Second)

	time.Sleep(6 * time.Second) // Wait for the task to complete
	fmt.Println("Task Completed:", taskID)
}

```

### Basic Example : `SetInterval`

Schedule a recurring task to run at regular intervals:

```go
package main

import (
	"fmt"
	"time"

	"github.com/kisshan13/taskout"
)

func main() {
	tm := taskout.NewTaskManager()

	taskID := tm.SetInterval(func(ctx context.Context) {
		fmt.Println("Recurring task executed")
	}, 2*time.Second)

	time.Sleep(10 * time.Second) // Allow the task to execute multiple times
	tm.Cancel(taskID, func() {
		fmt.Println("Recurring task canceled")
	})
}
```

### Extending a task

Extend the duration of a task. For one-shot or timeout task it will overwrite the timeout period and for interval based
task it overwrite the interval period.

```go
package main

import (
	"fmt"
	"time"

	"github.com/kisshan13/taskout"
)

func main() {
	tm := taskout.NewTaskManager()

	taskID := tm.SetTimeout(func(ctx context.Context) {
		fmt.Println("One-shot task executed after extension")
	}, 5*time.Second)

	time.Sleep(3 * time.Second) // Wait before extending
	tm.Extend(taskID, 5*time.Second) // Extend the task duration

	time.Sleep(6 * time.Second) // Wait for the extended duration
}
```

### Triggering a Task Immediately

You can trigger a task to execute immediately:
(NOTE : Interval Tasks will removed from the tasks after a successful execution)

```go 
package main

import (
	"fmt"
	"time"

	"github.com/kisshan13/taskout"
)

func main() {
	tm := taskout.NewTaskManager()

	taskID := tm.SetTimeout(func(ctx context.Context) {
		fmt.Println("One-shot task executed early")
	}, 10*time.Second)

	tm.Execute(taskID) // Trigger the task immediately
}
```

## Contributing

Contributions are welcome! Follow these steps to contribute:

1. Fork the repository: github.com/kisshan13/taskout
2. Create your feature branch: git checkout -b feature/my-feature
3. Commit your changes: git commit -m "Add my feature"
4. Push to the branch: git push origin feature/my-feature
5. Open a pull request


## Reporting Issues
If you find any bugs or have feature requests, please open an issue on the [GitHub Issues](https://github.com/kisshan13/taskout/issues) page.

## License

This project is licensed under the MIT License.

## Acknowledgments

This project was inspired by JavaScript's `setTimeout` and `setInterval` functionalities, adapted to provide lightweight task scheduling in Go.