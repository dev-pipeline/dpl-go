package resolve

import (
	"testing"

	"github.com/dev-pipeline/dpl-go/test/common"
)

func TestSingleComponentReverse(t *testing.T) {
	targets := []string{"foo"}
	project := &testcommon.ResolveProject{
		Comps: testcommon.ResolveComponents{
			targets[0]: testcommon.ResolveComponent{},
		},
	}
	tasks := []string{"build"}

	resolver, err := resolveReverse(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	taskChannel := make(chan []string)
	resolver.Resolve(taskChannel)

	ready := <-taskChannel
	if len(ready) != 1 {
		t.Fatalf("Unexpected ready length (expected 1, got %v)", len(ready))
	}
	if ready[0] != "foo.build" {
		t.Fatalf("Unexpected ready target (%v)", ready[0])
	}
	resolver.Complete(ready[0])
	ready = <-taskChannel
	if len(ready) != 0 {
		t.Fatalf("Unexpected ready result (%v)", ready)
	}
}

func TestSimpleDepsReverse(t *testing.T) {
	targets := []string{"foo", "bar"}
	project := &testcommon.ResolveProject{
		Comps: testcommon.ResolveComponents{
			targets[0]: testcommon.ResolveComponent{},
			targets[1]: testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {"foo"},
				},
			},
		},
	}
	tasks := []string{"build"}

	resolver, err := resolveReverse(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	taskChannel := make(chan []string)
	resolver.Resolve(taskChannel)

	ready := <-taskChannel
	if len(ready) != 1 {
		t.Fatalf("Unexpected ready length (expected 1, got %v)", len(ready))
	}
	if ready[0] != "foo.build" {
		t.Fatalf("Unexpected ready target (%v)", ready[0])
	}
	resolver.Complete(ready[0])
	ready = <-taskChannel
	if len(ready) != 1 {
		t.Fatalf("Unexpected ready length (expected 1, got %v)", len(ready))
	}
	if ready[0] != "bar.build" {
		t.Fatalf("Unexpected ready target (%v)", ready[0])
	}
	resolver.Complete(ready[0])
	ready = <-taskChannel
	if len(ready) != 0 {
		t.Fatalf("Unexpected ready result (%v)", ready)
	}
}

func TestDiamondDepsReverse(t *testing.T) {
	targets := []string{"foo", "bar", "baz", "biz"}
	project := &testcommon.ResolveProject{
		Comps: testcommon.ResolveComponents{
			targets[0]: testcommon.ResolveComponent{},
			targets[1]: testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {"foo"},
				},
			},
			targets[2]: testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {"foo"},
				},
			},
			targets[3]: testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {
						"bar",
						"baz",
					},
				},
			},
		},
	}
	tasks := []string{"build"}

	resolver, err := resolveReverse(project, targets[2:3], tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	taskChannel := make(chan []string)
	resolver.Resolve(taskChannel)

	ready := <-taskChannel
	if len(ready) != 1 {
		t.Fatalf("Unexpected ready length (expected 1, got %v)", len(ready))
	}
	if ready[0] != "baz.build" {
		t.Fatalf("Unexpected ready target (%v)", ready[0])
	}
	resolver.Complete(ready[0])

	ready = <-taskChannel
	if len(ready) != 1 {
		t.Fatalf("Unexpected ready length (expected 1, got %v)", len(ready))
	}
	if ready[0] != "biz.build" {
		t.Fatalf("Unexpected ready target (%v)", ready[0])
	}
	resolver.Complete(ready[0])

	ready = <-taskChannel
	if len(ready) != 0 {
		t.Fatalf("Unexpected ready result (%v)", ready)
	}
}

func TestFailDiamondDepsReverse(t *testing.T) {
	targets := []string{"foo", "bar", "baz", "biz"}
	project := &testcommon.ResolveProject{
		Comps: testcommon.ResolveComponents{
			targets[0]: testcommon.ResolveComponent{},
			targets[1]: testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {"foo"},
				},
			},
			targets[2]: testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {"foo"},
				},
			},
			targets[3]: testcommon.ResolveComponent{
				Data: map[string][]string{
					"depends.build": {
						"bar",
						"baz",
					},
				},
			},
		},
	}
	tasks := []string{"build"}

	resolver, err := resolveReverse(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	taskChannel := make(chan []string)
	resolver.Resolve(taskChannel)

	ready := <-taskChannel
	if len(ready) != 1 {
		t.Fatalf("Unexpected ready length (expected 1, got %v)", len(ready))
	}
	if ready[0] != "foo.build" {
		t.Fatalf("Unexpected ready target (%v)", ready[0])
	}
	failedTasks := resolver.Fail(ready[0])

	expectedFailures := map[string]struct{}{
		"bar.build": exists,
		"baz.build": exists,
		"biz.build": exists,
	}

	if len(failedTasks) != len(expectedFailures) {
		t.Fatalf("Unexpected number of failed jobs (expected %v, got %v)", len(expectedFailures), len(failedTasks))
	}
	for _, failedTask := range failedTasks {
		_, found := expectedFailures[failedTask]
		if !found {
			t.Fatalf("Unexpected failure (%v)", failedTask)
		}
	}

	ready = <-taskChannel
	if len(ready) != 0 {
		t.Fatalf("Unexpected ready result (%v)", ready)
	}
}
