package dpl

import (
	"errors"
	"fmt"

	"gopkg.in/ini.v1"
)

type ProjectValidator func(*Project) error

type ComponentValidator func(*Component) error

type ComponentValidationError struct {
	ValidatorName string
	OriginalError error
}

func (cve *ComponentValidationError) Error() string {
	return fmt.Sprintf("%v [%v]", cve.OriginalError, cve.ValidatorName)
}

var (
	componentValidators = make(map[string]ComponentValidator)
	projectValidators   = make(map[string]ProjectValidator)
)

func RegisterComponentValidator(name string, validator ComponentValidator) error {
	fmt.Printf("Registering check %v\n", name)
	componentValidators[name] = validator
	return nil
}

func RegisterProjectValidator(name string, validator ProjectValidator) error {
	projectValidators[name] = validator
	return nil
}

func validateComponent(component *Component) error {
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

func validateProject(project *Project) error {
	for name, validator := range projectValidators {
		err := validator(project)
		if err != nil {
			return errors.New(fmt.Sprintf("%v [%v]", err, name))
		}
	}
	return nil
}

func applyConfig(config *ini.File) (*Project, error) {
	project := NewProject()
	for _, component := range config.Sections() {
		if component.Name() != ini.DEFAULT_SECTION {
			projectComponent := NewComponent(component.Name())
			for _, key := range component.Keys() {
				projectComponent.Data[key.Name()] = key.Value()
			}
			err := validateComponent(projectComponent)
			if err != nil {
				return nil, err
			}
			project.Components[component.Name()] = projectComponent
		}
	}
	err := validateProject(project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func LoadRawConfig(data []byte) (*Project, error) {
	config, err := ini.Load(data)
	if err != nil {
		return nil, err
	}
	return applyConfig(config)
}

func LoadProjectConfig(path string) (*Project, error) {
	config, err := ini.Load(path)

	if err != nil {
		return nil, err
	}
	return applyConfig(config)
}
