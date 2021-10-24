package configfile

import (
	"testing"
)

func compareExpected(t *testing.T, expected []string, actual []string) {
	if len(expected) != len(actual) {
		t.Fatalf("Unexpected result sizes (%v vs %v)", expected, actual)
	}
	for index := range expected {
		if expected[index] != actual[index] {
			t.Fatalf("Mismatched result: (%v vs %v)", expected[index], actual[index])
		}
	}
}

func TestApplyPrepend(t *testing.T) {
	project, err := LoadRawConfig(
		[]byte(`
			[foo]
			a = world
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := NewModifierSet()
	modifierSet.prependValues["a"] = []string{"hello"}

	foo, found := project.GetComponent("foo")
	if !found {
		t.Fatalf("Couldn't find component")
	}
	ApplyComponentModifiers(foo, modifierSet)

	expectedValues := []string{
		"hello",
		"world",
	}
	actualValues := foo.GetValue("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyPrependEmpty(t *testing.T) {
	project, err := LoadRawConfig(
		[]byte(`
			[foo]
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := NewModifierSet()
	modifierSet.prependValues["a"] = []string{"hello"}

	foo, found := project.GetComponent("foo")
	if !found {
		t.Fatalf("Couldn't find component")
	}
	ApplyComponentModifiers(foo, modifierSet)

	expectedValues := []string{
		"hello",
	}
	actualValues := foo.GetValue("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyAppendEmpty(t *testing.T) {
	project, err := LoadRawConfig(
		[]byte(`
			[foo]
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := NewModifierSet()
	modifierSet.appendValues["a"] = []string{"hello"}

	foo, found := project.GetComponent("foo")
	if !found {
		t.Fatalf("Couldn't find component")
	}
	ApplyComponentModifiers(foo, modifierSet)

	expectedValues := []string{
		"hello",
	}
	actualValues := foo.GetValue("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyAppend(t *testing.T) {
	project, err := LoadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := NewModifierSet()
	modifierSet.appendValues["a"] = []string{"world"}

	foo, found := project.GetComponent("foo")
	if !found {
		t.Fatalf("Couldn't find component")
	}
	ApplyComponentModifiers(foo, modifierSet)

	expectedValues := []string{
		"hello",
		"world",
	}
	actualValues := foo.GetValue("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyOverride(t *testing.T) {
	project, err := LoadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := NewModifierSet()
	modifierSet.overrideValues["a"] = []string{"goodbye"}

	foo, found := project.GetComponent("foo")
	if !found {
		t.Fatalf("Couldn't find component")
	}
	ApplyComponentModifiers(foo, modifierSet)

	expectedValues := []string{"goodbye"}
	actualValues := foo.GetValue("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyErase(t *testing.T) {
	project, err := LoadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := NewModifierSet()
	modifierSet.eraseValues["a"] = struct{}{}

	foo, found := project.GetComponent("foo")
	if !found {
		t.Fatalf("Couldn't find component")
	}
	ApplyComponentModifiers(foo, modifierSet)

	var expectedValues []string = nil
	actualValues := foo.GetValue("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyOverrideAfter(t *testing.T) {
	project, err := LoadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := NewModifierSet()
	modifierSet.prependValues["a"] = []string{"say"}
	modifierSet.appendValues["a"] = []string{"world"}
	modifierSet.overrideValues["a"] = []string{"goodbye"}

	foo, found := project.GetComponent("foo")
	if !found {
		t.Fatalf("Couldn't find component")
	}
	ApplyComponentModifiers(foo, modifierSet)

	expectedValues := []string{"goodbye"}
	actualValues := foo.GetValue("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyEraseLast(t *testing.T) {
	project, err := LoadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := NewModifierSet()
	modifierSet.prependValues["a"] = []string{"say"}
	modifierSet.appendValues["a"] = []string{"world"}
	modifierSet.overrideValues["a"] = []string{"goodbye"}
	modifierSet.eraseValues["a"] = struct{}{}

	foo, found := project.GetComponent("foo")
	if !found {
		t.Fatalf("Couldn't find component")
	}
	ApplyComponentModifiers(foo, modifierSet)

	var expectedValues []string = nil
	actualValues := foo.GetValue("a")
	compareExpected(t, expectedValues, actualValues)
}
