package configure

import (
	"testing"
)

func TestNoExpand(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = bye
			b = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 1 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "bye" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
}

func TestSingleExpand(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${b}
			b = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 1 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "hello" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
}

func TestLeadingExpand(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = he${b}
			b = llo
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 1 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "hello" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
}

func TestTrailingExpand(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${b}llo
			b = he
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 1 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "hello" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
}

func TestCrossExpand(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${bar:x}
			[bar]
			x = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 1 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "hello" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
}

func TestCrossExpandSubKey(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${bar:x.y.z}
			[bar]
			x.y.z = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 1 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "hello" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
}

func TestCrossExpandLocal(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${bar:x}
			[bar]
			x = ${y}
			y = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 1 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "hello" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
}

func TestCrossKeyFailure(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${bar.x}
			[bar]
			y = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	_, err = foo.ExpandValues("a")
	if err == nil {
		t.Fatalf("Missing expected error")
	}
}

func TestCrossComponentFailure(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${baz.y}
			[bar]
			y = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	_, err = foo.ExpandValues("a")
	if err == nil {
		t.Fatalf("Missing expected error")
	}
}

func TestMultiExpand(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${b}${b}
			b = hello
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 1 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "hellohello" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
}

func TestRecursiveExpand(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${${b}}
			b = hello
			hello = bye
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 1 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "bye" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
}

func TestExpandEmpty(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${b}
			b = hello
			b = \e
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	expandedValues, err := foo.ExpandValues("a")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(expandedValues) != 2 {
		t.Fatalf("Unexpected number of expanded values: %v", len(expandedValues))
	}
	if expandedValues[0] != "hello" {
		t.Fatalf("Unexpected result: %v", expandedValues[0])
	}
	if expandedValues[1] != "" {
		t.Fatalf("Unexpected result: %v", expandedValues[1])
	}
}

func TestEscapeNothing(t *testing.T) {
	value, err := escapeValue(`hello`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value != `hello` {
		t.Fatalf("Unexpected result: %v", value)
	}
}

func TestEscapeDoubleSlashes(t *testing.T) {
	value, err := escapeValue(`\\\\`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value != `\\` {
		t.Fatalf("Unexpected result: %v", value)
	}
}

func TestEscapeLiteralPrefix(t *testing.T) {
	value, err := escapeValue(`\\hello`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value != `\hello` {
		t.Fatalf("Unexpected result: %v", value)
	}
}

func TestEscapeLiteralMiddle(t *testing.T) {
	value, err := escapeValue(`hel\\lo`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value != `hel\lo` {
		t.Fatalf("Unexpected result: %v", value)
	}
}

func TestEscapeEmptyPrefix(t *testing.T) {
	value, err := escapeValue(`\ehello`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value != "hello" {
		t.Fatalf("Unexpected result: %v", value)
	}
}

func TestEscapeEmptyMiddle(t *testing.T) {
	value, err := escapeValue(`hel\elo`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value != "hello" {
		t.Fatalf("Unexpected result: %v", value)
	}
}

func TestEscapeConsecutiveEmpty(t *testing.T) {
	value, err := escapeValue(`\e\e\e\ehello`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value != "hello" {
		t.Fatalf("Unexpected result: %v", value)
	}
}

func TestEscapeTrailingEmpty(t *testing.T) {
	value, err := escapeValue(`hello\e\e\e`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value != "hello" {
		t.Fatalf("Unexpected result: %v", value)
	}
}

func TestExpandLimit(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			a = ${a}
		`),
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	foo, _ := project.GetComponent("foo")
	_, err = foo.ExpandValues("a")
	if err == nil {
		t.Fatalf("Missing expected error")
	}
}
