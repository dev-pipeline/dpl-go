package dpl

import (
	"fmt"
	"testing"
)

func TestSingleExpandGood(t *testing.T) {
	key := "some-key"
	value := "some-value"
	component := &trivialComponent{
		ComponentName: "component",
		Data: map[string][]string{
			key: {value},
		},
	}

	actualValue, err := GetSingleComponentValue(component, key)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if actualValue != value {
		t.Fatalf("Unexpected value: %v (expected: %v)", actualValue, value)
	}
}

func TestSingleExpandMissing(t *testing.T) {
	key := "some-key"
	component := &trivialComponent{
		ComponentName: "component",
	}

	actualValue, err := GetSingleComponentValue(component, key)
	if _, ok := err.(*MissingKeyError); !ok {
		t.Fatalf("Unexpected error: %v", err)
	}
	if actualValue != "" {
		t.Fatalf("Unexpected value: %v", actualValue)
	}
}

func TestSingleExpandTooMany(t *testing.T) {
	key := "some-key"
	component := &trivialComponent{
		ComponentName: "component",
		Data: map[string][]string{
			key: {"some-value", "another-value"},
		},
	}

	actualValue, err := GetSingleComponentValue(component, key)
	if _, ok := err.(*TooManyValuesError); !ok {
		t.Fatalf("Unexpected error: %v", err)
	}
	if actualValue != "" {
		t.Fatalf("Unexpected value: %v", actualValue)
	}
}

func TestSingleExpandNoFallback(t *testing.T) {
	key := "some-key"
	value := "some-value"
	fallback := "another-value"
	component := &trivialComponent{
		ComponentName: "component",
		Data: map[string][]string{
			key: {value},
		},
	}

	actualValue, err := GetSingleComponentValueOrDefault(component, key, fallback)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if actualValue != value {
		t.Fatalf("Unexpected value: %v (expected: %v)", actualValue, value)
	}
}

func TestSingleExpandWithFallback(t *testing.T) {
	key := "some-key"
	fallback := "another-value"
	component := &trivialComponent{
		ComponentName: "component",
		Data:          map[string][]string{},
	}

	actualValue, err := GetSingleComponentValueOrDefault(component, key, fallback)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if actualValue != fallback {
		t.Fatalf("Unexpected value: %v (expected: %v)", actualValue, fallback)
	}
}

type erroringComponent struct {
	err error
}

func (*erroringComponent) Name() string {
	return ""
}

func (*erroringComponent) KeyNames() []string {
	return []string{}
}

func (*erroringComponent) GetValues(string) []string {
	return []string{}
}

func (ec *erroringComponent) ExpandValues(string) ([]string, error) {
	return nil, ec.err
}

func (*erroringComponent) SetValues(string, []string) {
}

func (*erroringComponent) EraseKey(string) {
}

func (*erroringComponent) GetSourceDir() string {
	return ""
}

func (*erroringComponent) GetWorkDir() string {
	return ""
}

func TestGetSingleComponentValueErr(t *testing.T) {
	component := &erroringComponent{
		err: fmt.Errorf("some error"),
	}

	_, err := GetSingleComponentValue(component, "some-key")
	if err != component.err {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestGetSingleComponentValueOrDefaultErr(t *testing.T) {
	component := &erroringComponent{
		err: fmt.Errorf("some error"),
	}

	_, err := GetSingleComponentValueOrDefault(component, "some-key", "some-fallback")
	if err != component.err {
		t.Fatalf("Unexpected error: %v", err)
	}
}
