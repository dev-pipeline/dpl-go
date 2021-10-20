package dpl

import (
	"fmt"
	"strings"
)

type depSet map[string]struct{}
type reverseDependencies map[string]depSet

func makeComponentTask(component string, task string) string {
	return fmt.Sprintf("%v.%v", component, task)
}

func insertKey(componentTask string, reverseDeps reverseDependencies) {
	_, found := reverseDeps[componentTask]
	if !found {
		reverseDeps[componentTask] = make(depSet)
	}
}

func addDeps(project *Project, target string, tasks []string, reverseDeps reverseDependencies) error {
	component, _ := project.ComponentInfo[target]
	for index, task := range tasks {
		depKey := fmt.Sprintf("depends.%v", task)
		componentTask := makeComponentTask(target, task)
		insertKey(componentTask, reverseDeps)

		rawDepends, found := component.Data[depKey]
		if found {
			// we have dependencies
			splitDepends := strings.Split(rawDepends, ",")
			for _, depend := range splitDepends {
				addDeps(project, depend, tasks[:index+1], reverseDeps)
				dependsTask := makeComponentTask(depend, task)
				insertKey(dependsTask, reverseDeps)
				reverseDeps[dependsTask][componentTask] = struct{}{}
			}
		}
		if index > 0 {
			reverseDeps[makeComponentTask(target, tasks[index-1])][makeComponentTask(target, tasks[index])] = struct{}{}
		}
	}
	return nil
}

func makeReverseDependencies(project *Project, targets []string, tasks []string) (reverseDependencies, error) {
	reverseDeps := make(reverseDependencies)
	for _, target := range targets {
		addDeps(project, target, tasks, reverseDeps)
		/*
			for index, task := range tasks[1:] {
				addDeps(project, target, task, reverseDeps)
			}
		*/
	}
	return reverseDeps, nil
}
