package log

import "strings"

var blacklisted = []string{
	"email",
	"postal_code",
	"password",
	"birthdate",
	"gender",
	"last_name",
}

func filterRequest(body string) string {
	for _, entry := range blacklisted {
		if strings.Contains(body, entry) {
			return "redacted"
		}
	}
	return body
}
