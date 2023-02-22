GO_BUILD_PATH ?= bin/go-vksave

TOOLS_BIN_DIR ?= tools/bin

GO_LINT_TOOL = $(TOOLS_BIN_DIR)/golangci-lint

TOOLS_MODFILE = tools/go.mod


#.PHONY: lint
#lint: $(GO_LINT_TOOL)
#	$(GO_LINT_TOOL) --version \
#    $(GO_LINT_TOOL) run --sort-results --max-issues-per-linter=0 --max-same-issues=0 --print-resources-usage -v
#
#
#.PHONY: test
#test:
#	go test

.PHONY: build
build:
	go build -o $(GO_BUILD_PATH) cmd/main.go

.PHONY: run
run:
	./$(GO_BUILD_PATH)

.PHONY: clean
clean:
	@rm -rf ./bin

#.PHONY: install-tools
#install-tools: $(GO_LINT_TOOL))
#
#TOOLS_MODFILE = tools/go.mod
#define install-go-tool
#	$(MAYBE_GO_BUILDER_EXEC) go build \
#		-o $(TOOLS_BIN_DIR) \
#		-ldflags "-s -w" \
#		-modfile $(TOOLS_MODFILE)
#endef
#
#$(GO_LINT_TOOL): | $(TOOLS_BIN_DIR)
#	$(install-go-tool) github.com/golangci/golangci-lint/cmd/golangci-lint