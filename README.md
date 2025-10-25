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
    nabu.FromMessage("Starting application").WithArgs("version", "1.0.0").Log()
}
```
```json
{"UUID":"abc123...","Date":"2025-02-10 21:12:36.388657","Args":["version","1.0.0"],"Msg":"Starting application","Level":1}
```

### Logging Errors

```go
func main() {
    err := process()
    if err != nil {
        nabu.FromError(err).WithArgs("operation", "database").Log()
    }
}
```
```json
{"UUID":"be985ee7...","Date":"2025-02-10 21:13:57.887870","Error":"Something went wrong","Args":["operation","database"],"Function":"main.main","Line":7,"Level":3}
```

### Error Chains and UUID Correlation

**Wrapping errors preserves the UUID** - all logs in a chain share the same UUID for correlation:

```go
func functionC() error {
    return errors.New("database connection failed")
}

func functionB(userID int) error {
    err := functionC()
    if err != nil {
        return nabu.FromError(err).WithArgs("userID", userID).WithMessage("query failed").Log()
    }
    return nil
}

func functionA(userID int) error {
    err := functionB(userID)
    if err != nil {
        return nabu.FromError(err).WithMessage("operation failed").Log()
    }
    return nil
}
```
```json
{"UUID":"0a1feb11...","Date":"...","Error":"database connection failed","Args":["userID",42],"Msg":"query failed","Function":"main.functionB","Line":9,"Level":3}
{"UUID":"0a1feb11...","Date":"...","Msg":"operation failed","Function":"main.functionA","Line":17,"Level":3}
```

**Key behaviors:**
- Logs in the same chain share the **same UUID**
- Use `WithMessage()` to add context at each level
- Works for both error chains and message chains

### Custom UUIDs for Cross-Service Correlation

Use `WithUuid()` to set a custom UUID for correlating logs across services:

```go
func handleRequest(traceID string) {
    err := validateInput()
    if err != nil {
        nabu.FromError(err).WithUuid(traceID).WithMessage("validation failed").Log()
    }
}
```
```json
{"UUID":"frontend-trace-12345","Date":"...","Error":"invalid email","Msg":"validation failed","Function":"main.handleRequest","Line":15,"Level":3}
```

## API Reference

**Creating Loggers:**
- `FromError(err error) *Logger` - Create from error (auto-generates UUID)
- `FromMessage(msg string) *Logger` - Create from message (auto-generates UUID)
- `New() *Logger` - Create empty logger

**Configuring Loggers:**
- `WithMessage(msg string)` - Add/update message
- `WithArgs(args ...any)` - Attach structured data
- `WithUuid(uuid string)` - Set custom UUID
- `WithLevel{Debug|Info|Warn|Error|Fatal}()` - Set log level
- `Log()` - Output the log

**Global Settings:**
- `SetLogLevel(level Level)` - Set minimum log level
- `SetLogOutput(output Output)` - Set output (stdout/stderr)

**Log Levels:** `LevelDebug` (1), `LevelInfo` (2), `LevelWarn` (3), `LevelError` (4), `LevelFatal` (5)

## Features

✅ Structured JSON logging  
✅ Automatic UUID generation for log correlation  
✅ Error chain tracking with preserved stack traces  
✅ Custom UUIDs for cross-service correlation  

# ☕ Support
Enjoying nabu?
If it saved you time or brought value to your project, feel free to show some support. Every bit is appreciated 🙂

[![Buy Me A Coffee](https://cdn.buymeacoffee.com/buttons/default-orange.png)](https://www.buymeacoffee.com/rah.0)
