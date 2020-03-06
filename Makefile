-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"


######################
# Testing
######################

GINKGO := artifacts/ginkgo/bin/ginkgo
$(GINKGO):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get -modfile tools.mod github.com/onsi/ginkgo/ginkgo

test:: artifacts/test/config.coverprofile $(GINKGO)
	-@mkdir -p "artifacts/test"
	$(GINKGO) -outputdir "artifacts/test/" -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --compilers=2 --nodes=2

######################
# Linting
######################

MISSPELL := artifacts/misspell/bin/misspell
$(MISSPELL):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get -modfile tools.mod github.com/client9/misspell/cmd/misspell

GOLANGCILINT := artifacts/golangci-lint/bin/golangci-lint
$(GOLANGCILINT):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(MF_PROJECT_ROOT)/$(@D)" v1.23.8

STATICCHECK := artifacts/staticcheck/bin/staticcheck
$(STATICCHECK):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get -modfile tools.mod honnef.co/go/tools/cmd/staticcheck

.PHONY: lint
lint:: $(MISSPELL) $(GOLANGCILINT) $(STATICCHECK)
	go vet ./...
	golint -set_exit_status ./...
	$(MISSPELL) -w -error -locale UK ./...
	$(GOLANGCILINT) run --enable-all ./...
	$(STATICCHECK) -checks all -fail "all,-U1001" ./...

