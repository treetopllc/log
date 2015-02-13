package log

import (
	"io"
	"log"
	"os"
)

type Basic struct {
	inner *log.Logger
}

//from https://golang.org/src/log/log.go#L25
const (
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LstdFlags     = log.LstdFlags
)

func NewBasic(out io.Writer, prefix string, flags ...int) Basic {
	if out == nil {
		out = os.Stderr
	}
	flag := LstdFlags
	if len(flags) != 0 {
		flag = 0
		for _, f := range flags {
			flag = flag | f
		}
	}
	return Basic{log.New(out, prefix, flag)}
}

func (b Basic) Log(msg ...interface{}) {
	b.inner.Print(msg...)
}
func (b Basic) Logf(msg string, args ...interface{}) {
	b.inner.Printf(msg, args...)
}
func (b Basic) Error(err error) {
	b.inner.Print(err)
}
func (b Basic) Errorf(msg string, args ...interface{}) {
	b.inner.Printf(msg, args...)
}
func (b Basic) Panic(msg interface{}) {
	b.inner.Panic(msg)
}
func (b Basic) Panicf(msg string, args ...interface{}) {
	b.inner.Panicf(msg, args...)
}
