package terraformcli

import (
	"fmt"
	"strings"
)

type options []option

type option struct {
	key   string
	value string
}

func (o *options) array() []string {
	array := []string{}
	for _, option := range *o {
		array = append(array, option.render())
	}
	return array
}

func (o *option) render() string {
	if emptyString(o.value) {
		return fmt.Sprintf("-%s", o.key)
	}
	return fmt.Sprintf("-%s=%s", o.key, o.value)
}

func emptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

func (t *Terraform) varFilesOptions(varFiles []string) *options {
	options := options{}

	for _, file := range varFiles {
		option := option{
			key:   "-var-file",
			value: file,
		}
		options = append(options, option)
	}

	return &options
}

func (t *Terraform) stateFileOption(stateFile string) *option {
	return &option{
		key:   "-state",
		value: stateFile,
	}
}
