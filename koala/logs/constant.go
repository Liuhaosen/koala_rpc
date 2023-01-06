package logs

const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelAccess
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

const (
	DefaultLogChanSize = 20000
	SpaceSep           = " "
	ColonSep           = ":"
	LineSep            = "\n"
)

type LogLevel int

func GetLogLevel(level string) LogLevel {
	switch level {
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	case "access":
		return LogLevelAccess
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError

	}
	return LogLevelDebug
}

func getLevelColor(level LogLevel) Color {
	switch level {
	case LogLevelDebug:
		return White
	case LogLevelTrace:
		return Yellow
	case LogLevelInfo:
		return Green
	case LogLevelAccess:
		return Blue
	case LogLevelWarn:
		return Cyan
	case LogLevelError:
		return Red
	}
	return Magenta
}

func getLevelText(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "debug"
	case LogLevelTrace:
		return "trace"
	case LogLevelAccess:
		return "access"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	}
	return "debug"
}
