MAKEFILE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: all
all: build

.PHONY: clean
	@echo "  >  "$@"ing $(DESTDIR)"
	$(MAKE) -C soletemplate/cmd/soletemplate clean

.PHONY: all
all:
	$(MAKE) -C soletemplate/cmd/soletemplate all

.PHONY: build
build:
	@echo "  >  "$@"ing $(DESTDIR)"
	$(MAKE) -C soletemplate/cmd/soletemplate build

.PHONY: install
install: install

DESTDIR :=
DESTDIR_PREFIX=$(shell echo $(DESTDIR) | head -c 1)
ifneq ($(DESTDIR_PREFIX),"/")
DESTDIR:=$(MAKEFILE_PATH)/$(DESTDIR)
endif

.PHONY: pack
pack: DESTDIR:=$(DESTDIR)/pack
pack:install
	@echo "  >  "$@" finished $(DESTDIR)"

.PHONY: unpack
unpack: DESTDIR:=$(DESTDIR)/pack
unpack:uninstall
	@echo "  >  "$@" finished $(DESTDIR)"

.PHONY: install
install:
	@echo "  >  "$@" $(DESTDIR)"
	$(MAKE) -C soletemplate/cmd/soletemplate install DESTDIR=$(DESTDIR)

.PHONY: uninstall
uninstall:
	@echo "  >  "$@"ing $(DESTDIR)"
	$(MAKE) -C soletemplate/cmd/soletemplate uninstall DESTDIR=$(DESTDIR)

# 编译proto
.PHONY: compile
compile:
	@echo "  >  "$@"ing $(DESTDIR)"
	$(MAKE) -C api/protobuf-spec/ compile

# 注释自动生成
.PHONY: comment
comment:
	@echo "  >  "$@"ing cmd"
	gocmt -i -d cmd
	@echo "  >  "$@"ing pkg"
	gocmt -i -d pkg
	@echo "  >  "$@"ing web"
	gocmt -i -d web

# 安装编译、代码生成等工具
.PHONY: tools
tools:
	@echo "  >  "$@"ing $(DESTDIR)"
	$(MAKE) -C api/protobuf-spec/ tools

# go mod tidy
## find . -type f -name "go.mod" -not -path "./.*" -exec bash -c 'cd $(dirname "$1"); go mod tidy' sh {} \;
# go test
## find . -type f -name "go.mod" -not -path "./.*" -exec bash -c 'cd $(dirname "$1");go test ./...' sh {} \;
# go mod tag
## TAG=v1.2.118 find . -type f -name "go.mod" -not -path "./.*" -not -path "./*/testdata/*" -exec bash -c 'path=$(dirname "${1#./}");if [ "$path" == "." ]; then git push origin :refs/"${TAG}"; else git push origin :refs/"${path}/${TAG}"; fi;' sh {} \;