#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

SHELL := /bin/bash -o pipefail
UNAME ?= $(shell uname)
YARN_EXEC ?= $(shell which yarn)
GO_EXEC ?= $(shell which go)

VERSION ?= latest
RELEASE_SRC = apache-apisix-dashboard-${VERSION}-src

export GO111MODULE=on

### help:		Show Makefile rules
.PHONY: help
help:
	@echo Makefile rules:
	@echo
	@grep -E '^### [-A-Za-z0-9_]+:' Makefile | sed 's/###/   /'


### build:		Build the Apache APISIX Dashboard, including web and manager-api
.PHONY: build
build: web-default api-default
	api/build.sh && cd ./web && export CYPRESS_INSTALL_BINARY=0  && yarn install --ignore-scripts && yarn build  && mkdir -p ../output/logs


.PHONY: web-default
web-default:
ifeq ("$(wildcard $(YARN_EXEC))", "")
	@echo "ERROR: Need to install yarn first"
	exit 1
endif


.PHONY: api-default
api-default:
ifeq ("$(wildcard $(GO_EXEC))", "")
	@echo "ERROR: Need to install golang 1.15+ first"
	exit 1
endif


### dag-lib:            Download the dag-lib
.PHONY: dag-lib
dag-lib:
ifeq ("$(wildcard api/dag-to-lua/dag-to-lua.lua)", "")
	curl -Lso /tmp/v1.1.tar.gz https://github.com/api7/dag-to-lua/archive/v1.1.tar.gz
	tar -zxvf /tmp/v1.1.tar.gz -C /tmp
	mkdir ./api/dag-to-lua
	cp -r /tmp/dag-to-lua-1.1/lib/* ./api/dag-to-lua
endif


### api-test:		Run the tests of manager-api
.PHONY: api-test
api-test: api-default dag-lib
	cd api/ && ENV=test go test -v -count=1 -race -cover -coverprofile=coverage.txt -covermode=atomic ./... && go tool cover -html=coverage.txt -o coverage.html


### api-run:		Run the manager-api in develop mode
.PHONY: api-run
api-run: api-default
	api/build.sh --dry-run


### api-stop:		Stop the manager-api
api-stop:
	cd api && go run -ldflags "${GOLDFLAGS}" ./main.go stop

### go-lint:		Lint Go source codes
.PHONY: go-lint
go-lint: ## Run the golangci-lint application (install if not found)
	@#Brew - MacOS
	@if [ "$(shell command -v golangci-lint)" = "" ] && [ "$(shell command -v brew)" != "" ] && [ "$(UNAME)" = "Darwin" ]; then brew install golangci-lint; fi;
	@#has sudo
	@if [ "$(shell command -v golangci-lint)" = "" ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.32.0 && sudo cp ./bin/golangci-lint $(go env GOPATH)/bin/; fi;
	@echo "running golangci-lint..."
	@cd api && golangci-lint run --tests=false ./...


### license-check:	Check source codes for Apache License
.PHONY: license-check
license-check:
ifeq ("$(wildcard .actions/openwhisk-utilities/scancode/scanCode.py)", "")
	git clone https://github.com/apache/openwhisk-utilities.git .actions/openwhisk-utilities
	cp .actions/ASF* .actions/openwhisk-utilities/scancode/
endif
	.actions/openwhisk-utilities/scancode/scanCode.py --config .actions/ASF-Release.cfg ./


### release-src:        Package source codes for release
.PHONY: release-src
release-src:
	./utils/version-check.sh ${VERSION}
	git clean -Xdf
	rm -f ./.githash && git log --pretty=format:"%h" -1 > ./.githash
	tar -zcvf $(RELEASE_SRC).tgz \
	--exclude .github \
	--exclude .git \
	--exclude .gitattributes \
	--exclude .idea \
	--exclude .vscode \
	--exclude .gitignore \
	--exclude .DS_Store \
	--exclude docs \
	--exclude release \
	--exclude api/internal/core/store/validate_mock.go \
	--exclude api/internal/core/storage/storage_mock.go \
	.

	gpg --batch --yes --armor --detach-sig $(RELEASE_SRC).tgz
	shasum -a 512 $(RELEASE_SRC).tgz > $(RELEASE_SRC).tgz.sha512

	mkdir -p release
	mv $(RELEASE_SRC).tgz release/$(RELEASE_SRC).tgz
	mv $(RELEASE_SRC).tgz.asc release/$(RELEASE_SRC).tgz.asc
	mv $(RELEASE_SRC).tgz.sha512 release/$(RELEASE_SRC).tgz.sha512
