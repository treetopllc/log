package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/treetopllc/elastilog"
)

type ElasticClient struct {
	Basic
	client elastilog.Client
}

func NewElastic(uri string, tags ...string) *ElasticClient {
	c := &ElasticClient{
		client: elastilog.NewClient(uri, tags...),
	}
	c.Basic = NewBasic(c, "", 0)
	return c
}

func (el ElasticClient) Send(e *ElasticEntry) {
	if e != nil {
		e.entry.Attributes["level"] = string(e.level)
		e.entry.Log = string(e.str)
		el.client.Send(e.entry)
	}
}

func (el ElasticClient) Write(msg []byte) (int, error) {
	hostname, _ := os.Hostname()
	el.client.Send(elastilog.Entry{
		Log:        string(msg),
		Host:       hostname,
		Timestamp:  time.Now(),
		Attributes: elastilog.Attributes{"level": "std"},
	})
	return len(msg), nil
}

type ElasticEntry struct {
	level LogLevel
	str   StringLogger
	entry elastilog.Entry
}

func NewElasticEntry() *ElasticEntry {
	hostname, _ := os.Hostname()
	return &ElasticEntry{
		level: INFO,
		entry: elastilog.Entry{
			Host:       hostname,
			Timestamp:  time.Now(),
			Attributes: make(elastilog.Attributes),
		},
	}
}

func (ee *ElasticEntry) set(key string, value string) {
	ee.entry.Attributes[key] = value
}

func (ee *ElasticEntry) SetRequest(req *http.Request) {
	ee.entry.Host = req.URL.Host
	ee.set("request.method", req.Method)
	ee.set("request.path", req.URL.Path)
	ee.set("request.query", req.URL.RawQuery)
	ee.set("request.proto", req.Proto)
	if req.Body != nil {
		b := new(bytes.Buffer)
		io.Copy(b, req.Body)
		rb := b.String()
		if !debug {
			rb = filterRequest(rb)
		}
		ee.set("request.body", rb+" ")
		req.Body = ioutil.NopCloser(b)
	}
	for k, h := range req.Header {
		ee.set("request.header."+k, strings.Join(h, ","))
	}
}
func (ee *ElasticEntry) SetUserID(id string) {
	ee.set("request.user_id", id)
}

func (ee *ElasticEntry) SetResponse(status int, body interface{}) {
	if debug || status >= 300 {
		switch body := body.(type) {
		case string:
			ee.set("response.body", fmt.Sprintf("%s ", body))
		default:
			b, _ := json.Marshal(body)
			ee.set("response.body", fmt.Sprintf("%s ", string(b)))
		}
	}
	ee.set("response.status", fmt.Sprintf("%v", status))
	ee.set("duration", fmt.Sprintf("%v", time.Since(ee.entry.Timestamp).Nanoseconds()/1000000)) //1 ms = 1000000ns
}

func (ee *ElasticEntry) Write(b []byte) (int, error) {
	ee.level.Set(INFO)
	ee.str.Log(string(b))
	return len(b), nil
}
func (ee *ElasticEntry) Log(msg ...interface{}) {
	ee.level.Set(INFO)
	ee.str.Log(msg...)
}
func (ee *ElasticEntry) Logf(msg string, args ...interface{}) {
	ee.level.Set(INFO)
	ee.str.Logf(msg, args...)
}
func (ee *ElasticEntry) Error(err error) {
	ee.level.Set(ERROR)
	ee.str.Error(err)
}
func (ee *ElasticEntry) Errorf(msg string, args ...interface{}) {
	ee.level.Set(ERROR)
	ee.str.Errorf(msg, args...)
}
func (ee *ElasticEntry) Panic(msg interface{}) {
	ee.level.Set(PANIC)
	ee.str.Panic(msg)
}
func (ee *ElasticEntry) Panicf(msg string, args ...interface{}) {
	ee.level.Set(PANIC)
	ee.str.Panicf(msg, args...)
}
