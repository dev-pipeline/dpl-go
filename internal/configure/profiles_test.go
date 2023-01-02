package configure

import (
	"path"
	"runtime"
	"testing"
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

func compareModSets(t *testing.T, expectedModSet modifierSet, actualModSet modifierSet) {
	compareStringMaps(t, expectedModSet.prependValues, actualModSet.prependValues)
	compareStringMaps(t, expectedModSet.appendValues, actualModSet.appendValues)
	compareStringMaps(t, expectedModSet.overrideValues, actualModSet.overrideValues)
	compareEmptyMaps(t, expectedModSet.eraseValues, actualModSet.eraseValues)
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

		expectedModSet := modifierSet{
			prependValues: map[string][]string{
				"x": {
					"a",
					"b",
				},
			},
			appendValues: map[string][]string{
				"x": {
					"y",
					"z",
				},
			},
		}
		compareModSets(t, expectedModSet, modSet)
	}, "configure")
}
