package resolve

import (
	"testing"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

var (
	exists struct{}
)

type resolveComponent struct {
	data map[string]string
}

func (rs *resolveComponent) Name() string {
	return ""
}

func (rs *resolveComponent) GetValue(key string) (string, bool) {
	value, found := rs.data[key]
	return value, found
}

type resolveComponents map[string]resolveComponent

type resolveProject struct {
	components resolveComponents
}

func (rp *resolveProject) GetComponent(name string) (dpl.Component, bool) {
	component, found := rp.components[name]
	if found {
		return &component, true
	}
	return nil, false
}

func (rp *resolveProject) Components() []string {
	return nil
}

func compareDeps(t *testing.T, expected reverseDependencies, actual reverseDependencies) {
	if len(expected) != len(actual) {
		t.Fatalf("Unexpected size (%v vs %v)", len(expected), len(actual))
	}
	for expectedName, expectedDepSet := range expected {
		actualDepSet, found := actual[expectedName]
		if !found {
			t.Fatalf("Missing dependency info for '%v'", expectedName)
		}
		if len(actualDepSet) != len(expectedDepSet) {
			t.Fatalf("Unexpected size for %v reverse dependencies (%v vs %v)", expectedName, len(actualDepSet), len(expectedDepSet))
		}
		for expectedDep := range expectedDepSet {
			_, found = actualDepSet[expectedDep]
			if !found {
				t.Fatalf("Missing expected dependency (%v should reverse depend on %v)", expectedName, expectedDep)
			}
		}
	}
}

func TestSingleReverseDeps(t *testing.T) {
	targets := []string{"foo"}
	project := &resolveProject{
		components: resolveComponents{
			targets[0]: resolveComponent{},
		},
	}
	tasks := []string{"build"}

	revDeps, err := makeReverseDependencies(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDeps := reverseDependencies{
		"foo.build": depSet{},
	}
	compareDeps(t, expectedDeps, revDeps)
}

func TestMultipleTasks(t *testing.T) {
	targets := []string{"foo"}
	project := &resolveProject{
		components: resolveComponents{
			targets[0]: resolveComponent{},
		},
	}
	tasks := []string{"checkout", "build"}

	revDeps, err := makeReverseDependencies(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDeps := reverseDependencies{
		"foo.checkout": depSet{
			"foo.build": exists,
		},
		"foo.build": depSet{},
	}
	compareDeps(t, expectedDeps, revDeps)
}

func TestIndependentReverseDeps(t *testing.T) {
	targets := []string{"foo", "bar"}
	project := &resolveProject{
		components: resolveComponents{
			targets[0]: resolveComponent{},
			targets[1]: resolveComponent{},
		},
	}
	tasks := []string{"build"}

	revDeps, err := makeReverseDependencies(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDeps := reverseDependencies{
		"foo.build": depSet{},
		"bar.build": depSet{},
	}
	compareDeps(t, expectedDeps, revDeps)
}

func TestLinearReverseDeps(t *testing.T) {
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

	revDeps, err := makeReverseDependencies(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDeps := reverseDependencies{
		"foo.build": depSet{
			"bar.build": exists,
		},
		"bar.build": depSet{},
	}
	compareDeps(t, expectedDeps, revDeps)
}

func TestImplicitComponentTasks(t *testing.T) {
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
	tasks := []string{"checkout", "build"}

	revDeps, err := makeReverseDependencies(project, targets[1:], tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDeps := reverseDependencies{
		"foo.checkout": depSet{
			"foo.build": exists,
		},
		"foo.build": depSet{
			"bar.build": exists,
		},
		"bar.checkout": depSet{
			"bar.build": exists,
		},
		"bar.build": depSet{},
	}
	compareDeps(t, expectedDeps, revDeps)
}

func TestDiamondReverseDeps(t *testing.T) {
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

	revDeps, err := makeReverseDependencies(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDeps := reverseDependencies{
		"foo.build": depSet{
			"bar.build": exists,
			"baz.build": exists,
		},
		"bar.build": depSet{
			"biz.build": exists,
		},
		"baz.build": depSet{
			"biz.build": exists,
		},
		"biz.build": depSet{},
	}
	compareDeps(t, expectedDeps, revDeps)
}

func TestCircularReverseDeps(t *testing.T) {
	targets := []string{"foo"}
	project := &resolveProject{
		components: resolveComponents{
			targets[0]: resolveComponent{
				data: map[string]string{
					"depends.build": targets[0],
				},
			},
		},
	}
	tasks := []string{"build"}

	revDeps, err := makeReverseDependencies(project, targets, tasks)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDeps := reverseDependencies{
		"foo.build": depSet{
			"foo.build": exists,
		},
	}
	compareDeps(t, expectedDeps, revDeps)
}

func TestMissingComponentReverseDeps(t *testing.T) {
	targets := []string{"foo", "missing"}
	project := &resolveProject{
		components: resolveComponents{
			targets[0]: resolveComponent{
				data: map[string]string{
					"depends.build": targets[1],
				},
			},
		},
	}
	tasks := []string{"build"}

	_, err := makeReverseDependencies(project, targets[:1], tasks)
	if err == nil {
		t.Fatalf("Missing expected error")
	}
	missingError, success := err.(*ComponentNotFoundError)
	if !success {
		t.Fatalf("Unexpected error: %v", err)
	}
	if missingError.Name != targets[1] {
		t.Fatalf("Name mismatch (%v vs %v)", missingError.Name, targets[1])
	}
}
