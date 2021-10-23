package configfile

import (
	"testing"
)

func TestNoExpand(t *testing.T) {
	project, err := LoadRawConfig(
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
	expandedValues, err := foo.ExpandValue("a")
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
	project, err := LoadRawConfig(
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
	expandedValues, err := foo.ExpandValue("a")
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

func TestMultiExpand(t *testing.T) {
	project, err := LoadRawConfig(
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
	expandedValues, err := foo.ExpandValue("a")
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
	project, err := LoadRawConfig(
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
	expandedValues, err := foo.ExpandValue("a")
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
