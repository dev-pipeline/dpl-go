package resolve

import (
	"fmt"
	"strings"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type depSet map[string]struct{}
type reverseDependencies map[string]depSet

var (
	exists struct{}
)

func insertKey(componentTask string, reverseDeps reverseDependencies) bool {
	_, found := reverseDeps[componentTask]
	if !found {
		reverseDeps[componentTask] = make(depSet)
		return true
	}
	return false
}

func addDeps(project dpl.Project, target string, tasks []string, reverseDeps reverseDependencies) error {
	for index, task := range tasks {
		componentTask := makeComponentTask(target, task)
		firstTime := insertKey(componentTask, reverseDeps)

		if firstTime {
			component, found := project.GetComponent(target)
			if !found {
				return &ComponentNotFoundError{
					Name: target,
				}
			}
			depKey := fmt.Sprintf("depends.%v", task)
			rawDepends, found := component.GetValue(depKey)
			if found {
				// we have dependencies
				splitDepends := strings.Split(rawDepends, ",")
				for _, depend := range splitDepends {
					err := addDeps(project, depend, tasks[:index+1], reverseDeps)
					if err != nil {
						return err
					}
					dependsTask := makeComponentTask(depend, task)
					insertKey(dependsTask, reverseDeps)
					reverseDeps[dependsTask][componentTask] = struct{}{}
				}
			}
			if index > 0 {
				reverseDeps[makeComponentTask(target, tasks[index-1])][makeComponentTask(target, tasks[index])] = struct{}{}
			}
		}
	}
	return nil
}

func makeReverseDependencies(project dpl.Project, targets []string, tasks []string) (reverseDependencies, error) {
	reverseDeps := make(reverseDependencies)
	for _, target := range targets {
		err := addDeps(project, target, tasks, reverseDeps)
		if err != nil {
			return nil, err
		}
	}
	return reverseDeps, nil
}
