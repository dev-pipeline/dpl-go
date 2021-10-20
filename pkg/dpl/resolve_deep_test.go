package dpl

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
}
