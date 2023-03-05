package dpl

import (
	"fmt"
	"testing"
)

var (
	errComponentError error = fmt.Errorf("error finding Component")
)

const (
	testValidatorName string = "some-test-validator"
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
		ComponentName: "component",
		Data: map[string][]string{
			badKey: {},
		},
	}

	err := ValidateComponent(&badField)
	if err == nil {
		t.Fatalf("Missing expected error")
	}

	realErr, ok := err.(*ComponentValidationError)
	if ok {
		origErr, ok := realErr.OriginalError.(*InvalidFieldNameError)
		if !ok {
			t.Fatalf("Unexpected error: %v", realErr.OriginalError)
		}
		if origErr.Name != badKey {
			t.Fatalf("Invalid field name (got '%v', expected '%v')", origErr.Name, badKey)
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

type invalidProject struct {
}

func (invalidProject) GetComponent(string) (Component, error) {
	return nil, errComponentError
}

func (invalidProject) ComponentNames() []string {
	return []string{}
}

func (invalidProject) Write() error {
	return nil
}

func TestInvalidProject(t *testing.T) {
	p := &invalidProject{}
	err := ValidateProject(p)
	if err == nil {
		t.Fatalf("Missing expected error")
	}
	realErr, ok := err.(*projectValidationError)
	if !ok {
		t.Fatalf("Unexpected error: %v", err)
	}
	if realErr.name != testValidatorName {
		t.Fatalf("Unexpected validator name: %v", realErr.name)
	}
}

func init() {
	RegisterProjectValidator(testValidatorName, func(p Project) error {
		_, ok := p.(*invalidProject)
		if ok {
			return fmt.Errorf("some error")
		}
		return nil
	})
}
