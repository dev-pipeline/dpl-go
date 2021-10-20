package dpl

type Component struct {
	Name string
	Data map[string]string
}

type Components map[string]*Component

type Project struct {
	ComponentInfo Components
}

func NewProject() *Project {
	return &Project{
		ComponentInfo: make(Components),
	}
}

func NewComponent(name string) *Component {
	return &Component{
		Name: name,
		Data: make(map[string]string),
	}
}
