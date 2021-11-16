package configfile

import (
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/ini.v1"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type IniComponent struct {
	config  *ini.Section
	project *IniProject
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

const (
	expandLimit = 100
)

func (ic *IniComponent) expandHelper(value string) (string, error) {
	pattern, err := regexp.Compile(`(^|[^\\])(?:\${(?:([a-z]+)\.)?([a-z]+)})`)
	if err != nil {
		return "", err
	}

	count := 0
	for count < expandLimit {
		groups := pattern.FindStringSubmatch(value)

		if groups == nil {
			// done expanding
			return value, nil
		}
		prefix := groups[1]
		component := groups[2]
		key := groups[3]

		iniComponent := ic.config
		if len(component) != 0 {
			var found bool
			iniComponent, found = ic.project.getConfigComponent(component)
			if !found {
				return "", errors.New("Couldn't find component")
			}
		}
		if iniComponent.HasKey(key) {
			rawKey := iniComponent.Key(key)
			if rawKey != nil {
				loc := pattern.FindStringIndex(value)
				value = fmt.Sprintf("%v%v%v%v", value[:loc[0]], prefix, rawKey.Value(), value[loc[1]:])
			}
		} else {
			return "", errors.New("Missing key")
		}
		count++
	}
	return "", errors.New("Too many expansions")
}

func (ic *IniComponent) ExpandValue(name string) ([]string, error) {
	rawValues := ic.GetValue(name)
	if rawValues == nil {
		return nil, nil
	}
	expandedValues := make([]string, len(rawValues))
	for index, value := range rawValues {
		expandedValue, err := ic.expandHelper(value)
		if err != nil {
			return nil, err
		}
		expandedValues[index] = expandedValue
	}
	return expandedValues, nil
}

func (ic *IniComponent) SetValue(name string, values []string) {
	ic.EraseValue(name)
	if len(values) > 0 {
		key, _ := ic.config.NewKey(name, values[0])
		for _, value := range values[1:] {
			key.AddShadow(value)
		}
	}
}

func (ic *IniComponent) EraseValue(name string) {
	ic.config.DeleteKey(name)
}

type IniProject struct {
	config *ini.File
}

func (ip *IniProject) getConfigComponent(name string) (*ini.Section, bool) {
	if name == ini.DEFAULT_SECTION {
		return nil, false
	}
	component, err := ip.config.GetSection(name)
	if err != nil {
		return nil, false
	}
	return component, true
}

func (ip *IniProject) GetComponent(name string) (dpl.Component, bool) {
	component, found := ip.getConfigComponent(name)
	if !found {
		return nil, false
	}
	return &IniComponent{
		config:  component,
		project: ip,
	}, true
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
			err := dpl.ValidateComponent(&projectComponent)
			if err != nil {
				return nil, err
			}
		}
	}
	err := dpl.ValidateProject(&project)
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
