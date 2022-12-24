package build

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

var (
	envPattern *regexp.Regexp
)

func findEnvIndex(env []string, actualKey string) int {
	for i := range env {
		if strings.HasPrefix(env[i], actualKey) {
			return i
		}
	}
	return -1
}

func makeEnvString(name string, values []string) string {
	return fmt.Sprintf("%v=%v", name, strings.Join(values, string(os.PathListSeparator)))
}

func extractEnvValues(variable string, value string) []string {
	return strings.Split(value[len(variable)+1:], string(os.PathListSeparator))
}

func prependEnvironment(originalEnv []string, variable string, extra []string) []string {
	actualKey := fmt.Sprintf("%v=", variable)
	index := findEnvIndex(originalEnv, actualKey)
	if index != -1 {
		originalEnv[index] = makeEnvString(variable, append(extra, extractEnvValues(variable, originalEnv[index])...))
		return originalEnv
	}
	return append(originalEnv, makeEnvString(variable, extra))
}

func appendEnvironment(originalEnv []string, variable string, extra []string) []string {
	actualKey := fmt.Sprintf("%v=", variable)
	index := findEnvIndex(originalEnv, actualKey)
	if index != -1 {
		originalEnv[index] = makeEnvString(variable, append(extractEnvValues(variable, originalEnv[index]), extra...))
		return originalEnv
	}
	return append(originalEnv, makeEnvString(variable, extra))
}

type environmentMap map[string][]string

type environmentChanges struct {
	prependValues environmentMap
	appendValues  environmentMap
}

func makeEnvMap(component dpl.Component) (environmentChanges, error) {
	ret := environmentChanges{
		prependValues: environmentMap{},
		appendValues:  environmentMap{},
	}
	configKeys := component.ValueNames()
	for index := range configKeys {
		groups := envPattern.FindStringSubmatch(configKeys[index])
		if groups != nil {
			expandedValues, err := component.ExpandValue(configKeys[index])
			if err != nil {
				return environmentChanges{}, err
			}
			key := strings.ToUpper(groups[1])
			switch groups[2] {
			case "append":
				ret.appendValues[key] = expandedValues
			case "prepend":
				ret.prependValues[key] = expandedValues
			}
		}
	}
	return ret, nil
}

func init() {
	var err error
	envPattern, err = regexp.Compile(`env\.(.*)\.((?:prepend)|(?:append))`)
	if err != nil {
		log.Fatalf("Error building pattern: %v", err)
	}
}
