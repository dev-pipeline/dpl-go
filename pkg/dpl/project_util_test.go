package dpl

import (
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
