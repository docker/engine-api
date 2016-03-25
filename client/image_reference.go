package client

import (
	"fmt"
	"github.com/docker/distribution/reference"
)

// parseReference parses the given references and return the repository and
// tag (if present) from it. If there is an error during parsing, it will
// return an error.
func parseReference(ref string) (string, string, error) {
	distributionRef, err := reference.ParseNamed(ref)
	if err != nil {
		return "", "", err
	}

	var tag string
	if tagged, isTagged := distributionRef.(reference.NamedTagged); isTagged {
		tag = tagged.Tag()
	}
	// This is in order to support a reference like :
	// tag@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
	if digested, isDigested := distributionRef.(reference.Digested); isDigested {
		tag = fmt.Sprintf("%s@%s", tag, digested.Digest().String())
	}
	return distributionRef.Name(), tag, nil
}
