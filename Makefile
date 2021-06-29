DIST := dist
WEBSITE := website
DIST_SCRIPT_DIR := $(DIST)/script
DIST_STYLE_DIR := $(DIST)/style
WEBSITE_SCRIPT_DIR := $(WEBSITE)/script
WEBSITE_STYLE_DIR := $(WEBSITE)/style

BREW_INSTALLED := $(shell command -v brew)

.PHONY: all debug release setup clean

all:
	@echo "make all"

debug: clean setup
	mkdir -p $(DIST_SCRIPT_DIR)
	mkdir -p $(DIST_STYLE_DIR)
	ln -sf ../../$(WEBSITE_SCRIPT_DIR)/main.js $(DIST_SCRIPT_DIR)/bundle.js
	sass --watch $(WEBSITE_STYLE_DIR)/main.scss $(DIST_STYLE_DIR)/bundle.css

release: clean setup
	@echo "make release"

setup:
# Install Sass if not exists.
	@if ! command -v sass &> /dev/null; then \
		echo "Installing Sass..."; \
		brew install sass/sass/sass && \
		echo "Sass installed: `sass --version`"; \
	fi

clean:
	go mod tidy
	rm -rf $(DIST)
