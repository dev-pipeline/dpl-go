package common

import (
	"testing"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/resolve"
	"github.com/dev-pipeline/dpl-go/test/common"
)

func TestCleanRun(t *testing.T) {
	project := &testcommon.ResolveProject{
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

	resolveFn := resolve.GetResolver("deep")
	tasks := []Task{
		Task{
			Name: "build",
			Work: func(component dpl.Component) error {
				// t.Logf("%v", component.Name())
				return nil
			},
		},
	}

	err := runTasks(project, []string{"foo", "bar", "baz", "biz"}, tasks, resolveFn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
