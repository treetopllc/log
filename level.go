package log

const (
	INFO  LogLevel = "info"
	ERROR LogLevel = "error"
	PANIC LogLevel = "panic"
)

type LogLevel string

func (ll *LogLevel) Set(s LogLevel) {
	switch *ll {
	case "":
		*ll = s
	case INFO:
		if s == ERROR || s == PANIC {
			*ll = s
		}
	case ERROR:
		if s == PANIC {
			*ll = s
		}
	case PANIC:
		//already at highest level
	}
}
