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
	// Bits or'ed together to control what's printed. There is no control over the
	// order they appear (the order listed here) or the format they present (as
	// described in the comments).  A colon appears after these items:
	//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota     // the date: 2009/01/23
	Ltime                         // the time: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
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
