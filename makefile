# Go parameters
GO_CMD=go
# GO_BUILD=$(GO_CMD) build
GO_BUILD = CGO_ENABLED=0 $(GO_CMD) build -a -ldflags '-s -w'
BINARY_NAME=GC2

# List of target OS/Arch combinations
TARGETS=\
    linux/amd64 \
    linux/arm64 \
    windows/amd64 \
    darwin/amd64 \
    darwin/arm64

# Default target executed when no arguments are given to make.
.PHONY: all
all: build

# Cross compile for all targets
.PHONY: build
build: $(TARGETS)

# Pattern rule to build for each target
$(TARGETS):
	@mkdir -p bin/$(subst /,_,$@)
	GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@)) $(GO_BUILD) -o bin/$(subst /,_,$@)/$(BINARY_NAME)

# Clean build files
.PHONY: clean
clean:
	rm -rf bin

# List of available targets
.PHONY: info
info:
	@echo "Binaries will be created for the available build targets:"
	@for target in $(TARGETS); do \
	    echo "  $$target"; \
	done
