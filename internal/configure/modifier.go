package configure

import (
	"gopkg.in/ini.v1"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/configfile"
)

func applyShadowSection(section *ini.Section, data map[string][]string) error {
	for _, key := range section.Keys() {
		data[key.Name()] = append(data[key.Name()], key.ValueWithShadows()...)
	}
	return nil
}

type shadowSection struct {
	name string
	data map[string][]string
}

func loadModifierConfig(filename string, modSet modifierSet) error {
	config, err := configfile.LoadProjectConfig(filename)
	if err != nil {
		return err
	}

	sections := []shadowSection{
		shadowSection{
			name: "prepend",
			data: modSet.prependValues,
		},
		shadowSection{
			name: "append",
			data: modSet.appendValues,
		},
		shadowSection{
			name: "override",
			data: modSet.overrideValues,
		},
	}

	for _, ss := range sections {
		section := config.Section(ss.name)
		if section != nil {
			err = applyShadowSection(section, ss.data)
			if err != nil {
				return err
			}
		}
	}

	section := config.Section("erase")
	if section != nil {
		for _, key := range section.Keys() {
			modSet.eraseValues[key.Name()] = struct{}{}
		}
	}

	return nil
}

type modifierSet struct {
	prependValues  map[string][]string
	appendValues   map[string][]string
	overrideValues map[string][]string
	eraseValues    map[string]struct{}
}

func newModifierSet() modifierSet {
	return modifierSet{
		prependValues:  make(map[string][]string),
		appendValues:   make(map[string][]string),
		overrideValues: make(map[string][]string),
		eraseValues:    make(map[string]struct{}),
	}
}

type projectModifiers struct {
	components map[string]modifierSet
}

func applyComponentModifiers(component dpl.Component, modifiers modifierSet) error {
	for key, prepends := range modifiers.prependValues {
		originalValues := component.GetValue(key)
		component.SetValue(key, append(prepends, originalValues...))
	}

	for key, appends := range modifiers.appendValues {
		originalValues := component.GetValue(key)
		component.SetValue(key, append(originalValues, appends...))
	}

	for key, overrides := range modifiers.overrideValues {
		component.SetValue(key, overrides)
	}

	for key, _ := range modifiers.eraseValues {
		component.EraseValue(key)
	}
	return nil
}
