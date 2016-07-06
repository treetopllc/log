package log

import (
	"time"
)

type TimeDeltaLogger struct {
	l    Logger
	last time.Time
}

func NewTimeDeltaLogger(l Logger) *TimeDeltaLogger {
	return &TimeDeltaLogger{
		l:    l,
		last: time.Now(),
	}
}

// Report logs a message, along with the number of milliseconds elapsed since the last call to Report. Alternatively, if
// this is the first call to this instance's Report(), it will display the number of milliseconds since initialization
// of the TimeDeltaLogger.
func (tdl *TimeDeltaLogger) Report(msg string) {
	now := time.Now()
	delta := now.Sub(tdl.last)
	ms := delta.Nanoseconds() / 1000 / 1000
	tdl.l.Logf("%dms elapsed: %s", ms, msg)
	tdl.last = now
}
