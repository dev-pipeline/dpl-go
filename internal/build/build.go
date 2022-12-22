package build

import (
	"fmt"
	"log"

	"github.com/dev-pipeline/dpl-go/internal/common"
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

const (
	buildToolKey string = "build.tool"
)

var (
	BuildTask common.Task = common.Task{
		Name: "build",
		Work: doBuild,
	}

	builders map[string]MakeBuilder = map[string]MakeBuilder{}

	errAlreadyRegistered error = fmt.Errorf("builder already registered")
	errTooManyBuilders   error = fmt.Errorf("multiple builders")
	errInvalidBuilder    error = fmt.Errorf("invalid builder")
)

type Builder interface {
	Configure() error
	Build() error
	Install() error
}

type MakeBuilder func(dpl.Component) (Builder, error)

func doBuild(component dpl.Component) error {
	builder := component.GetValue(buildToolKey)
	if len(builder) != 1 {
		if len(builder) == 0 {
			log.Printf("No builder specified for '%v'", component.Name())
			return nil
		}
		return errTooManyBuilders
	}
	builderMaker, found := builders[builder[0]]
	if !found {
		return errInvalidBuilder
	}

	actualBuilder, err := builderMaker(component)
	if err != nil {
		return err
	}

	err = actualBuilder.Configure()
	if err != nil {
		return err
	}
	err = actualBuilder.Build()
	if err != nil {
		return err
	}
	err = actualBuilder.Install()
	if err != nil {
		return err
	}

	return nil
}

func RegisterBuilder(identifier string, maker MakeBuilder) error {
	_, found := builders[identifier]
	if found {
		return errAlreadyRegistered
	}
	builders[identifier] = maker
	return nil
}
