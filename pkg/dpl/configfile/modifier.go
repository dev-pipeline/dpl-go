package configfile

import (
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type ModifierSet struct {
	PrependValues  map[string][]string
	AppendValues   map[string][]string
	OverrideValues map[string][]string
	EraseValues    map[string]struct{}
}

func NewModifierSet() ModifierSet {
	return ModifierSet{
		PrependValues:  make(map[string][]string),
		AppendValues:   make(map[string][]string),
		OverrideValues: make(map[string][]string),
		EraseValues:    make(map[string]struct{}),
	}
}

type ProjectModifiers struct {
	components map[string]ModifierSet
}

func ApplyComponentModifiers(component dpl.Component, modifiers ModifierSet) error {
	for key, prepends := range modifiers.PrependValues {
		originalValues := component.GetValue(key)
		component.SetValue(key, append(prepends, originalValues...))
	}

	for key, appends := range modifiers.AppendValues {
		originalValues := component.GetValue(key)
		component.SetValue(key, append(originalValues, appends...))
	}

	for key, overrides := range modifiers.OverrideValues {
		component.SetValue(key, overrides)
	}

	for key, _ := range modifiers.EraseValues {
		component.EraseValue(key)
	}
	return nil
}
