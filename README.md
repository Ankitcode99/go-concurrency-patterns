# Go Concurrency Patterns

This project demonstrates various Go concurrency patterns, including:

1. **Fan-in Pattern**: This pattern is used to merge multiple input channels into a single output channel. It is useful for aggregating data from multiple sources into a single stream. The `fanIn` function in `fan-in-pattern/main.go` demonstrates this pattern.
2. **Fan-out Pattern**: This pattern is used to split a single input channel into multiple output channels. It is useful for distributing work items into multiple uniform actors. The `processLogs` function in `fan-out-pattern/main.go` demonstrates this pattern.
3. **Function Returning a Channel**: This pattern is used to encapsulate the creation of a goroutine and a channel within a function. It simplifies the process of creating a channel and starting a goroutine that sends data on that channel. The `boring` function in `function-returning-channel/main.go` demonstrates this pattern.
4. **Worker Pool Pattern**: This pattern is used to manage a pool of workers that can be dynamically scaled up or down based on the workload. It is useful for efficiently handling a large number of tasks without creating a new goroutine for each task. The `workerPool` function in `worker-pool-pattern/main.go` demonstrates this pattern.


These patterns are essential in Go programming for handling concurrency effectively. They enable developers to write efficient, scalable, and concurrent code that can take full advantage of modern CPU architectures.

### Running the Examples

To run the examples, navigate to the respective directory and execute the `go run main.go` command. This will start the program and display the output in the console.

### Understanding the Code

Each example is accompanied by a detailed explanation in the form of comments within the code. These comments provide insight into the purpose and implementation of each pattern, making it easier to understand and learn from the examples.

### Contributing

This project is open to contributions. If you have a Go concurrency pattern you'd like to share, feel free to submit a pull request with your example.
