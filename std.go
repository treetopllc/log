package log

import "os"

var std = NewMulti(NewBasic(nil, ""))

func Std() Logger {
	return std
}
func Add(l ...Logger) {
	std.Add(l...)
}

func Print(msg ...interface{}) {
	std.Log(msg...)
}
func Println(msg ...interface{}) {
	std.Log(msg...)
}
func Printf(msg string, args ...interface{}) {
	std.Logf(msg, args...)
}
func Log(msg ...interface{}) {
	std.Log(msg...)
}
func Logf(msg string, args ...interface{}) {
	std.Logf(msg, args...)
}
func Error(err error) {
	std.Error(err)
}
func Errorf(msg string, args ...interface{}) {
	std.Errorf(msg, args...)
}
func Panic(msg interface{}) {
	std.Panic(msg)
}
func Panicf(msg string, args ...interface{}) {
	std.Panicf(msg, args...)
}
func Fatal(msg ...interface{}) {
	std.Log(msg...)
	os.Exit(1)
}
func Fatalf(msg string, args ...interface{}) {
	std.Logf(msg, args...)
	os.Exit(1)
}
