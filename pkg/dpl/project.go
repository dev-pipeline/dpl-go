package dpl

type Component interface {
	Name() string
	GetValue(string) []string
	ExpandValue(string) ([]string, error)
}

type Project interface {
	GetComponent(string) (Component, bool)
	ComponentNames() []string
	Components() []string
}
