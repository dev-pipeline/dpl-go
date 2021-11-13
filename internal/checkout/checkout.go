package checkout

import (
	"errors"
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
		handler := scm.GetHandler(scmInfo.Scheme)
		if handler == nil {
			return errors.New(fmt.Sprintf("No handler for %v", scmInfo.Scheme))
		}
		err = handler.Checkout(scmInfo, component)
		if err != nil {
			return err
		}
	}
	return nil
}
