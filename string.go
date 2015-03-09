package log

import "fmt"

type StringLogger string

func (s *StringLogger) addline(str string) {
	*s += StringLogger(str + "\n")
}

func (s *StringLogger) Write(b []byte) (int, error) {
	s.addline(string(b))
	return len(b), nil
}
func (s *StringLogger) Log(msg ...interface{}) {
	s.addline(fmt.Sprint(msg...))
}
func (s *StringLogger) Logf(msg string, args ...interface{}) {
	s.addline(fmt.Sprintf(msg, args...))
}
func (s *StringLogger) Error(err error) {
	s.addline(err.Error())
}
func (s *StringLogger) Errorf(msg string, args ...interface{}) {
	s.addline(fmt.Sprintf(msg, args...))
}
func (s *StringLogger) Panic(msg interface{}) {
	s.addline(fmt.Sprint(msg))
	panic(msg)
}
func (s *StringLogger) Panicf(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	s.addline(msg)
	panic(msg)
}
