package testcommon

import (
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type ResolveComponent struct {
	Data map[string][]string
}

func (rs *ResolveComponent) Name() string {
	return ""
}

func (rs *ResolveComponent) ValueNames() []string {
	return nil
}

func (rs *ResolveComponent) GetValue(key string) []string {
	value, found := rs.Data[key]
	if found {
		return value
	}
	return nil
}

func (rs *ResolveComponent) ExpandValue(key string) ([]string, error) {
	return rs.GetValue(key), nil
}

func (rs *ResolveComponent) SetValue(string, []string) {
}

func (rs *ResolveComponent) EraseValue(string) {
}

type ResolveComponents map[string]ResolveComponent

type ResolveProject struct {
	Comps ResolveComponents
}

func (rp *ResolveProject) GetComponent(name string) (dpl.Component, bool) {
	component, found := rp.Comps[name]
	if found {
		return &component, true
	}
	return nil, false
}

func (rp *ResolveProject) Components() []string {
	names := []string{}
	for name := range rp.Comps {
		names = append(names, name)
	}
	return names
}
