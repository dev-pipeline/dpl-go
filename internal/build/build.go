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
		Work: doFullBuild,
	}

	builders map[string]MakeBuilder = map[string]MakeBuilder{}

	errAlreadyRegistered      error = fmt.Errorf("builder already registered")
	errTooManyBuilders        error = fmt.Errorf("multiple builders")
	errInvalidBuilder         error = fmt.Errorf("invalid builder")
	errMultipleInstallMethods error = fmt.Errorf("multiple install methods specified")
	errInvalidInstallMethod   error = fmt.Errorf("unknown installation method")

	buildSteps []buildStep = []buildStep{
		doConfigure,
		doBuild,
		doInstall,
	}

	defaultInstallMethod []string = []string{"default"}

	installHandlers map[string]installFn = map[string]installFn{
		"default": defaultInstaller,
		"none":    noneInstaller,
	}
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

type buildStep func(Builder, dpl.Component, *BuildConfig) error

type installFn func(Builder, dpl.Component, *BuildConfig) error

func doConfigure(builder Builder, component dpl.Component, config *BuildConfig) error {
	return builder.Configure(config)
}

func doBuild(builder Builder, component dpl.Component, config *BuildConfig) error {
	return builder.Build(config)
}

func doInstall(builder Builder, component dpl.Component, config *BuildConfig) error {
	installMethod, err := component.ExpandValue("build.install_method")
	if err != nil {
		return err
	}
	if len(installMethod) > 1 {
		return errMultipleInstallMethods
	}
	if len(installMethod) == 0 {
		installMethod = defaultInstallMethod
	}
	installer, found := installHandlers[installMethod[0]]
	if !found {
		return errInvalidInstallMethod
	}
	return installer(builder, component, config)
}

func defaultInstaller(builder Builder, component dpl.Component, config *BuildConfig) error {
	return builder.Install(path.Join(component.GetWorkDir(), "install"))
}

func noneInstaller(Builder, dpl.Component, *BuildConfig) error {
	return nil
}

func doFullBuild(component dpl.Component) error {
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
	config := BuildConfig{
		Env: os.Environ(),
	}
	for k, v := range envChanges.prependValues {
		config.Env = prependEnvironment(config.Env, k, v)
	}
	for k, v := range envChanges.appendValues {
		config.Env = appendEnvironment(config.Env, k, v)
	}

	for i := range buildSteps {
		err := buildSteps[i](actualBuilder, component, &config)
		if err != nil {
			return err
		}
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
