package scm

import (
	"fmt"
	"log"
	"testing"

	"github.com/dev-pipeline/dpl-go/internal/test/common"
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

func buildTestUri(uri string) string {
	return fmt.Sprintf("test://%v", uri)
}

func buildErrorUri(uri string) string {
	return fmt.Sprintf("error://%v", uri)
}

func TestDoCheckout(t *testing.T) {
	uri := "some-uri"
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{
			scmUriKey: {buildTestUri(uri)},
		},
	}
	err := checkout(c)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestDoCheckoutError(t *testing.T) {
	uri := "some-uri"
	c := &testcommon.ResolveComponent{
		Data: map[string][]string{
			scmUriKey: {buildErrorUri(uri)},
		},
	}
	err := checkout(c)
	if err != errTestCheckout {
		t.Fatalf("Unexpected error: %v", err)
	}
}

type testScm struct {
}

func (testScm) Checkout(ScmInfo) error {
	return nil
}

func makeTestScm(dpl.Component) (ScmHandler, error) {
	return &testScm{}, nil
}

type errorScm struct {
}

var (
	errTestCheckout error = fmt.Errorf("some error")
)

func (errorScm) Checkout(ScmInfo) error {
	return errTestCheckout
}

func makeErrorScm(dpl.Component) (ScmHandler, error) {
	return &errorScm{}, nil
}

func init() {
	err := AddHandler("test", makeTestScm)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = AddHandler("error", makeErrorScm)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
