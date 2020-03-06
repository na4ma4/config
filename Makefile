-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

_GO_GTE_1_14 := $(shell expr `go version | cut -d' ' -f 3 | tr -d 'a-z' | cut -d'.' -f2` \>= 14)

ifeq "$(_GO_GTE_1_14)" "1"
_MODFILEARG := $(_MODFILEARG)
endif

######################
# Testing
######################

GINKGO := artifacts/ginkgo/bin/ginkgo
$(GINKGO):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) github.com/onsi/ginkgo/ginkgo

test:: $(GINKGO)
	-@mkdir -p "artifacts/test"
	$(GINKGO) -outputdir "artifacts/test/" -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --compilers=2 --nodes=2

######################
# Linting
######################

MISSPELL := artifacts/misspell/bin/misspell
$(MISSPELL):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) github.com/client9/misspell/cmd/misspell

GOLINT := artifacts/golint/bin/golint
$(GOLINT):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) golang.org/x/lint/golint

GOLANGCILINT := artifacts/golangci-lint/bin/golangci-lint
$(GOLANGCILINT):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(MF_PROJECT_ROOT)/$(@D)" v1.23.8

STATICCHECK := artifacts/staticcheck/bin/staticcheck
$(STATICCHECK):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) honnef.co/go/tools/cmd/staticcheck

.PHONY: lint
lint:: $(MISSPELL) $(GOLINT) $(GOLANGCILINT) $(STATICCHECK)
	go vet ./...
	$(GOLINT) -set_exit_status ./...
	$(MISSPELL) -w -error -locale UK ./...
	$(GOLANGCILINT) run --enable-all ./...
	$(STATICCHECK) -checks all -fail "all,-U1001" ./...

