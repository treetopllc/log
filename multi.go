package log

import "fmt"

type Multi []Logger

func NewMulti(ll ...Logger) *Multi {
	return (*Multi)(&ll)
}

func (m *Multi) Add(ll ...Logger) {
	for _, l := range ll {
		if l != nil {
			*m = append(*m, l)
		}
	}
}

func (m Multi) Log(msg ...interface{}) {
	for _, l := range m {
		l.Log(msg...)
	}
}
func (m Multi) Logf(msg string, args ...interface{}) {
	for _, l := range m {
		l.Logf(msg, args...)
	}
}
func (m Multi) Error(err error) {
	for _, l := range m {
		l.Error(err)
	}
}
func (m Multi) Errorf(msg string, args ...interface{}) {
	for _, l := range m {
		l.Errorf(msg, args...)
	}
}
func (m Multi) Panic(msg interface{}) {
	for _, l := range m {
		func() {
			defer recover() //Anywhere else this would be evil
			l.Panic(msg)
		}()
	}
	panic(msg)
}
func (m Multi) Panicf(msg string, args ...interface{}) {
	for _, l := range m {
		func() {
			defer recover() //Anywhere else this would be evil
			l.Panic(msg)
		}()
	}
	panic(fmt.Sprintf(msg, args...))
}
