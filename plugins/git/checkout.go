package git

import (
	// "encoding/hex"
	"fmt"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/dev-pipeline/dpl-go/pkg/dpl/scm"
)

var (
	errFoundRef       error = fmt.Errorf("found ref")
	errCantResolveRef error = fmt.Errorf("can't resolve reference")
)

func mapRef(ref *plumbing.Reference) (*gogit.CheckoutOptions, error) {
	if ref.Name().IsBranch() {
		return &gogit.CheckoutOptions{
			Branch: ref.Name(),
		}, nil
	}
	if ref.Name().IsTag() {
		return &gogit.CheckoutOptions{
			Hash: ref.Hash(),
		}, nil
	}
	return nil, errCantResolveRef
}

func tryHappyCheckout(r *gogit.Repository, ref string) (*gogit.CheckoutOptions, error) {
	reference, err := r.Reference(plumbing.ReferenceName(ref), true)
	if err == nil {
		return mapRef(reference)
	}
	return nil, errCantResolveRef
}

func findReference(r *gogit.Repository, ref string) (*gogit.CheckoutOptions, error) {
	refs, err := r.References()
	if err != nil {
		return nil, err
	}
	searchRef := fmt.Sprintf("/%v", ref)
	var realRef *plumbing.Reference
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if strings.HasSuffix(string(ref.Name()), searchRef) {
			realRef = ref
			return errFoundRef
		}
		return nil
	})

	if err == errFoundRef {
		return mapRef(realRef)
	}
	return nil, errCantResolveRef
}

func makeCheckoutOptions(r *gogit.Repository, info scm.ScmInfo) (*gogit.CheckoutOptions, error) {
	targetRef, found := info.Arguments["ref"]
	if !found {
		// if no ref was specified, then anything is fine
		return &gogit.CheckoutOptions{}, nil
	}
	// If ref was specified, then we need to find it.  go-git wants fully resolved references,
	// so we'll start by hoping for the best.
	options, err := tryHappyCheckout(r, targetRef)
	if err == nil {
		return options, err
	}

	// If the provided ref isn't good enough on its own, then try to turn it into a reference.
	// This should find either a branch or tag, based on whatever pattern hits first.
	options, err = findReference(r, targetRef)
	if err == nil {
		return options, err
	}

	// If we still don't have a resolved reference, all that's left is to hope it's
	// a commit hash (TODO in the future).
	return nil, errCantResolveRef
	// not a branch and not a tag, so see if we can get a hash out of it
	// before iterating commits, make sure we can actuaully
	/*
		_, err = hex.DecodeString(targetRef)
		if err != nil {
			return nil, errCantResolveRef
		}

		commits, err := r.CommitObjects()
		if err != nil {
			return err
		}
		var foundHash plumbing.Hash
		err = commits.ForEach(func ())

		return nil, nil
	*/
}

func doCheckout(r *gogit.Repository, info scm.ScmInfo) error {
	options, err := makeCheckoutOptions(r, info)
	if err != nil {
		return err
	}
	wt, err := r.Worktree()
	if err != nil {
		return err
	}
	err = wt.Checkout(options)
	if err != nil {
		return err
	}
	return nil
}
