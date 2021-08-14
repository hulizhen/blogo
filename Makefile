BUILD_DIR := build
WEB_DIR := web
WEB_SCRIPT_DIR := $(WEB_DIR)/script
WEB_STYLE_DIR := $(WEB_DIR)/style
DIST_DIR := dist
DIST_SCRIPT_DIR := $(DIST_DIR)/script
DIST_STYLE_DIR := $(DIST_DIR)/style


airless := false
.PHONY: debug
debug: prepare
	ln -sf bundle.css $(DIST_STYLE_DIR)/bundle.min.css
	ln -sf ../../$(WEB_SCRIPT_DIR)/theme.js $(DIST_SCRIPT_DIR)/theme.min.js
	ln -sf ../../$(WEB_SCRIPT_DIR)/main.js $(DIST_SCRIPT_DIR)/bundle.min.js
	@if [ $(airless) = true ]; then \
		$(MAKE) sass watchsass=true; \
	else \
		$(MAKE) -j2 sass watchsass=true air; \
	fi;


.PHONY: release
release: prepare
	$(MAKE) sass
	uglifycss $(DIST_STYLE_DIR)/bundle.css > $(DIST_STYLE_DIR)/bundle.min.css
	uglifyjs --compress --mangle --toplevel $(WEB_SCRIPT_DIR)/theme.js > $(DIST_SCRIPT_DIR)/theme.min.js
	uglifyjs --compress --mangle --toplevel $(WEB_SCRIPT_DIR)/main.js > $(DIST_SCRIPT_DIR)/bundle.min.js
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/blogod ./cmd/blogod/


.PHONY: prepare
prepare:
	mkdir -p $(DIST_STYLE_DIR)
	mkdir -p $(DIST_SCRIPT_DIR)


.PHONY: air
air:
	$(call install-if-needed,air,curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin)
	air


watchsass := false
.PHONY: sass
sass:
	@if [ $(watchsass) = true ]; then \
		sass --watch $(WEB_STYLE_DIR)/main.scss $(DIST_STYLE_DIR)/bundle.css; \
	else \
		sass $(WEB_STYLE_DIR)/main.scss $(DIST_STYLE_DIR)/bundle.css; \
	fi;


.PHONY: test
test:
	go test ./...


.PHONY: clean
clean:
	rm -rf __debug_bin
	rm -rf $(DIST_DIR)
	rm -rf $(BUILD_DIR)
	go mod tidy


# Install the $(1) with $(2) if it hasn't been installed.
define install-if-needed
	@if ! command -v $(1); then \
		echo "Start installing $(1)..."; \
		$(2) && \
		echo "Finished installing $(1)!"; \
	fi
endef
