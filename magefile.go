//go:build mage

package main

import (
	"context"

	//mage:import
	"github.com/dosquad/mage"
	"github.com/magefile/mage/mg"
)

// Local update, protoc, format, tidy, lint & test.
func Local(ctx context.Context) {
	mg.CtxDeps(ctx, mage.Test)
}

var Default = Local
