package dpl

import (
	"fmt"
)

type ProjectValidator func(Project) error

type ComponentValidator func(Component) error

type ComponentValidationError struct {
	ValidatorName string
	OriginalError error
}

func (cve ComponentValidationError) Error() string {
	return fmt.Sprintf("%v [%v]", cve.OriginalError, cve.ValidatorName)
}

var (
	componentValidators = make(map[string]ComponentValidator)
	projectValidators   = make(map[string]ProjectValidator)
)

func RegisterComponentValidator(name string, validator ComponentValidator) error {
	componentValidators[name] = validator
	return nil
}

func RegisterProjectValidator(name string, validator ProjectValidator) error {
	projectValidators[name] = validator
	return nil
}

func ValidateComponent(component Component) error {
	for name, validator := range componentValidators {
		err := validator(component)
		if err != nil {
			return &ComponentValidationError{
				ValidatorName: name,
				OriginalError: err,
			}
		}
	}
	return nil
}

type projectValidationError struct {
	name string
	err  error
}

func (pve *projectValidationError) Error() string {
	return fmt.Sprintf("%v [%v]", pve.err, pve.name)
}

func ValidateProject(project Project) error {
	for name, validator := range projectValidators {
		err := validator(project)
		if err != nil {
			return &projectValidationError{
				name: name,
				err:  err,
			}
		}
	}
	return nil
}
