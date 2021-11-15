package common

import (
	"errors"
	"testing"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/resolve"
	"github.com/dev-pipeline/dpl-go/test/common"
)

var (
	diamondProject dpl.Project = &testcommon.ResolveProject{
		Comps: testcommon.ResolveComponents{
			"foo": testcommon.ResolveComponent{},
			"bar": testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": []string{"foo"},
				},
			},
			"baz": testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": []string{"foo"},
				},
			},
			"biz": testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": []string{
						"bar",
						"baz",
					},
				},
			},
		},
	}
)

func TestCleanRun(t *testing.T) {
	resolveFn := resolve.GetResolver("deep")
	tasks := []Task{
		Task{
			Name: "build",
			Work: func(component dpl.Component) error {
				return nil
			},
		},
	}

	err := runTasks(diamondProject, []string{"foo", "bar", "baz", "biz"}, tasks, resolveFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestErrorRun(t *testing.T) {
	resolveFn := resolve.GetResolver("deep")
	tasks := []Task{
		Task{
			Name: "build",
			Work: func(component dpl.Component) error {
				return errors.New("Error")
			},
		},
	}

	err := runTasks(diamondProject, []string{"foo", "bar", "baz", "biz"}, tasks, resolveFn)
	if err == nil {
		t.Fatalf("Missing expected error")
	}
}
