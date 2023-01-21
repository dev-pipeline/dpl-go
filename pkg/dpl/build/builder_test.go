package build

import (
	"fmt"
	"log"
	"testing"

	"github.com/dev-pipeline/dpl-go/internal/test/common"
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

var (
	errCantMakeBuilder error = fmt.Errorf("can't make builder")
)

const (
	dummyBuilderName string = "dummy"
	errorBuilderName string = "error"
)

type dummyBuilder struct {
}

func (dummyBuilder) Configure(*BuildConfig) error {
	return nil
}

func (dummyBuilder) Build(*BuildConfig) error {
	return nil
}

func (dummyBuilder) Install(string) error {
	return nil
}

func TestMakeBuilderError(t *testing.T) {
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{
			buildToolKey: {errorBuilderName},
		},
	}
	err := doFullBuild(c)
	if err != errCantMakeBuilder {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestRegisterAgain(t *testing.T) {
	err := RegisterBuilder(dummyBuilderName, makeDummyBuilder)
	if err != errAlreadyRegistered {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestMissingBuilder(t *testing.T) {
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{
			buildToolKey: {"none2"},
		},
	}
	err := doFullBuild(c)
	if err != errInvalidBuilder {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func makeDummyBuilder(dpl.Component) (Builder, error) {
	return &dummyBuilder{}, nil
}

func makeErrorBuilder(dpl.Component) (Builder, error) {
	return nil, errCantMakeBuilder
}

func init() {
	err := RegisterBuilder(dummyBuilderName, makeDummyBuilder)
	if err != nil {
		log.Fatalf("Error registring builder: %v", err)
	}
	err = RegisterBuilder(errorBuilderName, makeErrorBuilder)
	if err != nil {
		log.Fatalf("Error registring builder: %v", err)
	}
}
