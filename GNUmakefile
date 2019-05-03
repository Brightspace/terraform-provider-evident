TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=evident

default: build

build: fmtcheck
	go build

build/win: fmtcheck
	GOOS=windows go build

get:
	go get -v -d

docker:
	docker run --rm -it \
		-v $(PWD):/go/src/github.com/Brightspace/terraform-provider-evident \
		--workdir /go/src/github.com/Brightspace/terraform-provider-evident \
		golang bash

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 EVIDENT_ACCESS_KEY=something EVIDENT_SECRET_KEY=somethingelse go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

.PHONY: build test testacc vet fmt fmtcheck errcheck

