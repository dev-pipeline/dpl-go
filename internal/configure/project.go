package configure

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path"
	"regexp"

	"gopkg.in/ini.v1"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

var (
	expandedPattern *regexp.Regexp

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

func (ic *IniComponent) expandRecursively(value string, count int) ([]string, error) {
	if count > expandLimit {
		return nil, errTooManyExpansions
	}
	groups := expandedPattern.FindStringSubmatch(value)
	if groups == nil {
		// done expanding
		if value == "<empty>" {
			return []string{""}, nil
		}
		return []string{value}, nil
	}
	prefix := groups[1]
	component := groups[2]
	key := groups[3]

	iniComponent := ic
	if len(component) != 0 {
		var found bool
		iniComponent, found = ic.project.getConfigComponent(component)
		if !found {
			return nil, errCantFindComponent
		}
	}
	if !iniComponent.config.HasKey(key) {
		return nil, errMissingKey
	}
	expanded, err := iniComponent.expandValueInternal(key, count+1)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for i := range expanded {
		loc := expandedPattern.FindStringIndex(value)
		nextValue := fmt.Sprintf("%v%v%v%v", value[:loc[0]], prefix, expanded[i], value[loc[1]:])
		nextExpanded, err := ic.expandRecursively(nextValue, count+1)
		if err != nil {
			return nil, err
		}
		ret = append(ret, nextExpanded...)
	}
	return ret, nil
}

func (ic *IniComponent) expandValueInternal(name string, count int) ([]string, error) {
	rawValues := ic.GetValue(name)
	if rawValues == nil {
		return nil, nil
	}
	ret := []string{}
	for index := range rawValues {
		expandedValues, err := ic.expandRecursively(rawValues[index], count)
		if err != nil {
			return nil, err
		}
		ret = append(ret, expandedValues...)
	}
	return ret, nil
}

func (ic *IniComponent) ExpandValue(name string) ([]string, error) {
	return ic.expandValueInternal(name, 0)
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

func (ip *IniProject) getConfigComponent(name string) (*IniComponent, bool) {
	if name == ini.DefaultSection {
		return nil, false
	}
	component, err := ip.config.GetSection(name)
	if err != nil {
		return nil, false
	}
	return &IniComponent{
		config:  component,
		project: ip,
	}, true
}

func (ip *IniProject) GetComponent(name string) (dpl.Component, bool) {
	return ip.getConfigComponent(name)
}

func (ip *IniProject) Components() []string {
	return ip.config.SectionStrings()[1:]
}

func (ip *IniProject) Write(writer io.Writer) error {
	_, err := ip.config.WriteTo(writer)
	return err
}

func init() {
	var err error
	expandedPattern, err = regexp.Compile(`(^|[^\\])(?:\${(?:([a-z_\-]+):)?((?:[a-zA-Z0-9_]+\.)*[a-zA-Z0-9)_]+)})`)
	if err != nil {
		log.Fatalf("Error compiling pattern: %v", err)
	}
}
