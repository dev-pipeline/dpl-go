package testcommon

import (
	"fmt"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

var (
	errMissingComponent error = fmt.Errorf("missing component")
)

type ResolveComponent struct {
	Data map[string][]string
}

func (rs *ResolveComponent) Name() string {
	return ""
}

func (rs *ResolveComponent) KeyNames() []string {
	return nil
}

func (rs *ResolveComponent) GetValues(key string) []string {
	value, found := rs.Data[key]
	if found {
		return value
	}
	return nil
}

func (rs *ResolveComponent) ExpandValues(key string) ([]string, error) {
	return rs.GetValues(key), nil
}

func (rs *ResolveComponent) SetValues(string, []string) {
}

func (rs *ResolveComponent) EraseKey(string) {
}

func (rs *ResolveComponent) GetSourceDir() string {
	return ""
}

func (rs *ResolveComponent) GetWorkDir() string {
	return ""
}

type ResolveComponents map[string]ResolveComponent

type ResolveProject struct {
	Comps ResolveComponents
}

func (rp *ResolveProject) GetComponent(name string) (dpl.Component, error) {
	component, found := rp.Comps[name]
	if found {
		return &component, nil
	}
	return nil, errMissingComponent
}

func (rp *ResolveProject) ComponentNames() []string {
	names := []string{}
	for name := range rp.Comps {
		names = append(names, name)
	}
	return names
}

func (rp *ResolveProject) Write() error {
	return nil
}
