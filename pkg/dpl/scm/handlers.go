package scm

import (
	"fmt"
	"net/url"

	"github.com/dev-pipeline/dpl-go/pkg/dpl"
)

type ScmInfo struct {
	Scheme    string
	Path      string
	Arguments map[string]string
}

type ScmHandler interface {
	Checkout(ScmInfo, dpl.Component) error
}

var (
	scms map[string]ScmHandler = map[string]ScmHandler{}
)

func AddHandler(protocol string, handler ScmHandler) error {
	scms[protocol] = handler
	return nil
}

func GetHandler(protocol string) ScmHandler {
	handler, found := scms[protocol]
	if found {
		return handler
	}
	return nil
}

func BuildScmInfo(uriInfo string) (ScmInfo, error) {
	scmInfo := ScmInfo{}
	uri, err := url.Parse(uriInfo)
	if err != nil {
		return scmInfo, err
	}

	trimmedPath, arguments, err := extractArguments(uri.Path, ";")
	if err != nil {
		return scmInfo, err
	}

	userPrefix := ""
	if uri.User != nil {
		userPrefix = fmt.Sprintf("%v@", uri.User)
	}
	scmInfo.Scheme = uri.Scheme
	scmInfo.Path = fmt.Sprintf("%v%v%v", userPrefix, uri.Host, trimmedPath)
	scmInfo.Arguments = arguments

	return scmInfo, nil
}
