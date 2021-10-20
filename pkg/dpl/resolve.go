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

func addDeps(project Project, target string, tasks []string, reverseDeps reverseDependencies) error {
	component, _ := project.GetComponent(target)
	for index, task := range tasks {
		depKey := fmt.Sprintf("depends.%v", task)
		componentTask := makeComponentTask(target, task)
		insertKey(componentTask, reverseDeps)

		rawDepends, found := component.GetValue(depKey)
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

func makeReverseDependencies(project Project, targets []string, tasks []string) (reverseDependencies, error) {
	reverseDeps := make(reverseDependencies)
	for _, target := range targets {
		addDeps(project, target, tasks, reverseDeps)
	}
	return reverseDeps, nil
}

type DeepResolver struct {
	revDeps    reverseDependencies
	depCounts  map[string]int
	readyTasks []string
}

func ResolveDeep(project Project, targets []string, tasks []string) (*DeepResolver, error) {
	revDeps, err := makeReverseDependencies(project, targets, tasks)
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for baseTask, revDep := range revDeps {
		_, found := counts[baseTask]
		if !found {
			counts[baseTask] = 0
		}
		for task := range revDep {
			count, found := counts[task]
			if found {
				counts[task] = count + 1
			} else {
				counts[task] = 1
			}
		}
	}
	fmt.Printf("%v", counts)

	ready := []string{}
	for name, count := range counts {
		if count == 0 {
			ready = append(ready, name)
		}
	}

	for _, name := range ready {
		delete(counts, name)
	}

	return &DeepResolver{
		revDeps:    revDeps,
		depCounts:  counts,
		readyTasks: ready,
	}, nil
}
