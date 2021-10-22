package resolve

import (
	"errors"
	"sync"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type DeepResolver struct {
	cond       *sync.Cond
	revDeps    reverseDependencies
	depCounts  map[string]int
	readyTasks []string
}

func (dr *DeepResolver) Resolve(taskChannel chan []string) {
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

func (dr *DeepResolver) Complete(task string) {
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

func (dr *DeepResolver) failHelper(task string, failChain map[string]struct{}) {
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

func (dr *DeepResolver) Fail(task string) []string {
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

func validateResolution(resolver *DeepResolver) error {
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

func resolveDeep(project dpl.Project, targets []string, tasks []string) (*DeepResolver, error) {
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
	resolver := &DeepResolver{
		cond:       sync.NewCond(&m),
		revDeps:    revDeps,
		depCounts:  counts,
		readyTasks: ready,
	}
	err = validateResolution(resolver)
	if err != nil {
		return nil, err
	}
	return resolver, nil
}

func init() {
	RegisterResolver("deep", func(project dpl.Project, targets []string, tasks []string) (Resolver, error) {
		return resolveDeep(project, targets, tasks)
	})
}
