package terraformcli

import (
	"fmt"
	"strings"
)

type options struct {
	array []option
}

type option struct {
	key   string
	value string
}

func (o *option) render() string {
	if emptyString(o.value) {
		return fmt.Sprintf("%s", o.key)
	}
	return fmt.Sprintf("%s=%s", o.key, o.value)
}

func emptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
