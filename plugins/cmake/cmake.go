package cmake

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/build"
)

type cmakeBuilder struct {
	component dpl.Component
}

type cmakeFlags struct {
	args []string
	env  []string
}

func (cb cmakeBuilder) runCmake(cf cmakeFlags) error {
	cmd := exec.Command("cmake", cf.args...)
	cmd.Dir = cb.component.GetWorkDir()
	if len(cf.env) > 0 {
		cmd.Env = cf.env
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing cmake: %v", string(output))
		return err
	}
	return nil
}

func (cb cmakeBuilder) Configure(config *build.BuildConfig) error {
	err := os.MkdirAll(cb.component.GetWorkDir(), 0755)
	if err != nil {
		return err
	}
	args := []string{}
	keys := cb.component.KeyNames()
	for i := range keys {
		fn, found := flagHandlers[keys[i]]
		if found {
			values, err := cb.component.ExpandValues(keys[i])
			if err != nil {
				return err
			}
			flags, err := fn(keys[i], values)
			if err != nil {
				return err
			}
			args = append(args, flags...)
		}
	}
	return cb.runCmake(cmakeFlags{
		args: append(args, cb.component.GetSourceDir()),
		env:  config.Env,
	})
}

func (cb cmakeBuilder) Build(config *build.BuildConfig) error {
	return cb.runCmake(cmakeFlags{
		args: []string{
			"--build",
			cb.component.GetWorkDir(),
		},
		env: config.Env,
	})
}

func (cb cmakeBuilder) Install(destdir string) error {
	env := os.Environ()
	if len(destdir) > 0 {
		env = append(env, fmt.Sprintf("DESTDIR=%v", destdir))
	}
	return cb.runCmake(cmakeFlags{
		args: []string{
			"--build",
			cb.component.GetWorkDir(),
			"--target",
			"install",
		},
		env: env,
	})
}

func init() {
	build.RegisterBuilder("cmake", func(component dpl.Component) (build.Builder, error) {
		return &cmakeBuilder{
			component: component,
		}, nil
	})
}
