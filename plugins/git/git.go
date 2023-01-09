package git

import (
	"log"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/scm"
)

type gitHandler struct {
	component dpl.Component
}

func (gh *gitHandler) Checkout(info scm.ScmInfo) error {
	r, err := getRepository(gh.component.GetSourceDir(), info)
	if err != nil {
		return err
	}
	err = doCheckout(r, info)
	return err
}

func makeGit(component dpl.Component) (scm.ScmHandler, error) {
	return &gitHandler{
		component: component,
	}, nil
}

func init() {
	err := scm.AddHandler("git", makeGit)
	if err != nil {
		log.Fatalf("Error registering git handler: %v", err)
	}
}
