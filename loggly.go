package log

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	loggly "github.com/treetopllc/go-loggly"
)

type LogglyClient struct {
	Basic
	client loggly.Client
}

func NewLogglyClient(token string, tags ...string) *LogglyClient {
	c := &LogglyClient{
		client: loggly.NewClient(token, tags...),
	}
	c.Basic = NewBasic(c, "", 0)
	return c
}

func (lc *LogglyClient) Send(e *LogglyEntry) {
	if e != nil {
		b, err := json.Marshal(e) //no fields here can json err
		if err != nil {
			fmt.Println("Error sending to loggly: ", err)
		} else {
			lc.Write(b)
		}
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
		Proto       string
		Method      string
		Host        string
		Path        string
		Query       string
		Body        string
		Header      http.Header
		UserID      string
		ProductType string
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
		le.Request.Body = requestBody(req)
	}
}
func (le *LogglyEntry) SetUserID(id string) {
	le.Request.UserID = id
}
func (le *LogglyEntry) SetResponse(status int, body interface{}) {
	if debug {
		le.Response.Body = responseBody(body)
	}
	le.Response.Status = status
	le.Response.Duration = time.Since(le.startTime).Nanoseconds() / 1000000 //1 ms = 1000000ns
}

func (le *LogglyEntry) SetProductType(pt string) {
	le.Request.ProductType = pt
}

func (le *LogglyEntry) Write(b []byte) (int, error) {
	le.Level.Set(INFO)
	le.Logs.Log(string(b))
	return len(b), nil
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
