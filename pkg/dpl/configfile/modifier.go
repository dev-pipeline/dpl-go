package configfile

import (
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type ModifierSet struct {
	prependValues  map[string][]string
	appendValues   map[string][]string
	overrideValues map[string][]string
	eraseValues    map[string]struct{}
}

func NewModifierSet() ModifierSet {
	return ModifierSet{
		prependValues:  make(map[string][]string),
		appendValues:   make(map[string][]string),
		overrideValues: make(map[string][]string),
		eraseValues:    make(map[string]struct{}),
	}
}

type ProjectModifiers struct {
	components map[string]ModifierSet
}

func ApplyComponentModifiers(component dpl.Component, modifiers ModifierSet) {
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
}
