package configure

import (
	"path"
	"runtime"
	"testing"

	"github.com/dev-pipeline/dpl-go/pkg/dpl/configfile"
)

func compareStringMaps(t *testing.T, expected map[string][]string, actual map[string][]string) {
	if len(expected) != len(actual) {
		t.Fatalf("Length mismatch (%v vs %v)", len(expected), len(actual))
	}
	for key, expectedValue := range expected {
		actualValue, found := actual[key]
		if !found {
			t.Fatalf("Couldn't find %v", key)
		}
		if len(expectedValue) != len(actualValue) {
			t.Fatalf("Length mismatch (%v vs %v)", len(expectedValue), len(actualValue))
		}
		for index := range expectedValue {
			if expectedValue[index] != actualValue[index] {
				t.Fatalf("Value mismatch (%v vs %v)", expectedValue[index], actualValue[index])
			}
		}
	}
}

func compareEmptyMaps(t *testing.T, expected map[string]struct{}, actual map[string]struct{}) {
}

func compareModSets(t *testing.T, expectedModSet configfile.ModifierSet, actualModSet configfile.ModifierSet) {
	compareStringMaps(t, expectedModSet.PrependValues, actualModSet.PrependValues)
	compareStringMaps(t, expectedModSet.AppendValues, actualModSet.AppendValues)
	compareStringMaps(t, expectedModSet.OverrideValues, actualModSet.OverrideValues)
	compareEmptyMaps(t, expectedModSet.EraseValues, actualModSet.EraseValues)
}

func buildConfigPath(testFunc func(string), chunks ...string) {
	_, currentFilename, _, _ := runtime.Caller(0)
	fullArgs := append([]string{
		currentFilename,
		"..",
		"..",
		"..",
		"test_files",
	}, chunks...)
	dataDir := path.Join(fullArgs...)
	testFunc(dataDir)
}

func TestLoadMultipleProfiles(t *testing.T) {
	buildConfigPath(func(dataDir string) {
		modSet, err := loadProfiles(dataDir, []string{
			"prepend-xab",
			"append-xyz",
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expectedModSet := configfile.ModifierSet{
			PrependValues: map[string][]string{
				"x": []string{
					"a",
					"b",
				},
			},
			AppendValues: map[string][]string{
				"x": []string{
					"y",
					"z",
				},
			},
		}
		compareModSets(t, expectedModSet, modSet)
	}, "configure")
}
