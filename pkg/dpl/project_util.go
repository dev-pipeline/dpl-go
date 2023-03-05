package dpl

import (
	"fmt"
)

type MissingKeyError struct {
	component  Component
	missingKey string
}

func (mke *MissingKeyError) Error() string {
	return fmt.Sprintf("%v missing value for key '%v'", mke.component.Name(), mke.missingKey)
}

type TooManyValuesError struct {
	component Component
	key       string
	values    []string
}

func (tmve *TooManyValuesError) Error() string {
	return fmt.Sprintf("%v has too many values for key '%v' (%v)", tmve.component.Name(), tmve.key, tmve.values)
}

func GetSingleComponentValue(component Component, key string) (string, error) {
	vals, err := component.ExpandValues(key)
	if err != nil {
		return "", err
	}
	if len(vals) == 0 {
		return "", &MissingKeyError{
			component:  component,
			missingKey: key,
		}
	}
	if len(vals) > 1 {
		return "", &TooManyValuesError{
			component: component,
			key:       key,
			values:    vals,
		}
	}
	return vals[0], nil
}

func GetSingleComponentValueOrDefault(component Component, key string, fallback string) (string, error) {
	val, err := GetSingleComponentValue(component, key)
	if err == nil {
		return val, nil
	}
	if _, ok := err.(*MissingKeyError); ok {
		return fallback, nil
	}
	return "", err
}
