package configure

import (
	"testing"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

func compareProject(t *testing.T, project dpl.Project, expected map[string]map[string][]string) {
	if len(project.ComponentNames()) != len(expected) {
		t.Fatalf("Mismatched lengths (%v vs %v)", len(project.ComponentNames()), len(expected))
	}
	for componentName, valueMap := range expected {
		component, err := project.GetComponent(componentName)
		if err != nil {
			t.Fatalf("Missing component %v (%v)", componentName, err)
		}
		for keyName, values := range valueMap {
			actualValues, err := component.ExpandValues(keyName)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if actualValues == nil {
				t.Fatalf("Missing values for %v", keyName)
			}
			if len(values) != len(actualValues) {
				t.Fatalf("Unequal number of values (%v vs %v)", len(values), len(actualValues))
			}
			for index := range values {
				if values[index] != actualValues[index] {
					t.Fatalf("Value mismatch (%v vs %v)", values[index], actualValues[index])
				}
			}
		}
	}
}

func TestLoadOverride(t *testing.T) {
	project, err := loadRawConfig([]byte(`
		[foo]
		x = b
		[bar]
		x = b
	`))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	buildConfigPath(func(dataDir string) {
		err := applyOverrides(dataDir, []string{
			"override-prepend",
		}, project)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expectedValues := map[string]map[string][]string{
			"foo": {
				"x": {
					"a",
					"b",
				},
			},
			"bar": {
				"x": {
					"b",
				},
			},
		}
		compareProject(t, project, expectedValues)
	}, "configure")
}
