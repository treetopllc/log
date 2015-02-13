package log

type Logger interface {
	Log(...interface{})
	Logf(string, ...interface{})
	Error(error)
	Errorf(string, ...interface{})
	Panic(interface{})
	Panicf(string, ...interface{})
}
