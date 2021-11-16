package dpl

type Component interface {
	Name() string
	ValueNames() []string
	GetValue(string) []string
	ExpandValue(string) ([]string, error)
	SetValue(string, []string)
	EraseValue(string)
}

type Project interface {
	GetComponent(string) (Component, bool)
	Components() []string
}
