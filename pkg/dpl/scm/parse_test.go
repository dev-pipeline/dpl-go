package scm

import (
	"fmt"
	"strings"
	"testing"
)

func buildUri(scheme string, uri string, arguments map[string]string) string {
	parts := []string{fmt.Sprintf("%v://%v", scheme, uri)}
	for key, value := range arguments {
		parts = append(parts, fmt.Sprintf("%v=%v", key, value))
	}
	return strings.Join(parts, ";")
}

func compareArgs(t *testing.T, expected map[string]string, actual map[string]string) {
	if len(expected) != len(actual) {
		t.Fatalf("Mismatched sizes (%v vs %v)", len(expected), len(actual))
	}
	for key, expectedValue := range expected {
		actualValue, found := actual[key]
		if !found {
			t.Fatalf("Missing argument for key %v", key)
		}
		if expectedValue != actualValue {
			t.Fatalf("Mismatch for key %v (%v vs %v)", key, expectedValue, actualValue)
		}
	}
}

func TestParse(t *testing.T) {
	scheme := "git"
	scmUri := "git@github.com/foo/bar.git"
	expectedArguments := map[string]string{
		"protocol": "ssh",
	}
	fullUri := buildUri(scheme, scmUri, expectedArguments)

	scmInfo, err := BuildScmInfo(fullUri)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if scmInfo.Scheme != scheme {
		t.Fatalf("Unexpected scheme (%v vs %v)", scmInfo.Scheme, scheme)
	}
	if scmInfo.Path != scmUri {
		t.Fatalf("Unexpected uri (%v vs %v)", scmInfo.Path, scmUri)
	}
	compareArgs(t, expectedArguments, scmInfo.Arguments)
}

func TestBadParse(t *testing.T) {
	scmUri := "git://git@github.com/foo/bar.git;some-junk-without-a-value"

	_, err := BuildScmInfo(scmUri)
	if err == nil {
		t.Fatalf("Missing expected error")
	}
	if _, ok := err.(*extractionError); !ok {
		t.Fatalf("Unexpected error type: %T", err)
	}
}
