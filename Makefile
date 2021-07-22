TMP_DIR := tmp
WEB_DIR := web
WEB_SCRIPT_DIR := $(WEB_DIR)/script
WEB_STYLE_DIR := $(WEB_DIR)/style
DIST_DIR := dist
DIST_SCRIPT_DIR := $(DIST_DIR)/script
DIST_STYLE_DIR := $(DIST_DIR)/style
DIST_EMBED_FILE := $(DIST_DIR)/dist.go
DIST_EMBED_CONTENT := "package dist\n\nimport \"embed\"\n\n//go:embed *\nvar _ embed.FS"


airless := false
.PHONY: debug
debug: clean prepare
	ln -sf bundle.css $(DIST_STYLE_DIR)/bundle.min.css
	ln -sf ../../$(WEB_SCRIPT_DIR)/theme.js $(DIST_SCRIPT_DIR)/theme.min.js
	ln -sf ../../$(WEB_SCRIPT_DIR)/main.js $(DIST_SCRIPT_DIR)/bundle.min.js
	@if [ $(airless) = true ]; then \
		$(MAKE) sass watchsass=true; \
	else \
		$(MAKE) -j2 sass watchsass=true air; \
	fi;


.PHONY: release
release: clean prepare
	echo $(DIST_EMBED_CONTENT) > $(DIST_EMBED_FILE)
	$(MAKE) sass
	$(MAKE) uglifycss
	$(MAKE) uglifyjs
	go build -o $(TMP_DIR)/blogo ./cmd/blogo/
	GIN_MODE=release $(TMP_DIR)/blogo


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
	$(call install-if-needed,sass,brew install sass/sass/sass)
	@if [ $(watchsass) = true ]; then \
		sass --watch $(WEB_STYLE_DIR)/main.scss $(DIST_STYLE_DIR)/bundle.css; \
	else \
		sass $(WEB_STYLE_DIR)/main.scss $(DIST_STYLE_DIR)/bundle.css; \
	fi;


.PHONY: migrate
migrate:
	$(call install-if-needed,migrate,brew install golang-migrate)
	migrate -path store/migration -database 'mysql://hulz:xxxxxx@tcp(localhost:3306)/blogo?charset=utf8mb4&parseTime=true&loc=Local' -verbose $(cmd)


.PHONY: uglifycss
uglifycss:
	$(call install-if-needed,uglifycss,npm install uglifycss -g)
	uglifycss $(DIST_STYLE_DIR)/bundle.css > $(DIST_STYLE_DIR)/bundle.min.css


.PHONY: uglifyjs
uglifyjs:
	$(call install-if-needed,uglifyjs,npm install uglify-js -g)
	uglifyjs --compress --mangle --toplevel $(WEB_SCRIPT_DIR)/theme.js > $(DIST_SCRIPT_DIR)/theme.min.js
	uglifyjs --compress --mangle --toplevel $(WEB_SCRIPT_DIR)/main.js > $(DIST_SCRIPT_DIR)/bundle.min.js


.PHONY: test
test:
	go test ./...


.PHONY: clean
clean:
	find $(DIST_DIR)/* -not -name "*.go" -prune -exec rm -rf {} 2>/dev/null \;
	rm -rf $(TMP_DIR)
	go mod tidy


# Install the $(1) with $(2) if it hasn't been installed.
define install-if-needed
	@if ! command -v $(1) &> /dev/null; then \
		echo "Start installing $(1)..."; \
		$(2) && \
		echo "Finished installing $(1)!"; \
	fi
endef
