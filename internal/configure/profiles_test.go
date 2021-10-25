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

func TestLoadSingleProfile(t *testing.T) {
	_, currentFilename, _, _ := runtime.Caller(0)
	configDir := path.Join(currentFilename, "..", "..", "..", "test_files", "configure")
	modSet, err := loadProfiles(configDir, []string{"foo"})
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
	}
	compareModSets(t, expectedModSet, modSet)
}
