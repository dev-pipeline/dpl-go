package resolve

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type depSet map[string]struct{}
type reverseDependencies map[string]depSet

var (
	exists struct{}
)

type commonResolver struct {
	cond       *sync.Cond
	revDeps    reverseDependencies
	depCounts  map[string]int
	readyTasks []string
}

func (cr *commonResolver) Resolve(taskChannel chan []string) {
	go func() {
		stillWork := func() bool {
			return len(cr.revDeps) > 0 || len(cr.depCounts) > 0 || len(cr.readyTasks) > 0
		}
		cr.cond.L.Lock()
		for stillWork() {
			for len(cr.readyTasks) > 0 {
				toSend := cr.readyTasks
				cr.readyTasks = []string{}
				cr.cond.L.Unlock()
				taskChannel <- toSend
				cr.cond.L.Lock()
			}

			if stillWork() {
				cr.cond.Wait()
			}
		}
		cr.cond.L.Unlock()
		close(taskChannel)
	}()
}

func (cr *commonResolver) Complete(task string) {
	cr.cond.L.Lock()
	defer cr.cond.L.Unlock()

	rev, found := cr.revDeps[task]
	if found {
		for task := range rev {
			count, found := cr.depCounts[task]
			if found {
				if count == 1 {
					// this was the last blocking dependency
					cr.readyTasks = append(cr.readyTasks, task)
					delete(cr.depCounts, task)
				} else {
					cr.depCounts[task] = count - 1
				}
			}
		}
		delete(cr.revDeps, task)
	}

	cr.cond.Signal()
}

func (cr *commonResolver) failHelper(task string, failChain map[string]struct{}) {
	rev, found := cr.revDeps[task]
	if found {
		for task := range rev {
			cr.failHelper(task, failChain)
		}
		delete(cr.revDeps, task)
	}
	delete(cr.depCounts, task)
	failChain[task] = struct{}{}
}

func (cr *commonResolver) Fail(task string) []string {
	failChain := func() map[string]struct{} {
		cr.cond.L.Lock()
		defer cr.cond.L.Unlock()

		failChain := map[string]struct{}{}
		cr.failHelper(task, failChain)
		delete(failChain, task)
		cr.cond.Signal()
		return failChain
	}()

	failures := make([]string, len(failChain))
	index := 0
	for failure := range failChain {
		failures[index] = failure
		index++
	}
	return failures
}

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

func deepCopyReverseDeps(revDeps reverseDependencies) reverseDependencies {
	ret := reverseDependencies{}
	for rev, deps := range revDeps {
		ret[rev] = depSet{}
		for dep := range deps {
			ret[rev][dep] = struct{}{}
		}
	}
	return ret
}

func deepCopyDepCounts(depCounts map[string]int) map[string]int {
	ret := map[string]int{}
	for task, count := range depCounts {
		ret[task] = count
	}
	return ret
}

func deepCopyReadyTasks(tasks []string) []string {
	ret := make([]string, len(tasks))
	for index, task := range tasks {
		ret[index] = task
	}
	return ret
}

func validateResolution(resolver *commonResolver) error {
	backupRevDeps := deepCopyReverseDeps(resolver.revDeps)
	backupCounts := deepCopyDepCounts(resolver.depCounts)
	backupReady := deepCopyReadyTasks(resolver.readyTasks)

	for len(resolver.readyTasks) > 0 {
		resolver.Complete(resolver.readyTasks[0])
		resolver.readyTasks = resolver.readyTasks[1:]
	}

	if len(resolver.revDeps) > 0 || len(resolver.depCounts) > 0 || len(resolver.readyTasks) > 0 {
		return errors.New("Unable to resolve targets")
	}

	resolver.revDeps = backupRevDeps
	resolver.depCounts = backupCounts
	resolver.readyTasks = backupReady
	return nil
}

func resolveCommon(revDeps reverseDependencies) (commonResolver, error) {
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

	ready := []string{}
	for name, count := range counts {
		if count == 0 {
			ready = append(ready, name)
		}
	}

	for _, name := range ready {
		delete(counts, name)
	}

	m := sync.Mutex{}
	resolver := commonResolver{
		cond:       sync.NewCond(&m),
		revDeps:    revDeps,
		depCounts:  counts,
		readyTasks: ready,
	}
	err := validateResolution(&resolver)
	if err != nil {
		return commonResolver{}, err
	}
	return resolver, nil
}
