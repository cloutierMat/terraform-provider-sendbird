ifneq ($(origin PKG), undefined)
	PKG_NAME = internal/service/$(PKG)
	TEST = ./$(PKG_NAME)/...
endif

ifneq ($(origin TESTS), undefined)
	RUNARGS = -run='$(TESTS)'
endif

default: build

build:
	go install

docs: fmtcheck
	terraform fmt -recursive ./examples/
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

fmt: 
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./internal ./main.go

fmtcheck:
	@sh -c "'$(CURDIR)/.scripts/gofmtcheck.sh'"

test: fmtcheck
	go test ./$(PKG_NAME)/...

testacc: fmtcheck
	TF_ACC=1 go test ./$(PKG_NAME)/... -v $(RUNARGS) $(TESTARGS)
