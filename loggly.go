package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	loggly "github.com/treetopllc/go-loggly"
)

type LogglyClient struct {
	Basic
	debug  bool
	client loggly.Client
}

func NewLogglyClient(token string, debug bool, tags ...string) *LogglyClient {
	c := &LogglyClient{
		client: loggly.NewClient(token, tags...),
		debug:  debug,
	}
	c.Basic = NewBasic(c, "", 0)
	return c
}

func (lc *LogglyClient) Send(e *LogglyEntry) {
	if e != nil {
		go func() {
			if !lc.debug {
				e.Response.Body = nil
				e.Request.Body = filterRequest(e.Request.Body)
			}
			b, err := json.Marshal(e) //no fields here can json err
			if err != nil {
				fmt.Println("Error sending to loggly: ", err)
			} else {
				lc.Write(b)
			}
		}()
	}
}

func (lc *LogglyClient) Write(msg []byte) (int, error) {
	go lc.client.Send(msg)
	return len(msg), nil
}

type LogglyEntry struct {
	Logs  StringLogger
	Level LogLevel

	Request struct {
		Proto  string
		Method string
		Host   string
		Path   string
		Query  string
		Body   string
		Header http.Header
	}
	Response struct {
		Body     interface{}
		Status   int
		Duration int64
	}
	startTime time.Time //Used for duration, not encoded
}

func NewLogglyEntry() *LogglyEntry {
	return &LogglyEntry{
		startTime: time.Now(),
		Level:     INFO,
	}
}

func (le *LogglyEntry) SetRequest(req *http.Request) {
	le.Request.Method = req.Method
	le.Request.Proto = req.Proto
	le.Request.Host = req.URL.Host
	le.Request.Path = req.URL.Path
	le.Request.Query = req.URL.RawQuery
	le.Request.Header = req.Header
	if req.Body != nil {
		b := new(bytes.Buffer)
		io.Copy(b, req.Body)
		le.Request.Body = b.String()
		req.Body = ioutil.NopCloser(b)
	}
}
func (le *LogglyEntry) SetResponse(status int, body interface{}) {
	le.Response.Body = body
	le.Response.Status = status
	le.Response.Duration = time.Since(le.startTime).Nanoseconds() / 1000000 //1 ms = 1000000ns
}

func (le *LogglyEntry) Log(msg ...interface{}) {
	le.Level.Set(INFO)
	le.Logs.Log(msg...)
}
func (le *LogglyEntry) Logf(msg string, args ...interface{}) {
	le.Level.Set(INFO)
	le.Logs.Logf(msg, args...)
}
func (le *LogglyEntry) Error(err error) {
	le.Level.Set(ERROR)
	le.Logs.Error(err)
}
func (le *LogglyEntry) Errorf(msg string, args ...interface{}) {
	le.Level.Set(ERROR)
	le.Logs.Errorf(msg, args...)
}
func (le *LogglyEntry) Panic(msg interface{}) {
	le.Level.Set(PANIC)
	le.Logs.Panic(msg)
}
func (le *LogglyEntry) Panicf(msg string, args ...interface{}) {
	le.Level.Set(PANIC)
	le.Logs.Panicf(msg, args...)
}
