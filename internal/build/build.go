package build

import (
	"fmt"
	"log"
	"os"
	"path"

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

type BuildConfig struct {
	Env []string
}

type Builder interface {
	Configure(*BuildConfig) error
	Build(*BuildConfig) error
	Install(string) error
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

	envChanges, err := makeEnvMap(component)
	if err != nil {
		return err
	}
	config := &BuildConfig{
		Env: os.Environ(),
	}
	for k, v := range envChanges.prependValues {
		config.Env = prependEnvironment(config.Env, k, v)
	}
	for k, v := range envChanges.appendValues {
		config.Env = appendEnvironment(config.Env, k, v)
	}

	err = actualBuilder.Configure(config)
	if err != nil {
		return err
	}
	err = actualBuilder.Build(config)
	if err != nil {
		return err
	}
	err = actualBuilder.Install(path.Join(component.GetWorkDir(), "install"))
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
