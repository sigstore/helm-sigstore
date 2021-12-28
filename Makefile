#
# Copyright 2021 The Sigstore Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifdef MAKEFILE_SHELL
SHELL=$(MAKEFILE_SHELL)
endif
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

BINDIR     := $(CURDIR)/bin
DISTDIR    := $(CURDIR)/dist
BINNAME    := helm-sigstore

# Set version variables for LDFLAGS
GO := go
CGO_ENABLED := 0
GIT_VERSION ?= $(shell git describe --tags --always --dirty)
GIT_HASH ?= $(shell git rev-parse HEAD)
DATE_FMT = +'%Y-%m-%dT%H:%M:%SZ'
SOURCE_DATE_EPOCH ?= $(shell git log -1 --pretty=%ct)
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
GIT_TREESTATE = "clean"
DIFF = $(shell git diff --quiet >/dev/null 2>&1; if [ $$? -eq 1 ]; then echo "1"; fi)
ifeq ($(DIFF), 1)
    GIT_TREESTATE = "dirty"
endif

PKG=github.com/sigstore/helm-sigstore/cmd
LDFLAGS=-X $(PKG).gitVersion=$(GIT_VERSION) -X $(PKG).gitCommit=$(GIT_HASH) -X $(PKG).gitTreeState=$(GIT_TREESTATE) -X $(PKG).buildDate=$(BUILD_DATE)

.PHONY: all lint test clean sigstore build

all: sigstore

SRCS = $(shell find . -type f -name '*.go')

sigstore: $(SRCS)
	CGO_ENABLED=$(CGO_ENABLED) go build -trimpath -ldflags "$(LDFLAGS)" -o '$(BINDIR)/$(BINNAME)' .

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
golangci-lint:
	rm -f $(GOLANGCI_LINT) || :
	set -e ;\
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell dirname $(GOLANGCI_LINT)) v1.39.0 ;\

lint: golangci-lint ## Runs golangci-lint linter
	$(GOLANGCI_LINT) run  -n

test: ## Runs go tests
	go test ./...

##################
# release section
##################

.PHONY: release
release: ## Runs goreleaser in release mode
	LDFLAGS="$(LDFLAGS)" goreleaser release --rm-dist

# used when need to validate the goreleaser
.PHONY: snapshot
snapshot: ## Runs goreleaser in snapshot mode
	LDFLAGS="$(LDFLAGS)" goreleaser release --skip-sign --skip-publish --snapshot --rm-dist

clean:
	rm -rf $(BINDIR)
	rm -rf $(DISTDIR)

##################
# help
##################

help: # Display help
	@awk -F ':|##' \
		'/^[^\t].+?:.*?##/ {\
			printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
		}' $(MAKEFILE_LIST) | sort
