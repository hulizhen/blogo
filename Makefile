DIST_DIR := dist
WEBSITE_DIR := website
DIST_SCRIPT_DIR := $(DIST_DIR)/script
DIST_STYLE_DIR := $(DIST_DIR)/style
WEBSITE_SCRIPT_DIR := $(WEBSITE_DIR)/script
WEBSITE_STYLE_DIR := $(WEBSITE_DIR)/style


.PHONY: debug
debug:
	mkdir -p $(DIST_SCRIPT_DIR)
	mkdir -p $(DIST_STYLE_DIR)
	ln -sf ../../$(WEBSITE_SCRIPT_DIR)/main.js $(DIST_SCRIPT_DIR)/bundle.js
	$(MAKE) -j2 sass air


.PHONY: air
air:
	@if ! command -v air &> /dev/null; then \
		echo "Installing Air..."; \
		curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin && \
		echo "Air installed: `air -v`"; \
	fi
	air


.PHONY: sass
sass:
	@if ! command -v sass &> /dev/null; then \
		echo "Installing Sass..."; \
		brew install sass/sass/sass && \
		echo "Sass installed: `sass --version`"; \
	fi
	sass --watch $(WEBSITE_STYLE_DIR)/main.scss $(DIST_STYLE_DIR)/bundle.css


.PHONY: clean
clean:
	go mod tidy
	rm -rf $(DIST_DIR)
