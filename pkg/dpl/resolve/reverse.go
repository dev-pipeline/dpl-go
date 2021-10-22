package resolve

import (
	"errors"
	"strings"
	"sync"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type ReverseResolver struct {
	cond       *sync.Cond
	revDeps    reverseDependencies
	depCounts  map[string]int
	readyTasks []string
}

func (dr *ReverseResolver) Resolve(taskChannel chan []string) {
	go func() {
		stillWork := func() bool {
			return len(dr.revDeps) > 0 || len(dr.depCounts) > 0 || len(dr.readyTasks) > 0
		}
		dr.cond.L.Lock()
		for stillWork() {
			for len(dr.readyTasks) > 0 {
				toSend := dr.readyTasks
				dr.readyTasks = []string{}
				dr.cond.L.Unlock()
				taskChannel <- toSend
				dr.cond.L.Lock()
			}

			if stillWork() {
				dr.cond.Wait()
			}
		}
		dr.cond.L.Unlock()
		close(taskChannel)
	}()
}

func (dr *ReverseResolver) Complete(task string) {
	dr.cond.L.Lock()
	defer dr.cond.L.Unlock()

	rev, found := dr.revDeps[task]
	if found {
		for task := range rev {
			count, found := dr.depCounts[task]
			if found {
				if count == 1 {
					// this was the last blocking dependency
					dr.readyTasks = append(dr.readyTasks, task)
					delete(dr.depCounts, task)
				} else {
					dr.depCounts[task] = count - 1
				}
			}
		}
		delete(dr.revDeps, task)
	}

	dr.cond.Signal()
}

func (dr *ReverseResolver) failHelper(task string, failChain map[string]struct{}) {
	rev, found := dr.revDeps[task]
	if found {
		for task := range rev {
			dr.failHelper(task, failChain)
		}
		delete(dr.revDeps, task)
	}
	delete(dr.depCounts, task)
	failChain[task] = struct{}{}
}

func (dr *ReverseResolver) Fail(task string) []string {
	failChain := func() map[string]struct{} {
		dr.cond.L.Lock()
		defer dr.cond.L.Unlock()

		failChain := map[string]struct{}{}
		dr.failHelper(task, failChain)
		delete(failChain, task)
		dr.cond.Signal()
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

func addRevDep(fullDeps reverseDependencies, trimmedDeps reverseDependencies, target string, tasks []string) {
	for index, task := range tasks {
		componentTask := makeComponentTask(target, task)
		_, found := trimmedDeps[componentTask]
		if !found {
			trimmedDeps[componentTask] = fullDeps[componentTask]
			localRevDeps := trimmedDeps[componentTask]
			for revDep := range localRevDeps {
				addRevDep(fullDeps, trimmedDeps, strings.Split(revDep, ".")[0], tasks[index+1:])
			}
		}
	}
}

func trimReverseDependencies(fullDeps reverseDependencies, targets []string, tasks []string) reverseDependencies {
	required := make(reverseDependencies)
	for _, target := range targets {
		addRevDep(fullDeps, required, target, tasks)
	}
	return required
}

func validateReverseResolution(resolver *ReverseResolver) error {
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

func resolveReverse(project dpl.Project, targets []string, tasks []string) (*ReverseResolver, error) {
	revDeps, err := makeReverseDependencies(project, project.ComponentNames(), tasks)
	if err != nil {
		return nil, err
	}
	trimmedDeps := trimReverseDependencies(revDeps, targets, tasks)

	counts := make(map[string]int)
	for baseTask, revDep := range trimmedDeps {
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
	resolver := &ReverseResolver{
		cond:       sync.NewCond(&m),
		revDeps:    trimmedDeps,
		depCounts:  counts,
		readyTasks: ready,
	}
	err = validateReverseResolution(resolver)
	if err != nil {
		return nil, err
	}
	return resolver, nil

	return nil, nil
}

func init() {
	RegisterResolver("reverse", func(project dpl.Project, targets []string, tasks []string) (Resolver, error) {
		return resolveReverse(project, targets, tasks)
	})
}
