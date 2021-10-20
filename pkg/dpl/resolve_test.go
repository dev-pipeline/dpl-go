package dpl

import (
	"testing"
)

var (
	exists struct{}
)

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
	project := &Project{
		ComponentInfo: Components{
			targets[0]: &Component{
				Name: targets[0],
			},
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
	project := &Project{
		ComponentInfo: Components{
			targets[0]: &Component{
				Name: targets[0],
			},
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
	project := &Project{
		ComponentInfo: Components{
			targets[0]: &Component{
				Name: targets[0],
			},
			targets[1]: &Component{
				Name: targets[1],
			},
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
	project := &Project{
		ComponentInfo: Components{
			targets[0]: &Component{
				Name: targets[0],
			},
			targets[1]: &Component{
				Name: targets[1],
				Data: map[string]string{
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
	project := &Project{
		ComponentInfo: Components{
			targets[0]: &Component{
				Name: targets[0],
			},
			targets[1]: &Component{
				Name: targets[1],
				Data: map[string]string{
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
	project := &Project{
		ComponentInfo: Components{
			targets[0]: &Component{
				Name: targets[0],
			},
			targets[1]: &Component{
				Name: targets[1],
				Data: map[string]string{
					"depends.build": "foo",
				},
			},
			targets[2]: &Component{
				Name: targets[2],
				Data: map[string]string{
					"depends.build": "foo",
				},
			},
			targets[3]: &Component{
				Name: targets[3],
				Data: map[string]string{
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
