package util

import (
	"net/url"
	"strings"
)

func FormQuery(queryArr map[string]string) string {
	var qu []string
	for key, value := range queryArr {
		qu = append(qu, key+"="+value)
	}
	query, _ := url.ParseQuery(strings.Join(qu, "&"))
	return query.Encode()
}
