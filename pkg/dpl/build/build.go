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
	defaultInstallMethod string = "default"
	defaultInstallPath   string = "install"

	buildToolKey     string = "build.tool"
	installPathKey   string = "build.install_path"
	installMethodKey string = "build.install_method"
)

var (
	BuildTask common.Task = common.Task{
		Name: "build",
		Work: doFullBuild,
	}

	builders map[string]MakeBuilder = map[string]MakeBuilder{}

	errAlreadyRegistered    error = fmt.Errorf("builder already registered")
	errTooManyBuilders      error = fmt.Errorf("multiple builders")
	errInvalidBuilder       error = fmt.Errorf("invalid builder")
	errInvalidInstallMethod error = fmt.Errorf("unknown installation method")

	buildSteps []buildStep = []buildStep{
		doConfigure,
		doBuild,
		doInstall,
	}

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
	installMethod, err := dpl.GetSingleComponentValueOrDefault(component, installMethodKey, defaultInstallMethod)
	if err != nil {
		return err
	}
	installer, found := installHandlers[installMethod]
	if !found {
		return errInvalidInstallMethod
	}
	return installer(builder, component, config)
}

func defaultInstaller(builder Builder, component dpl.Component, config *BuildConfig) error {
	installPath, err := dpl.GetSingleComponentValueOrDefault(component, installPathKey, defaultInstallPath)
	if err != nil {
		return err
	}
	return builder.Install(path.Join(component.GetWorkDir(), installPath))
}

func noneInstaller(Builder, dpl.Component, *BuildConfig) error {
	return nil
}

func doFullBuild(component dpl.Component) error {
	builder, err := dpl.GetSingleComponentValue(component, buildToolKey)
	if err != nil {
		if _, ok := err.(*dpl.MissingKeyError); ok {
			log.Printf("No builder specified for '%v'", component.Name())
			return nil
		} else if _, ok := err.(*dpl.TooManyValuesError); ok {
			return errTooManyBuilders
		}
		return err
	}
	builderMaker, found := builders[builder]
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

	err = findAllArtifacts(component, buildArtifactPath, component.GetWorkDir())
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
