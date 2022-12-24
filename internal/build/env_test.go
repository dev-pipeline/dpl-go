package build

import (
	"testing"
)

func compareEnvironments(t *testing.T, actual []string, expected []string) {
	if len(actual) != len(expected) {
		t.Fatalf("Length mismatch: %v vs %v", len(actual), len(expected))
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("Mismatched environment at index %v: '%v' vs '%v'", i, actual[i], expected[i])
		}
	}
}

func TestPrependEmpty(t *testing.T) {
	env := []string{}
	name := "FOO"
	extra := []string{"a", "b", "c"}

	expected := []string{
		makeEnvString(name, extra),
	}
	newEnv := prependEnvironment(env, name, extra)
	compareEnvironments(t, newEnv, expected)
}

func TestPrependExisting(t *testing.T) {
	name := "FOO"
	starting := []string{"x", "y", "z"}
	env := []string{
		makeEnvString(name, starting),
	}
	extra := []string{"a", "b", "c"}

	expected := []string{
		makeEnvString(name, append(extra, starting...)),
	}
	newEnv := prependEnvironment(env, name, extra)
	compareEnvironments(t, newEnv, expected)
}

func TestAppendEmpty(t *testing.T) {
	env := []string{}
	name := "FOO"
	extra := []string{"a", "b", "c"}

	expected := []string{
		makeEnvString(name, extra),
	}
	newEnv := appendEnvironment(env, name, extra)
	compareEnvironments(t, newEnv, expected)
}

func TestAppendExisting(t *testing.T) {
	name := "FOO"
	starting := []string{"a", "b", "c"}
	env := []string{
		makeEnvString(name, starting),
	}
	extra := []string{"x", "y", "z"}

	expected := []string{
		makeEnvString(name, append(starting, extra...)),
	}
	newEnv := appendEnvironment(env, name, extra)
	compareEnvironments(t, newEnv, expected)
}
