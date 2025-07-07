[![Go Report Card](https://goreportcard.com/badge/github.com/rah-0/nabu)](https://goreportcard.com/report/github.com/rah-0/nabu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<a href="https://www.buymeacoffee.com/rah.0" target="_blank">
  <img src="https://cdn.buymeacoffee.com/buttons/v2/arial-orange.png" alt="Buy Me A Coffee" height="50" style="height:50px;">
</a>

# nabu
`nabu` is a **structured logging library** for Go that provides **error tracking**, **log levels**, and **traceable logs** without requiring external dependencies.

With `nabu`, logs can:
- **Propagate errors** while preserving their stack trace.
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
    nabu.SetLogLevel(nabu.LevelDebug)
	
    nabu.FromMessage("Starting application").WithArgs("version", "1.0.0").WithLevelInfo().Log()
}
```
Log output:
```json
{"Date":"2025-02-10 21:12:36.388657","Args":["version","1.0.0"],"Msg":"Starting application","Level":1}
```

### Logging Errors with Stack Traces
```go
func process() error {
    return errors.New("Something went wrong")
}
func main() {
    err := process()
    if err != nil {
        nabu.FromError(err).WithArgs("operation", "database").Log()
    }
}
```
Log output:
```json
{"UUID":"be985ee7-d3e3-42c8-a9db-4422f1f32e96","Date":"2025-02-10 21:13:57.887870","Error":"Something went wrong","Args":["operation","database"],"Function":"github.com/rah-0/nabu.TestSomething","Line":7,"Level":3}
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
{"UUID":"0a1feb11-b250-4790-bad9-3a187df6f0f6","Date":"2025-02-10 21:15:24.790412","Error":"database connection failed","Args":[42,"delete_account"],"Function":"github.com/rah-0/nabu.functionB","Line":9,"Level":3}
{"UUID":"0a1feb11-b250-4790-bad9-3a187df6f0f6","Date":"2025-02-10 21:15:24.790458","Error":"database connection failed","Function":"github.com/rah-0/nabu.functionA","Line":17,"Level":3}
```

# ☕ Support
Enjoying nabu?
If it saved you time or brought value to your project, feel free to show some support. Every bit is appreciated 🙂

[![Buy Me A Coffee](https://cdn.buymeacoffee.com/buttons/default-orange.png)](https://www.buymeacoffee.com/rah.0)
