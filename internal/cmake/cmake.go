package cmake

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/dev-pipeline/dpl-go/internal/build"
	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type cmakeBuilder struct {
	sourceDir string
	workDir   string
}

type cmakeFlags struct {
	args []string
	env  []string
}

func (cb cmakeBuilder) runCmake(cf cmakeFlags) error {
	cmd := exec.Command("cmake", cf.args...)
	cmd.Dir = cb.workDir
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
	log.Printf("creating %v", cb.workDir)
	err := os.MkdirAll(cb.workDir, 0755)
	if err != nil {
		return err
	}
	return cb.runCmake(cmakeFlags{
		args: []string{
			cb.sourceDir,
		},
	})
}

func (cb cmakeBuilder) Build(config *build.BuildConfig) error {
	return cb.runCmake(cmakeFlags{
		args: []string{
			"--build",
			cb.workDir,
		},
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
			cb.workDir,
			"--target",
			"install",
		},
		env: env,
	})
}

func init() {
	build.RegisterBuilder("cmake", func(component dpl.Component) (build.Builder, error) {
		return &cmakeBuilder{
			sourceDir: component.GetSourceDir(),
			workDir:   component.GetWorkDir(),
		}, nil
	})
}
