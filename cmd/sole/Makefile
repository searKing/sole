# note: call scripts from /scripts

TARGET := sole
MAKEFILE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PKG_CONFIG_PATH := ${MAKEFILE_PATH}/../../pkgconfig
git_tag=$(shell git describe --long --tags --dirty --tags --always)
git_commit=$(shell git rev-parse HEAD)
git_build_time=$(shell TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')

GO := go
GO_DOWNLOAD_FLAGS := --insecure
GO_BUILD_FLAGS := -gcflags "-N -l" \
-ldflags "-s -X 'github.com/searKing/sole/internal/pkg/provider/viper.Version=${git_tag}' \
-X 'github.com/searKing/sole/internal/pkg/provider/viper.BuildTime=${git_build_time}' \
-X 'github.com/searKing/sole/internal/pkg/provider/viper.GitHash=${git_commit}'"

.PHONY: all
all: all-mod

.PHONY: all-mod
all-mod: go-get-mod clean build-mod

.PHONY: all-vendor
all-vendor: go-get-vendor clean build-vendor

.PHONY: build-vendor
build-vendor: GO_BUILD_FLAGS:=-mod vendor $(GO_BUILD_FLAGS)
build-vendor: git-info build-mod

.PHONY: build-mod
build-mod: git-info $(TARGET)

.PHONY: $(TARGET)
$(TARGET):
	@echo "  >  Building binary $(TARGET) ..."
	@PKG_CONFIG_PATH="${PKG_CONFIG_PATH}" ${GO} build ${GO_BUILD_FLAGS} -o $@ .


.PHONY: go-get-mod
go-get-mod:
	@echo "  >  Checking if there is any missing dependencies..."
	@PKG_CONFIG_PATH="${PKG_CONFIG_PATH}" ${GO} get ${GO_DOWNLOAD_FLAGS}
	@${GO} mod vendor

.PHONY: go-get-vendor
go-get-vendor: go-get-mod
	@echo "  >  Copying vendor from mod cache..."
	@${GO} mod vendor

.IGNORE: clean
clean:
	@echo "  >  Cleaning build cache"
	@${GO} clean -cache -testcache
	@if [ -f ${TARGET} ] ; then rm ${TARGET} ; fi


PREFIX :=
ifeq ($(DESTDIR),)
PREFIX=.
endif

.PHONY: install
install: install-mod

.PHONY: install-mod
install-mod: build-mod install-only

.PHONY: install-vendor
install-vendor: build-vendor install-only

.PHONY: install-only
install-only:
	@echo "  >  installing target to $(DESTDIR)$(PREFIX)"
	@mkdir -p "$(DESTDIR)$(PREFIX)/bin"
	@cp -Rv "$(TARGET)" "$(DESTDIR)$(PREFIX)/bin/$(TARGET)"
	@mkdir -p "$(DESTDIR)$(PREFIX)/conf"
	@cp -Rv "${MAKEFILE_PATH}/../../configs/$(TARGET).yml" "$(DESTDIR)$(PREFIX)/conf/$(TARGET).yml"
	@mkdir -p "$(DESTDIR)$(PREFIX)/web/webapp"
	@cp -Rv "${MAKEFILE_PATH}/../../web/webapp"/* "$(DESTDIR)$(PREFIX)/web/webapp/"
	@mkdir -p "$(DESTDIR)$(PREFIX)/website"
	@cp -Rv "${MAKEFILE_PATH}/../../website"/* "$(DESTDIR)$(PREFIX)/website/"

.PHONY: uninstall
uninstall:
	@echo "  >  uninstalling $(DESTDIR)$(PREFIX)"
	@if [[ -d "$(DESTDIR)$(PREFIX)/bin" ]]; then rm -Rv "$(DESTDIR)$(PREFIX)/bin"; fi
	@if [[ -d "$(DESTDIR)$(PREFIX)/conf" ]]; then rm -Rv "$(DESTDIR)$(PREFIX)/conf"; fi
	@if [[ -d "$(DESTDIR)$(PREFIX)/web" ]]; then rm -Rv "$(DESTDIR)$(PREFIX)/web"; fi
	@if [[ -d "$(DESTDIR)$(PREFIX)/website" ]]; then rm -Rv "$(DESTDIR)$(PREFIX)/website"; fi

.PHONY: git-info
git-info:
	@echo "  >  git_tag $(git_tag) ..."
	@echo "  >  git_commit $(git_commit) ..."
	@echo "  >  git_build_time $(git_build_time) ..."