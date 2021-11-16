package dpl

import (
	"fmt"
	"regexp"
)

type InvalidComponentNameError struct {
	Name string
}

func (icne *InvalidComponentNameError) Error() string {
	return fmt.Sprintf("Invalid name: %v", icne.Name)
}

func validateComponentName(component Component) error {
	matched, err := regexp.Match("^([a-zA-Z](?:([-_])?[a-zA-Z0-9])+)+$", []byte(component.Name()))

	if err != nil {
		return err
	}
	if !matched {
		return &InvalidComponentNameError{
			Name: component.Name(),
		}
	}
	return nil
}

func init() {
	RegisterComponentValidator("component-name", validateComponentName)
}
