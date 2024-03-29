package dpl

import (
	"fmt"
	"log"
	"regexp"
)

var (
	legalComponentName *regexp.Regexp
	legalFieldName     *regexp.Regexp
	dplPrefix          *regexp.Regexp
)

type InvalidComponentNameError struct {
	Name string
}

func (icne *InvalidComponentNameError) Error() string {
	return fmt.Sprintf("invalid name: %v", icne.Name)
}

func validateComponentName(component Component) error {
	matched := legalComponentName.Match([]byte(component.Name()))

	if !matched {
		return &InvalidComponentNameError{
			Name: component.Name(),
		}
	}
	return nil
}

type InvalidFieldNameError struct {
	Name string
}

func (ifne *InvalidFieldNameError) Error() string {
	return fmt.Sprintf("invalid key name: %v", ifne.Name)
}

func validateFieldName(component Component) error {
	for _, name := range component.KeyNames() {
		matched := legalFieldName.Match([]byte(name))
		if !matched {
			return &InvalidFieldNameError{
				Name: name,
			}
		}

		matched = dplPrefix.Match([]byte(name))
		if matched {
			return &InvalidFieldNameError{
				Name: name,
			}
		}
	}
	return nil
}

func init() {
	var err error
	legalComponentName, err = regexp.Compile("^([a-zA-Z](?:([-_])?[a-zA-Z0-9])+)+$")
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	legalFieldName, err = regexp.Compile(`^([a-z][a-z0-9_]*(?:\.)?)+$`)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	dplPrefix, err = regexp.Compile(`^dpl\.`)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	RegisterComponentValidator("component-name", validateComponentName)
	RegisterComponentValidator("field-name", validateFieldName)
}
