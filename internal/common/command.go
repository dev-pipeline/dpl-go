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
	MaxTasks     int
}

type TaskFn func(dpl.Component) error

type Task struct {
	Name string
	Work TaskFn
}

type work struct {
	fn        TaskFn
	name      string
	component dpl.Component
}

type taskComplete struct {
	name string
	err  error
}

func executeTasks(taskChannel chan work, doneChannel chan taskComplete) {
	for {
		workUnit, ok := <-taskChannel
		if !ok {
			return
		}
		log.Printf("Executing %v", workUnit.name)
		err := workUnit.fn(workUnit.component)
		doneChannel <- taskComplete{
			name: workUnit.name,
			err:  err,
		}
	}
}

func makeTaskContainers(tasks []Task) ([]string, map[string]TaskFn) {
	taskList := []string{}
	taskMap := map[string]TaskFn{}

	for _, task := range tasks {
		taskMap[task.Name] = task.Work
		taskList = append(taskList, task.Name)
	}

	return taskList, taskMap
}

type taskCompleteFn func(taskComplete)

func startDrainComplete(doneChannel chan taskComplete, completeFn taskCompleteFn) {
	go func() {
		for {
			completedTask, ok := <-doneChannel
			if !ok {
				return
			}
			completeFn(completedTask)
		}
	}()
}

func startResolve(project dpl.Project, resolver resolve.Resolver, taskMap map[string]TaskFn) chan work {
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
	return workChannel
}

type failedTask struct {
	originalError error
	name          string
	dependents    []string
}

func (ft *failedTask) Error() string {
	return ""
}

func runTasks(project dpl.Project, components []string, tasks []Task, resolveFn resolve.ResolveFn, keepGoing bool, maxTasks int) error {
	taskList, taskMap := makeTaskContainers(tasks)
	resolver, err := resolveFn(project, components, taskList)
	if err != nil {
		return err
	}

	workChannel := startResolve(project, resolver, taskMap)
	doneChannel := make(chan taskComplete)
	defer close(doneChannel)
	wg := sync.WaitGroup{}

	for i := 0; i < maxTasks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			executeTasks(workChannel, doneChannel)
		}()
	}

	errors := []error{}
	m := sync.Mutex{}
	completeTask := func(completedTask taskComplete) {
		if completedTask.err == nil {
			resolver.Complete(completedTask.name)
		} else {
			log.Printf("Error executing task '%v': %v", completedTask.name, completedTask.err)
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
	log.Printf("%v total error(s)", len(errors))
	return errors[0]
}

func DoCommand(components []string, args Args, tasks []Task) {
	project, err := dpl.LoadProject()
	if err != nil {
		log.Fatalf("Failed to load project: %v", err)
	}

	if len(components) == 0 {
		components = project.Components()
	}
	resolveFn := resolve.GetResolver(args.Dependencies)
	if resolveFn == nil {
		log.Fatalf("No resolver '%v'", args.Dependencies)
	} else {
		err := runTasks(project, components, tasks, resolveFn, args.KeepGoing, args.MaxTasks)
		project.Write()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}
