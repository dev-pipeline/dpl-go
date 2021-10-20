package dpl

import (
	"fmt"
	"strings"
)

type depSet map[string]bool
type reverseDependencies map[string]depSet

func makeComponentTask(component string, task string) string {
	return fmt.Sprintf("%v.%v", component, task)
}

func addDeps(project *Project, target string, task string, reverseDeps reverseDependencies) error {
	component, _ := project.ComponentInfo[target]
	depKey := fmt.Sprintf("depends.%v", task)
	componentTask := makeComponentTask(target, task)
	rawDepends, found := component.Data[depKey]
	if found {
		// we have dependencies
		splitDepends := strings.Split(rawDepends, ",")
		for _, depend := range splitDepends {
			addDeps(project, depend, task, reverseDeps)
			dependsTask := makeComponentTask(depend, task)
			reverseDeps[dependsTask][componentTask] = true
		}
	}
	_, found = reverseDeps[componentTask]
	if !found {
		reverseDeps[componentTask] = make(depSet)
	}
	return nil
}

func makeReverseDependencies(project *Project, targets []string, tasks []string) (reverseDependencies, error) {
	reverseDeps := make(reverseDependencies)
	for _, target := range targets {
		for _, task := range tasks {
			addDeps(project, target, task, reverseDeps)
		}
	}
	return reverseDeps, nil
}
