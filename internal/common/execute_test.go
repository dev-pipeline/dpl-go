package common

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/dev-pipeline/dpl-go/internal/test/common"
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/resolve"
)

var (
	diamondProject dpl.Project = &testcommon.ResolveProject{
		Comps: testcommon.ResolveComponents{
			"foo": testcommon.ResolveComponent{},
			"bar": testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {"foo"},
				},
			},
			"baz": testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {"foo"},
				},
			},
			"biz": testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {
						"bar",
						"baz",
					},
				},
			},
		},
	}

	parallelProject dpl.Project = &testcommon.ResolveProject{
		Comps: testcommon.ResolveComponents{
			"foo": testcommon.ResolveComponent{},
			"bar": testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {"foo"},
				},
			},
			"baz": testcommon.ResolveComponent{},
			"biz": testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {
						"baz",
					},
				},
			},
		},
	}
)

func TestCleanRun(t *testing.T) {
	executeCount := 0
	resolveFn := resolve.GetResolver("deep")
	tasks := []Task{
		{
			Name: "build",
			Work: func(component dpl.Component) error {
				executeCount++
				return nil
			},
		},
	}

	err := runTasks(diamondProject, []string{"foo", "bar", "baz", "biz"}, tasks, resolveFn, false, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if executeCount != 4 {
		t.Fatalf("Executed wrong number of tasks (%v)", executeCount)
	}
}

func TestErrorRun(t *testing.T) {
	executeCount := 0
	resolveFn := resolve.GetResolver("deep")
	tasks := []Task{
		{
			Name: "build",
			Work: func(component dpl.Component) error {
				executeCount++
				return errors.New("Error")
			},
		},
	}

	err := runTasks(diamondProject, []string{"foo", "bar", "baz", "biz"}, tasks, resolveFn, false, 1)
	if err == nil {
		t.Fatalf("Missing expected error")
	}
	if executeCount != 1 {
		t.Fatalf("Executed too many tasks (%v)", executeCount)
	}
}

func TestRecoverError(t *testing.T) {
	executeCount := atomic.Int32{}
	resolveFn := resolve.GetResolver("deep")
	tasks := []Task{
		{
			Name: "build",
			Work: func(component dpl.Component) error {
				previous := executeCount.Add(1)
				if previous == 1 {
					// only fail the first one
					return errors.New("Error")
				}
				return nil
			},
		},
	}

	err := runTasks(parallelProject, []string{"foo", "bar", "baz", "biz"}, tasks, resolveFn, true, 4)
	if err == nil {
		t.Fatalf("Missing expected error")
	}
	if executeCount.Load() != 3 {
		t.Fatalf("Executed too many tasks (%v)", executeCount.Load())
	}
}
