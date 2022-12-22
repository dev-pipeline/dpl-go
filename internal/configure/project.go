package configure

import (
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"

	"gopkg.in/ini.v1"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

var (
	errCantFindComponent error = errors.New("couldn't find component")
	errMissingKey        error = errors.New("missing key")
	errTooManyExpansions error = errors.New("too many expansions")
)

type IniComponent struct {
	config  *ini.Section
	project *IniProject
}

func (ic *IniComponent) Name() string {
	return ic.config.Name()
}

func (ic *IniComponent) ValueNames() []string {
	return ic.config.KeyStrings()
}

func (ic *IniComponent) GetValue(name string) []string {
	if ic.config.HasKey(name) {
		return ic.config.Key(name).ValueWithShadows()
	} else {
		defaultComponent, err := ic.project.config.GetSection(ini.DefaultSection)
		if err != nil {
			// swallow error on purpose
			return nil
		}
		if defaultComponent.HasKey(name) {
			return defaultComponent.Key(name).ValueWithShadows()
		}
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
				return "", errCantFindComponent
			}
		}
		if iniComponent.HasKey(key) {
			rawKey := iniComponent.Key(key)
			if rawKey != nil {
				loc := pattern.FindStringIndex(value)
				value = fmt.Sprintf("%v%v%v%v", value[:loc[0]], prefix, rawKey.Value(), value[loc[1]:])
			}
		} else {
			return "", errMissingKey
		}
		count++
	}
	return "", errTooManyExpansions
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

func (ic *IniComponent) GetSourceDir() string {
	return path.Join(ic.project.srcDir, ic.Name())
}

func (ic *IniComponent) GetWorkDir() string {
	return path.Join(ic.project.workDir, ic.Name())
}

type IniProject struct {
	config  *ini.File
	workDir string
	srcDir  string
}

func (ip *IniProject) getDefaultComponent() (dpl.Component, error) {
	section, err := ip.config.GetSection(ini.DefaultSection)
	if err != nil {
		return nil, err
	}
	return &IniComponent{
		config:  section,
		project: ip,
	}, nil
}

func (ip *IniProject) getConfigComponent(name string) (*ini.Section, bool) {
	if name == ini.DefaultSection {
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

func (ip *IniProject) Write(writer io.Writer) error {
	_, err := ip.config.WriteTo(writer)
	return err
}
