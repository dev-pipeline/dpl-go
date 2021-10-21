package resolve

import (
	"fmt"
)

type ComponentNotFoundError struct {
	Name string
}

type Resolver interface {
	Resolve(taskChannel chan []string)
	Complete(task string)
	Fail(task string) []string
}

func (cnfe *ComponentNotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find component: %v", cnfe.Name)
}

func makeComponentTask(component string, task string) string {
	return fmt.Sprintf("%v.%v", component, task)
}
