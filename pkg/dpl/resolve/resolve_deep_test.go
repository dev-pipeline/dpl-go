package resolve

import (
	"testing"
)

func compareCounts(t *testing.T, expectedCounts map[string]int, actualCounts map[string]int) {
	if len(expectedCounts) != len(actualCounts) {
		t.Fatalf("Mismatched count sizes (%v vs %v)", len(expectedCounts), len(actualCounts))
	}
	for task, expectedCount := range expectedCounts {
		actualCount, found := actualCounts[task]
		if !found {
			t.Fatalf("Missing count for %v", task)
		}
		if expectedCount != actualCount {
			t.Fatalf("Mismatched counts for task %v (%v vs %v)", task, expectedCount, actualCount)
		}
	}
}

func compareReadyMaps(t *testing.T, setName string, first []string, second []string) {
	lookupMap := map[string]struct{}{}
	for _, task := range first {
		lookupMap[task] = struct{}{}
	}

	for _, task := range second {
		_, found := lookupMap[task]
		if !found {
			t.Fatalf("Missing task %v from %v set", task, setName)
		}
	}
}

func compareReady(t *testing.T, expectedReady []string, actualReady []string) {
	if len(expectedReady) != len(actualReady) {
		t.Fatalf("Mismatched count sizes (%v vs %v)", len(expectedReady), len(actualReady))
	}

	compareReadyMaps(t, "expected", expectedReady, actualReady)
	compareReadyMaps(t, "actual", actualReady, expectedReady)
}

func TestSingleComponent(t *testing.T) {
	targets := []string{"foo"}
	project := &resolveProject{
		components: resolveComponents{
			targets[0]: resolveComponent{},
		},
	}
	tasks := []string{"build"}

	resolver, err := ResolveDeep(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedCounts := map[string]int{}
	compareCounts(t, expectedCounts, resolver.depCounts)

	expectedReady := []string{"foo.build"}
	compareReady(t, expectedReady, resolver.readyTasks)

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

func TestSimpleDeps(t *testing.T) {
	targets := []string{"foo", "bar"}
	project := &resolveProject{
		components: resolveComponents{
			targets[0]: resolveComponent{},
			targets[1]: resolveComponent{
				data: map[string]string{
					"depends.build": "foo",
				},
			},
		},
	}
	tasks := []string{"build"}

	resolver, err := ResolveDeep(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedCounts := map[string]int{
		"bar.build": 1,
	}
	compareCounts(t, expectedCounts, resolver.depCounts)

	expectedReady := []string{"foo.build"}
	compareReady(t, expectedReady, resolver.readyTasks)

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

func TestDiamondDeps(t *testing.T) {
	targets := []string{"foo", "bar", "baz", "biz"}
	project := &resolveProject{
		components: resolveComponents{
			targets[0]: resolveComponent{},
			targets[1]: resolveComponent{
				data: map[string]string{
					"depends.build": "foo",
				},
			},
			targets[2]: resolveComponent{
				data: map[string]string{
					"depends.build": "foo",
				},
			},
			targets[3]: resolveComponent{
				data: map[string]string{
					"depends.build": "bar,baz",
				},
			},
		},
	}
	tasks := []string{"build"}

	resolver, err := ResolveDeep(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedCounts := map[string]int{
		"bar.build": 1,
		"baz.build": 1,
		"biz.build": 2,
	}
	compareCounts(t, expectedCounts, resolver.depCounts)

	expectedReady := []string{"foo.build"}
	compareReady(t, expectedReady, resolver.readyTasks)

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
	if len(ready) != 2 {
		t.Fatalf("Unexpected ready length (expected 2, got %v)", len(ready))
	}
	readyTasks := map[string]struct{}{}
	for _, task := range ready {
		readyTasks[task] = struct{}{}
	}
	for _, expectedTask := range []string{"bar.build", "baz.build"} {
		_, found := readyTasks[expectedTask]
		if !found {
			t.Fatalf("Missing expected task: %v", expectedTask)
		}
		resolver.Complete(expectedTask)
	}

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

func TestFailDiamondDeps(t *testing.T) {
	targets := []string{"foo", "bar", "baz", "biz"}
	project := &resolveProject{
		components: resolveComponents{
			targets[0]: resolveComponent{},
			targets[1]: resolveComponent{
				data: map[string]string{
					"depends.build": "foo",
				},
			},
			targets[2]: resolveComponent{
				data: map[string]string{
					"depends.build": "foo",
				},
			},
			targets[3]: resolveComponent{
				data: map[string]string{
					"depends.build": "bar,baz",
				},
			},
		},
	}
	tasks := []string{"build"}

	resolver, err := ResolveDeep(project, targets, tasks)
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
