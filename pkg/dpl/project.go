package dpl

type Component interface {
	Name() string
	GetValue(string) (string, bool)
}

type Project interface {
	GetComponent(string) (Component, bool)
	Components() []string
}
