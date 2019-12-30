PROJECT_DIR  = $(shell readlink -f .)
BUILD_DIR    = "$(PROJECT_DIR)/build"
RECEIVER_DIR = "$(PROJECT_DIR)/cmd/receiver"
RECEIVER_BIN = "$(BUILD_DIR)/wathola-receiver"
SENDER_DIR   = "$(PROJECT_DIR)/cmd/sender"
SENDER_BIN   = "$(BUILD_DIR)/wathola-sender"

GO           ?= go
RICHGO       ?= rich$(GO)

RESET          = \033[0m
make_std_color = \033[3$1m      # defined for 1 through 7
make_color     = \033[38;5;$1m  # defined for 1 through 255
BLUE   = $(strip $(call make_color,38))
PINK   = $(strip $(call make_color,210))
RED    = $(strip $(call make_color,206))
GREEN  = $(strip $(call make_color,120))
DGREEN = $(strip $(call make_color,106))
GRAY   = $(strip $(call make_color,224))

.PHONY: default
default: binaries

.PHONY: builddeps
builddeps:
	@GO111MODULE=off $(GO) get github.com/kyoh86/richgo
	@GO111MODULE=off $(GO) get github.com/mgechev/revive

.PHONY: builddir
builddir:
	@mkdir -p build

.PHONY: clean
clean: builddeps
	@echo " $(GRAY)🛁 Cleaning$(RESET)"
	@rm -frv $(BUILD_DIR)

.PHONY: check
check: builddeps
	@echo " $(PINK)🛂 Checking$(RESET)"
	revive -config revive.toml -formatter stylish ./...

.PHONY: test
test: builddir check
	@echo " $(GREEN)✔️ Testing$(RESET)"
	$(RICHGO) test -v -covermode=count -coverprofile=build/coverage.out ./...

.PHONY: binaries
binaries: builddir test
	@echo " $(BLUE)🔨 Building$(RESET)"
	CGO_ENABLED=0 GOOS=linux $(RICHGO) build -o $(SENDER_BIN) $(SENDER_DIR)
	CGO_ENABLED=0 GOOS=linux $(RICHGO) build -o $(RECEIVER_BIN) $(RECEIVER_DIR)

.PHONY: run
run: builddir binaries
	@echo " $(RED)🏃 Running$(RESET)"
	cd $(BUILD_DIR) && $(RECEIVER_BIN) $(args)

.PHONY: images
images:
	@echo " $(BLUE)🔨 Building images$(RESET)"
	docker build --tag cardil/wathola-receiver --file images/receiver/Dockerfile .
	docker build --tag cardil/wathola-sender --file images/sender/Dockerfile .
