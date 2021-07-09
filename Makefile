TMP_DIR := ./tmp
DIST_DIR := ./dist
WEBSITE_DIR := ./website
DIST_SCRIPT_DIR := $(DIST_DIR)/script
DIST_STYLE_DIR := $(DIST_DIR)/style
DIST_ASSET_DIR := $(DIST_DIR)/asset
WEBSITE_SCRIPT_DIR := $(WEBSITE_DIR)/script
WEBSITE_STYLE_DIR := $(WEBSITE_DIR)/style
WEBSITE_ASSET_DIR := $(WEBSITE_DIR)/asset
WEBSITE_STYLE_SRC := $(shell find $(DIST_STYLE_DIR) -name "*.css")
WEBSITE_SCRIPT_SRC := $(shell find $(WEBSITE_SCRIPT_DIR) -name "*.js")


airless := false
.PHONY: debug
debug: clean
	mkdir -p $(DIST_STYLE_DIR)
	mkdir -p $(DIST_SCRIPT_DIR)
	ln -sf bundle.css $(DIST_STYLE_DIR)/bundle.min.css
	ln -sf ../../$(WEBSITE_SCRIPT_DIR)/theme.js $(DIST_SCRIPT_DIR)/theme.min.js
	ln -sf ../../$(WEBSITE_SCRIPT_DIR)/main.js $(DIST_SCRIPT_DIR)/bundle.min.js
	ln -sf ../$(WEBSITE_ASSET_DIR) $(DIST_ASSET_DIR)
	@if [ $(airless) = true ]; then \
		$(MAKE) sass watchsass=true; \
	else \
		$(MAKE) -j2 sass watchsass=true air; \
	fi;


.PHONY: release
release: clean
	mkdir -p $(DIST_STYLE_DIR)
	mkdir -p $(DIST_SCRIPT_DIR)
	cp -a $(WEBSITE_ASSET_DIR) $(DIST_ASSET_DIR)
	$(MAKE) sass
	$(MAKE) uglifycss
	$(MAKE) uglifyjs
	go build -o $(TMP_DIR)/blogo ./cmd/blogo/
	GIN_MODE=release $(TMP_DIR)/blogo


.PHONY: air
air:
	$(call install-if-needed,air,curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin)
	air


watchsass := false
.PHONY: sass
sass:
	$(call install-if-needed,sass,brew install sass/sass/sass)
	@if [ $(watchsass) = true ]; then \
		sass --watch $(WEBSITE_STYLE_DIR)/main.scss $(DIST_STYLE_DIR)/bundle.css; \
	else \
		sass $(WEBSITE_STYLE_DIR)/main.scss $(DIST_STYLE_DIR)/bundle.css; \
	fi;


.PHONY: migrate
migrate:
	$(call install-if-needed,migrate,brew install golang-migrate)
	migrate -path store/migration -database 'mysql://hulz:xxxxxx@tcp(localhost:3306)/blogo?charset=utf8mb4&parseTime=true&loc=Local' -verbose $(cmd)


.PHONY: uglifycss
uglifycss: $(WEBSITE_STYLE_SRC)
	echo $(WEBSITE_STYLE_SRC)
	$(call install-if-needed,uglifycss,npm install uglifycss -g)
	uglifycss $(DIST_STYLE_DIR)/bundle.css > $(DIST_STYLE_DIR)/bundle.min.css


.PHONY: uglifyjs
uglifyjs: $(WEBSITE_SCRIPT_SRC)
	echo $(WEBSITE_SCRIPT_SRC)
	$(call install-if-needed,uglifyjs,npm install uglify-js -g)
	uglifyjs --compress --mangle --toplevel $(WEBSITE_SCRIPT_DIR)/theme.js > $(DIST_SCRIPT_DIR)/theme.min.js
	uglifyjs --compress --mangle --toplevel $(WEBSITE_SCRIPT_DIR)/main.js > $(DIST_SCRIPT_DIR)/bundle.min.js


.PHONY: test
test:
	go test ./...


.PHONY: clean
clean:
	rm -rf $(DIST_DIR)
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
