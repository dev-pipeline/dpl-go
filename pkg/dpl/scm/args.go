package scm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	argumentPattern *regexp.Regexp
)

func extractArguments(raw string, separator string) (string, map[string]string, error) {
	chunks := strings.Split(raw, separator)
	args := map[string]string{}
	for _, arg := range chunks[1:] {
		groups := argumentPattern.FindStringSubmatch(arg)
		if groups == nil {
			return "", nil, errors.New(fmt.Sprintf("Couldn't parse '%v'", arg))
		}
		args[groups[1]] = groups[2]
	}
	return chunks[0], args, nil
}

func init() {
	var err error
	argumentPattern, err = regexp.Compile(`(.+)=(.+)`)
	if err != nil {
		panic(err)
	}
}
