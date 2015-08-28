package log

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var blacklisted []string

func Blacklist(strs ...string) {
	blacklisted = append(blacklisted, strs...)
}

func filterRequest(body string) string {
	for _, entry := range blacklisted {
		if strings.Contains(body, entry) {
			return "redacted"
		}
	}
	return body
}

const MAX_BODY_SIZE = 10000

func truncateBody(body string) string {
	if len(body) > MAX_BODY_SIZE {
		return body[:MAX_BODY_SIZE] + "[truncated]"
	}
	return body
}

func stringifyBody(body interface{}) string {
	switch body := body.(type) {
	case string:
		return body + " "
	default:
		b, _ := json.Marshal(body)
		return string(b) + " "
	}
}

func requestBody(req *http.Request) string {
	b := new(bytes.Buffer)
	io.Copy(b, req.Body)
	req.Body = ioutil.NopCloser(b)

	rb := b.String()
	if !debug {
		rb = filterRequest(rb)
	}
	return truncateBody(rb)
}

func responseBody(body interface{}) string {
	return truncateBody(stringifyBody(body))
}
