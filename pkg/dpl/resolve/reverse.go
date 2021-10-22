package resolve

import (
	"strings"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type reverseResolver struct {
	commonResolver
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

func resolveReverse(project dpl.Project, targets []string, tasks []string) (*reverseResolver, error) {
	revDeps, err := makeReverseDependencies(project, project.ComponentNames(), tasks)
	if err != nil {
		return nil, err
	}
	trimmedDeps := trimReverseDependencies(revDeps, targets, tasks)

	common, err := resolveCommon(trimmedDeps)
	if err != nil {
		return nil, err
	}
	return &reverseResolver{
		commonResolver: common,
	}, nil
}

func init() {
	RegisterResolver("reverse", func(project dpl.Project, targets []string, tasks []string) (Resolver, error) {
		return resolveReverse(project, targets, tasks)
	})
}
