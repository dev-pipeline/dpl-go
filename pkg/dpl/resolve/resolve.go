package resolve

import (
	"errors"
	"fmt"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type Resolver interface {
	Resolve(taskChannel chan []string)
	Complete(task string)
	Fail(task string) []string
}

type ResolveFn func(dpl.Project, []string, []string) (Resolver, error)

var (
	resolvers = map[string]ResolveFn{}
)

func RegisterResolver(name string, resolver ResolveFn) error {
	_, found := resolvers[name]
	if found {
		return errors.New(fmt.Sprintf("Resolver %v already registered", name))
	}
	resolvers[name] = resolver
	return nil
}

type ComponentNotFoundError struct {
	Name string
}

func (cnfe *ComponentNotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find component: %v", cnfe.Name)
}

func makeComponentTask(component string, task string) string {
	return fmt.Sprintf("%v.%v", component, task)
}
