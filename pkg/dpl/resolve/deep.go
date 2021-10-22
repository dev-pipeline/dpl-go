package resolve

import (
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type deepResolver struct {
	commonResolver
}

func resolveDeep(project dpl.Project, targets []string, tasks []string) (*deepResolver, error) {
	revDeps, err := makeReverseDependencies(project, targets, tasks)
	if err != nil {
		return nil, err
	}
	common, err := resolveCommon(revDeps)
	if err != nil {
		return nil, err
	}
	return &deepResolver{
		commonResolver: common,
	}, nil
}

func init() {
	RegisterResolver("deep", func(project dpl.Project, targets []string, tasks []string) (Resolver, error) {
		return resolveDeep(project, targets, tasks)
	})
}
