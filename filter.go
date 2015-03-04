package log

import "strings"

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
