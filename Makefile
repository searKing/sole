MAKEFILE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: all
all: all-mod

.PHONY: clean
	@echo "  >  "$@"ing $(DESTDIR)"
	$(MAKE) -C cmd/sole clean

.PHONY: all-mod
all-mod:generate
	$(MAKE) -C cmd/sole all-mod

.PHONY: all-vendor
all-vendor:generate
	$(MAKE) -C cmd/sole all-vendor

.PHONY: build-mod
build-mod:
	@echo "  >  "$@"ing $(DESTDIR)"
	$(MAKE) -C cmd/sole build-mod

.PHONY: build-vendor
build-vendor:
	@echo "  >  "$@"ing $(DESTDIR)"
	$(MAKE) -C cmd/sole build-vendor

.PHONY: generate
generate:
	@echo "  >  "$@"ing $(DESTDIR)"
	go generate github.com/searKing/sole/api/protobuf-spec

.PHONY: install
install: install-mod

DESTDIR :=
DESTDIR_PREFIX=$(shell echo $(DESTDIR) | head -c 1)
ifneq ($(DESTDIR_PREFIX),"/")
DESTDIR:=$(MAKEFILE_PATH)/$(DESTDIR)
endif

.PHONY: pack-mod
pack-mod: DESTDIR:=$(DESTDIR)/pack
pack-mod:install-mod
	@echo "  >  "$@" finished $(DESTDIR)"

.PHONY: pack-vendor
pack-vendor: DESTDIR:=$(DESTDIR)/pack
pack-vendor:install-vendor
	@echo "  >  "$@" finished $(DESTDIR)"

.PHONY: unpack
unpack: DESTDIR:=$(DESTDIR)/pack
unpack:uninstall
	@echo "  >  "$@" finished $(DESTDIR)"

.PHONY: install-mod
install-mod:
	@echo "  >  "$@" $(DESTDIR)"
	$(MAKE) -C cmd/sole install-mod DESTDIR=$(DESTDIR)

.PHONY: install-vendor
install-vendor:
	@echo "  >  "$@" $(DESTDIR)"
	$(MAKE) -C cmd/sole install-vendor DESTDIR=$(DESTDIR)

.PHONY: uninstall
uninstall:
	@echo "  >  "$@"ing $(DESTDIR)"
	$(MAKE) -C cmd/sole uninstall DESTDIR=$(DESTDIR)