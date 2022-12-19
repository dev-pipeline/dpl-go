package common

import (
	"log"
	"strings"
	"sync"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/resolve"
)

type Args struct {
	KeepGoing    bool
	Executor     string
	Dependencies string
}

type taskFn func(dpl.Component) error

type Task struct {
	Name string
	Work taskFn
}

type work struct {
	fn        taskFn
	name      string
	component dpl.Component
}

type taskComplete struct {
	name string
	err  error
}

func executeTasks(wg *sync.WaitGroup, taskChannel chan work, doneChannel chan taskComplete) error {
	wg.Add(1)
	go func() {
		for {
			workUnit := <-taskChannel
			if len(workUnit.name) == 0 {
				wg.Done()
				return
			}
			err := workUnit.fn(workUnit.component)
			doneChannel <- taskComplete{
				name: workUnit.name,
				err:  err,
			}
		}
	}()

	return nil
}

func makeTaskContainers(tasks []Task) ([]string, map[string]taskFn, error) {
	taskList := []string{}
	taskMap := map[string]taskFn{}

	for _, task := range tasks {
		taskMap[task.Name] = task.Work
		taskList = append(taskList, task.Name)
	}

	return taskList, taskMap, nil
}

type taskCompleteFn func(taskComplete)

func startDrainComplete(doneChannel chan taskComplete, completeFn taskCompleteFn) {
	go func() {
		for {
			completedTask := <-doneChannel
			if len(completedTask.name) != 0 {
				completeFn(completedTask)
			} else {
				return
			}
		}
	}()
}

func startResolve(project dpl.Project, resolver resolve.Resolver, taskMap map[string]taskFn) (chan work, error) {
	workChannel := make(chan work)
	go func() {
		defer close(workChannel)

		readyTaskChannel := make(chan []string)
		resolver.Resolve(readyTaskChannel)
		readyTasks := <-readyTaskChannel
		for readyTasks != nil {
			for _, taskToExecute := range readyTasks {
				taskChunks := strings.Split(taskToExecute, ".")
				if len(taskChunks) != 2 {
					log.Fatalf("Internal error: improper extraction of task '%v'", taskChunks)
				}
				component, found := project.GetComponent(taskChunks[0])
				if !found {
					log.Fatalf("Internal error: cannot get component %v", taskChunks[0])
				}
				workFn, found := taskMap[taskChunks[1]]
				if !found {
					log.Fatalf("Internal error: no handler for task %v", taskChunks[1])
				}

				workChannel <- work{
					fn:        workFn,
					name:      taskToExecute,
					component: component,
				}
			}
			readyTasks = <-readyTaskChannel
		}
	}()
	return workChannel, nil
}

type failedTask struct {
	originalError error
	name          string
	dependents    []string
}

func (ft *failedTask) Error() string {
	return ""
}

func runTasks(project dpl.Project, components []string, tasks []Task, resolveFn resolve.ResolveFn, keepGoing bool) error {
	taskList, taskMap, err := makeTaskContainers(tasks)
	if err != nil {
		return err
	}

	resolver, err := resolveFn(project, components, taskList)
	if err != nil {
		return err
	}

	workChannel, err := startResolve(project, resolver, taskMap)
	if err != nil {
		return err
	}

	doneChannel := make(chan taskComplete)
	defer close(doneChannel)
	wg := sync.WaitGroup{}

	err = executeTasks(&wg, workChannel, doneChannel)
	if err != nil {
		return err
	}

	errors := []error{}
	m := sync.Mutex{}
	completeTask := func(completedTask taskComplete) {
		if completedTask.err == nil {
			resolver.Complete(completedTask.name)
		} else {
			var dependents []string
			if keepGoing {
				dependents = resolver.Fail(completedTask.name)
			} else {
				dependents, _ = resolver.Abort()
			}
			m.Lock()
			defer m.Unlock()
			errors = append(errors, &failedTask{
				originalError: completedTask.err,
				name:          completedTask.name,
				dependents:    dependents,
			})
		}
	}
	startDrainComplete(doneChannel, completeTask)

	wg.Wait()

	m.Lock()
	defer m.Unlock()
	if len(errors) == 0 {
		return nil
	}
	return errors[0]
}

func DoCommand(components []string, args Args, tasks []Task) {
	var project dpl.Project = nil // TODO: load project

	if len(components) == 0 {
		components = project.Components()
	}
	resolveFn := resolve.GetResolver(args.Dependencies)
	if resolveFn == nil {
		log.Fatalf("No resolver '%v'", args.Dependencies)
	} else {
		runTasks(project, components, tasks, resolveFn, args.KeepGoing)
	}
}
