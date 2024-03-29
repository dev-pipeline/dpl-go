package configure

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
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = world
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := newModifierSet()
	modifierSet.prependValues["a"] = []string{"hello"}

	foo, err := project.GetComponent("foo")
	if err != nil {
		t.Fatalf("Couldn't find component (%v)", err)
	}
	applyComponentModifiers(foo, modifierSet)

	expectedValues := []string{
		"hello",
		"world",
	}
	actualValues := foo.GetValues("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyPrependEmpty(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := newModifierSet()
	modifierSet.prependValues["a"] = []string{"hello"}

	foo, err := project.GetComponent("foo")
	if err != nil {
		t.Fatalf("Couldn't find component (%v)", err)
	}
	applyComponentModifiers(foo, modifierSet)

	expectedValues := []string{
		"hello",
	}
	actualValues := foo.GetValues("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyAppendEmpty(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := newModifierSet()
	modifierSet.appendValues["a"] = []string{"hello"}

	foo, err := project.GetComponent("foo")
	if err != nil {
		t.Fatalf("Couldn't find component (%v)", err)
	}
	applyComponentModifiers(foo, modifierSet)

	expectedValues := []string{
		"hello",
	}
	actualValues := foo.GetValues("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyAppend(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := newModifierSet()
	modifierSet.appendValues["a"] = []string{"world"}

	foo, err := project.GetComponent("foo")
	if err != nil {
		t.Fatalf("Couldn't find component (%v)", err)
	}
	applyComponentModifiers(foo, modifierSet)

	expectedValues := []string{
		"hello",
		"world",
	}
	actualValues := foo.GetValues("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyOverride(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := newModifierSet()
	modifierSet.overrideValues["a"] = []string{"goodbye"}

	foo, err := project.GetComponent("foo")
	if err != nil {
		t.Fatalf("Couldn't find component (%v)", err)
	}
	applyComponentModifiers(foo, modifierSet)

	expectedValues := []string{"goodbye"}
	actualValues := foo.GetValues("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyErase(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := newModifierSet()
	modifierSet.eraseValues["a"] = struct{}{}

	foo, err := project.GetComponent("foo")
	if err != nil {
		t.Fatalf("Couldn't find component (%v)", err)
	}
	applyComponentModifiers(foo, modifierSet)

	var expectedValues []string = nil
	actualValues := foo.GetValues("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyOverrideAfter(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := newModifierSet()
	modifierSet.prependValues["a"] = []string{"say"}
	modifierSet.appendValues["a"] = []string{"world"}
	modifierSet.overrideValues["a"] = []string{"goodbye"}

	foo, err := project.GetComponent("foo")
	if err != nil {
		t.Fatalf("Couldn't find component (%v)", err)
	}
	applyComponentModifiers(foo, modifierSet)

	expectedValues := []string{"goodbye"}
	actualValues := foo.GetValues("a")
	compareExpected(t, expectedValues, actualValues)
}

func TestApplyEraseLast(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	modifierSet := newModifierSet()
	modifierSet.prependValues["a"] = []string{"say"}
	modifierSet.appendValues["a"] = []string{"world"}
	modifierSet.overrideValues["a"] = []string{"goodbye"}
	modifierSet.eraseValues["a"] = struct{}{}

	foo, err := project.GetComponent("foo")
	if err != nil {
		t.Fatalf("Couldn't find component (%v)", err)
	}
	applyComponentModifiers(foo, modifierSet)

	var expectedValues []string = nil
	actualValues := foo.GetValues("a")
	compareExpected(t, expectedValues, actualValues)
}
