package log

type Logger interface {
	Log(...interface{})
	Logf(string, ...interface{})
	Error(error)
	Errorf(string, ...interface{})
	Panic(interface{})
	Panicf(string, ...interface{})
	Write([]byte) (int, error)
}

var debug bool = false

func SetDebug(v bool) {
	debug = v
}
