package configfile

import (
	"errors"
	"fmt"

	"gopkg.in/ini.v1"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type ProjectValidator func(dpl.Project) error

type ComponentValidator func(dpl.Component) error

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

func validateComponent(component dpl.Component) error {
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

func validateProject(project dpl.Project) error {
	for name, validator := range projectValidators {
		err := validator(project)
		if err != nil {
			return errors.New(fmt.Sprintf("%v [%v]", err, name))
		}
	}
	return nil
}

type IniComponent struct {
	config *ini.Section
}

func (ic *IniComponent) Name() string {
	return ic.config.Name()
}

func (ic *IniComponent) GetValue(name string) []string {
	if ic.config.HasKey(name) {
		return ic.config.Key(name).ValueWithShadows()
	}
	return nil
}

type IniProject struct {
	config *ini.File
}

func (ip *IniProject) GetComponent(name string) (dpl.Component, bool) {
	if name == ini.DEFAULT_SECTION {
		return nil, false
	}
	component, err := ip.config.GetSection(name)
	if err != nil {
		return nil, false
	}
	return &IniComponent{
		config: component,
	}, true
}

func (ip *IniProject) ComponentNames() []string {
	return ip.config.SectionStrings()
}

func (ip *IniProject) Components() []string {
	return ip.config.SectionStrings()[1:]
}

func applyConfig(config *ini.File) (dpl.Project, error) {
	project := IniProject{
		config: config,
	}
	for _, component := range config.Sections() {
		if component.Name() != ini.DEFAULT_SECTION {
			projectComponent := IniComponent{
				config: component,
			}
			err := validateComponent(&projectComponent)
			if err != nil {
				return nil, err
			}
		}
	}
	err := validateProject(&project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func LoadRawConfig(data []byte) (dpl.Project, error) {
	config, err := ini.ShadowLoad(data)
	if err != nil {
		return nil, err
	}
	return applyConfig(config)
}

func LoadProjectConfig(path string) (dpl.Project, error) {
	config, err := ini.ShadowLoad(path)

	if err != nil {
		return nil, err
	}
	return applyConfig(config)
}
