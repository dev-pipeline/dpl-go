package checkout

import (
	"fmt"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
	"github.com/dev-pipeline/dpl-go/pkg/dpl/scm"

	"github.com/dev-pipeline/dpl-go/internal/common"
)

var (
	Task = common.Task{
		Name: "scm",
		Work: checkout,
	}
)

func checkout(component dpl.Component) error {
	scmUris, err := component.ExpandValue("scm.uri")
	if err != nil {
		return err
	}
	for _, uri := range scmUris {
		scmInfo, err := scm.BuildScmInfo(uri)
		if err != nil {
			return err
		}
		scmBuilder := scm.GetHandler(scmInfo.Scheme)
		if scmBuilder == nil {
			return fmt.Errorf("no handler for %v", scmInfo.Scheme)
		}
		handler, err := scmBuilder(component)
		if err != nil {
			return err
		}
		err = handler.Checkout(scmInfo)
		if err != nil {
			return err
		}
	}
	return nil
}
