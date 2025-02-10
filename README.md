[![Go Report Card](https://goreportcard.com/badge/github.com/rah-0/nabu)](https://goreportcard.com/report/github.com/rah-0/nabu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# nabu
`nabu` is a **structured logging library** for Go that provides **error tracking**, **log levels**, and **traceable logs** without requiring external dependencies.

With `nabu`, logs can:
- **Propagate errors automatically** while preserving their stack trace.
- **Attach structured metadata** to errors using `WithArgs()`.
- **Support multiple log outputs** (stdout, stderr).

## Installation

```sh
go get github.com/rah-0/nabu
```

## Usage

### Basic Logging

```go
func main() {
    nabu.SetLogLevel(nabu.LevelDebug) // Set log level
    nabu.SetLogOutput(nabu.OutputStdout) // Log to stdout

    nabu.FromMessage("Starting application").WithArgs("version", "1.0.0").WithLevelInfo().Log()
}
```
Log output:
```json
{"UUID":"f95a630a-94dd-444e-a749-40a2864ba32c","Date":"2025-02-10 14:52:07.624683","Args":["version","1.0.0"],"Msg":"Starting application","Function":"github.com/rah-0/nabu.TestSomething","Line":5,"Level":1}
```

### Logging Errors with Stack Traces
```go
func process() error {
    return errors.New("Something went wrong")
}
func main() {
    err := process()
    if err != nil {
        nabu.FromError(err).WithArgs("operation", "database").WithLevelError().Log()
    }
}
```
Log output:
```json
{"UUID":"14e37d51-2420-44f7-a7aa-403c13a048f0","Date":"2025-02-10 14:55:14.992355","Error":"Something went wrong","Args":["operation","database"],"Function":"github.com/rah-0/nabu.TestSomething","Line":7,"Level":3}
```

### Logging a Full Stack Trace chain
```go
// Function C (Deepest in the stack)
func functionC(userID int, action string) error {
    return errors.New("database connection failed")
}
// Function B (Calls functionC and wraps error)
func functionB(userID int, action string) error {
    err := functionC(userID, action)
    if err != nil {
        return nabu.FromError(err).WithArgs(userID, action).Log()
    }
    return nil
}
// Function A (Calls functionB and wraps error)
func functionA(userID int, action string) error {
    err := functionB(userID, action)
    if err != nil {
        return nabu.FromError(err).Log()
    }
    return nil
}
func main() {
    // Set up logging
    nabu.SetLogLevel(nabu.LevelDebug)
    nabu.SetLogOutput(nabu.OutputStdout)
    
    // Simulate an operation
    functionA(42, "delete_account")
}
```
Log output:
```json lines
{"UUID":"e5f04f91-f272-4562-b6d4-146bbd44d8dc","Date":"2025-02-10 14:57:27.163813","Error":"database connection failed","Args":[42,"delete_account"],"Function":"github.com/rah-0/nabu.TestSomething.functionB","Line":9}
{"UUID":"e5f04f91-f272-4562-b6d4-146bbd44d8dc","Date":"2025-02-10 14:57:27.163875","Error":"database connection failed","Function":"github.com/rah-0/nabu.TestSomething.functionA","Line":17}
```
