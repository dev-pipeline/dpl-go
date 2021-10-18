package dpl

type Component struct {
	Name string
	Data map[string]string
}

type Project struct {
	Components map[string]*Component
}

func NewProject() *Project {
	return &Project{
		Components: make(map[string]*Component),
	}
}

func NewComponent(name string) *Component {
	return &Component{
		Name: name,
		Data: make(map[string]string),
	}
}
