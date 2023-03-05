package configure

import (
	"fmt"
	"testing"
)

func TestParseSimple(t *testing.T) {
	componentNames := []string{
		"foo",
		"bar",
	}
	project, err := loadRawConfig(
		[]byte(fmt.Sprintf(`
			[%v]
			[%v]
		`, componentNames[0], componentNames[1])),
	)

	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(project.ComponentNames()) != 2 {
		t.Fatalf("Wrong number of components (expected %v)", 2)
	}
	for _, name := range componentNames {
		_, err := project.GetComponent(name)
		if err != nil {
			t.Fatalf("Missing component %v (%v)", name, err)
		}
	}
}

func TestParseMultiValue(t *testing.T) {
	project, err := loadRawConfig(
		[]byte(`
			[foo]
			build.depends = bar
			build.depends = baz
		`),
	)

	if err != nil {
		t.Fatalf("%v", err)
	}

	foo, _ := project.GetComponent("foo")
	expectedDepends := []string{
		"bar",
		"baz",
	}
	depends := foo.GetValues("build.depends")
	if len(expectedDepends) != len(depends) {
		t.Fatalf("Unexpected key counts (%v vs %v)", len(expectedDepends), len(depends))
	}

	for index := range expectedDepends {
		if expectedDepends[index] != depends[index] {
			t.Fatalf("Unexpected value at index %v (%v vs %v)", index, expectedDepends[index], depends[index])
		}
	}
}
