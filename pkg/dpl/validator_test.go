package dpl

import (
	"testing"
)

func TestRestrictedComponentName(t *testing.T) {
	badName := trivialComponent{
		ComponentName: "/hello",
	}

	err := validateComponentName(&badName)
	if err == nil {
		t.Fatalf("Missing expected error")
	}

	origErr, ok := err.(*InvalidComponentNameError)
	if ok {
		if origErr.Name != badName.ComponentName {
			t.Fatalf("Invalid component name (got '%v', expected '%v')", origErr.Name, badName.ComponentName)
		}
	} else {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func testBadFieldHelper(t *testing.T, badKey string) {
	badField := trivialComponent{
		Data: map[string][]string{
			badKey: []string{},
		},
	}

	err := validateFieldName(&badField)
	if err == nil {
		t.Fatalf("Missing expected error")
	}

	origErr, ok := err.(*InvalidFieldNameError)
	if ok {
		if origErr.Name != badKey {
			t.Fatalf("Invalid component name (got '%v', expected '%v')", origErr.Name, badKey)
		}
	} else {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestFieldStartsWithNumber(t *testing.T) {
	testBadFieldHelper(t, "123")
}

func TestDplPrefix(t *testing.T) {
	testBadFieldHelper(t, "dpl.a")
}
