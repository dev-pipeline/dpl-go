package configure

import (
	"gopkg.in/ini.v1"

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

func LoadModifierConfig(filename string, modSet configfile.ModifierSet) error {
	config, err := ini.ShadowLoad(filename)
	if err != nil {
		return err
	}

	sections := []shadowSection{
		shadowSection{
			name: "prepend",
			data: modSet.PrependValues,
		},
		shadowSection{
			name: "append",
			data: modSet.AppendValues,
		},
		shadowSection{
			name: "override",
			data: modSet.OverrideValues,
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
			modSet.EraseValues[key.Name()] = struct{}{}
		}
	}

	return nil
}
