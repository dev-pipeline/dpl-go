package configure

import (
	"fmt"
	"testing"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

func TestClean(t *testing.T) {
	badName := "hello"
	_, err := loadRawConfig([]byte(fmt.Sprintf("[%v]", badName)))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestRestrictedName(t *testing.T) {
	badName := "/hello"
	_, err := loadRawConfig([]byte(fmt.Sprintf("[%v]", badName)))

	if err == nil {
		t.Fatalf("Expected error")
	}

	realErr, ok := err.(*dpl.ComponentValidationError)
	if ok {
		origErr, ok := realErr.OriginalError.(*dpl.InvalidComponentNameError)
		if !ok {
			t.Fatalf("Unexpected error: %v", realErr.OriginalError)
		}
		if origErr.Name != badName {
			t.Fatalf("Invalid component name (got '%v', expected '%v')", origErr.Name, badName)
		}
	} else {
		t.Fatalf("Unexpected error: %v", err)
	}
}
