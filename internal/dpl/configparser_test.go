package dplint

import (
	"fmt"
	"testing"
)

func TestParseSimple(t *testing.T) {
	componentNames := []string{
		"foo",
		"bar",
	}
	project, err := LoadRawConfig(
		[]byte(fmt.Sprintf(`
			[%v]
			[%v]
		`, componentNames[0], componentNames[1])),
	)

	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(project.Components) != 2 {
		t.Fatalf("Wrong number of components (expected %v)", 2)
	}
	for _, name := range componentNames {
		_, found := project.Components[name]
		if !found {
			t.Fatalf("Missing component %v", name)
		}
	}
}
