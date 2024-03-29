package configure

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

var (
	expandedPattern     *regexp.Regexp
	reservedNamePattern *regexp.Regexp

	errCantFindComponent     error = errors.New("couldn't find component")
	errMissingKey            error = errors.New("missing key")
	errTooManyValues         error = errors.New("too many values")
	errTooManyExpansions     error = errors.New("too many expansions")
	errMissingEscapeSequence error = errors.New("missing escape sequence")
)

type IniComponent struct {
	config  *ini.Section
	project *IniProject
}

func (ic *IniComponent) Name() string {
	return ic.config.Name()
}

func (ic *IniComponent) KeyNames() []string {
	return ic.config.KeyStrings()
}

func (ic *IniComponent) GetValues(name string) []string {
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

func (ic *IniComponent) GetSingleValue(name string) (string, error) {
	vals := ic.GetValues(name)
	switch len(vals) {
	case 0:
		return "", errMissingKey

	case 1:
		return vals[0], nil

	default:
		return "", errTooManyValues
	}
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
		return []string{value}, nil
	}
	prefix := groups[1]
	component := groups[2]
	key := groups[3]

	iniComponent := ic
	if len(component) != 0 {
		var err error
		iniComponent, err = ic.project.getConfigComponent(component)
		if err != nil {
			return nil, err
		}
	}
	if !iniComponent.config.HasKey(key) {
		if iniComponent != ic {
			return nil, errMissingKey
		}
		defaultComponent, err := iniComponent.project.getDefaultComponent()
		if err != nil || !defaultComponent.config.HasKey(key) {
			return nil, errMissingKey
		}
		iniComponent = defaultComponent
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

var (
	specialEscapes map[byte]string = map[byte]string{
		'e': "",
	}
)

func escapeValue(value string) (string, error) {
	sb := strings.Builder{}
	lastIndex := 0
	for {
		ss := value[lastIndex:]
		index := strings.Index(ss, `\`)
		if index == -1 {
			sb.WriteString(ss)
			return sb.String(), nil
		}
		if index == len(value)-1 {
			return "", errMissingEscapeSequence
		} else {
			sb.WriteString(ss[:index])
			if specialString, found := specialEscapes[value[index+1]]; found {
				sb.WriteString(specialString)
			} else {
				sb.WriteByte(value[index+1])
			}
			lastIndex += 2 + index
		}
	}
}

func (ic *IniComponent) expandValueInternal(name string, count int) ([]string, error) {
	rawValues := ic.GetValues(name)
	if rawValues == nil {
		return nil, nil
	}
	ret := []string{}
	for index := range rawValues {
		expandedValues, err := ic.expandRecursively(rawValues[index], count)
		if err != nil {
			return nil, err
		}
		for index := range expandedValues {
			expandedValues[index], err = escapeValue(expandedValues[index])
			if err != nil {
				return nil, err
			}
		}
		ret = append(ret, expandedValues...)
	}
	return ret, nil
}

func (ic *IniComponent) ExpandValues(name string) ([]string, error) {
	return ic.expandValueInternal(name, 0)
}

func (ic *IniComponent) SetValues(name string, values []string) {
	ic.EraseKey(name)
	if len(values) > 0 {
		key, _ := ic.config.NewKey(name, values[0])
		for _, value := range values[1:] {
			key.AddShadow(value)
		}
	}
}

func (ic *IniComponent) EraseKey(name string) {
	ic.config.DeleteKey(name)
}

func (ic *IniComponent) GetSourceDir() string {
	srcDir, err := dpl.GetSingleComponentValue(ic, sourceDirKey)
	if err != nil {
		log.Fatalf("Error getting %v's source dir: %v", ic.Name(), err)
	}
	return srcDir
}

func (ic *IniComponent) GetWorkDir() string {
	workDir, err := dpl.GetSingleComponentValue(ic, workDirKey)
	if err != nil {
		log.Fatalf("Error getting %v's work dir: %v", ic.Name(), err)
	}
	return workDir
}

type IniProject struct {
	config     *ini.File
	configFile string
	dirty      bool
}

func (ip *IniProject) getAnyComponent(name string) (*IniComponent, error) {
	section, err := ip.config.GetSection(name)
	if err != nil {
		return nil, errCantFindComponent
	}
	return &IniComponent{
		config:  section,
		project: ip,
	}, nil
}

func (ip *IniProject) getDefaultComponent() (*IniComponent, error) {
	return ip.getAnyComponent(ini.DefaultSection)
}

func (ip *IniProject) getConfigComponent(name string) (*IniComponent, error) {
	if name == ini.DefaultSection {
		return nil, errCantFindComponent
	}
	if strings.HasPrefix(name, "dpl.") {
		return nil, errCantFindComponent
	}
	component, err := ip.config.GetSection(name)
	if err != nil {
		return nil, err
	}
	return &IniComponent{
		config:  component,
		project: ip,
	}, nil
}

func (ip *IniProject) GetComponent(name string) (dpl.Component, error) {
	return ip.getConfigComponent(name)
}

func (ip *IniProject) ComponentNames() []string {
	ret := []string{}
	rawNames := ip.config.SectionStrings()[1:]
	for i := range rawNames {
		if !reservedNamePattern.MatchString(rawNames[i]) {
			ret = append(ret, rawNames[i])
		}
	}
	return ret
}

func (ip *IniProject) Write() error {
	outConfig, err := os.Create(ip.configFile)
	if err != nil {
		return err
	}
	defer outConfig.Close()
	return ip.write(outConfig)
}

func (ip *IniProject) write(writer io.Writer) error {
	_, err := ip.config.WriteTo(writer)
	if err == nil {
		ip.dirty = false
	}
	return err
}

func init() {
	var err error
	expandedPattern, err = regexp.Compile(`(^|[^\\])(?:\${(?:([a-z_\-]+):)?((?:[a-zA-Z0-9_]+\.)*[a-zA-Z0-9)_]+)})`)
	if err != nil {
		log.Fatalf("Error compiling pattern: %v", err)
	}
	reservedNamePattern, err = regexp.Compile(`^dpl\..+`)
	if err != nil {
		log.Fatalf("Error compiling pattern: %v", err)
	}
}
