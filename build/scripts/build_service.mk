# note: call scripts from /scripts
#MAKEFILE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
#PROJECT_ROOT_PATH := $(shell git rev-parse --show-toplevel)
#PKG_CONFIG_ENV := ${MAKEFILE_PATH}/pkgconfig:${PKG_CONFIG_PATH}
#PKG_CONFIG_PATH := ${MAKEFILE_PATH}/pkgconfig
#target_name := $(shell basename ${MAKEFILE_PATH})

# dynamically load lib*.so

# A literal space.
space :=
space +=
# Joins elements of the list in arg 2 with the given separator.
#   1. Element separator.
#   2. The list.
join-with = $(subst $(space),$1,$(strip $2))

git_tag=$(shell git describe --long --tags --dirty --tags --always)
git_commit=$(shell git rev-parse HEAD)
git_build_time=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
git_root=$(shell git rev-parse --show-toplevel)
PROJECT_ROOT_PATH ?= ${git_root}
$(info ************ PROJECT_ROOT_PATH=${PROJECT_ROOT_PATH} ************)

GO := go
GO_DOWNLOAD_FLAGS :=
GO_BUILD_FLAGS := -a -gcflags=all="-N -l" \
-ldflags "-X 'github.com/searKing/golang/go/version.GitTag=${git_tag}' \
-X 'github.com/searKing/golang/go/version.BuildTime=${git_build_time}' \
-X 'github.com/searKing/golang/go/version.GitHash=${git_commit}'"
CONF_NAME ?=
DEPS_NAME ?= deps.yaml
DO_PARSE_DEPS ?= ON
GO_BUILD_TAG ?=
build_tags := ${GO_BUILD_TAG}
ifeq ($(ENABLE_PPROF),ON)
	ifneq ($(build_tags),)
		build_tags :=$(build_tags),
	endif
	build_tags :=$(build_tags)enable_pprof
endif

ifeq ($(WITH_LICENSE),ON)
	ifneq ($(build_tags),)
		build_tags :=$(build_tags),
	endif
	build_tags :=$(build_tags)license_force_enabled
endif

ifneq ($(build_tags),)
	GO_BUILD_FLAGS :=$(GO_BUILD_FLAGS) -tags $(build_tags)
endif

.PHONY: all
all: build

.IGNORE: clean
clean:
	@echo "  >  Cleaning build cache"
	@${GO} clean -cache -testcache
	@if [ -f ${target_name} ] ; then rm ${target_name} ; fi

ifeq (,$(wildcard $(shell echo $(DEPS_NAME)) ))
$(info ************ GO PROJECT ************)
.PHONY: build
build: go-get clean git-info build-only

.PHONY: build-only
build-only: $(target_name)

.PHONY: $(target_name)
$(target_name):
	@echo "  >  Building binary $(target_name) ..."
	${GO} build ${GO_BUILD_FLAGS} -o $@ .

.PHONY: go-get
go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	${GO} get ${GO_DOWNLOAD_FLAGS}
	@${GO} mod tidy

.PHONY: install
install: build install-only-go
else
$(info ************ CGO PROJECT ************)
.PHONY: build
build: go-get clean git-info pkg-config build-only

.PHONY: build-only
build-only: $(target_name)

.PHONY: $(target_name)
$(target_name):
	@echo "  >  Building binary $(target_name) ..."
	@$(eval COMPILE_THIRD_LIB_PATHS := $(shell find -L third_path/ -maxdepth 3 -mindepth 2 -type d \( -iname "lib*" -o -iname "stubs" \) -print0 |xargs -0 -I {} sh -c 'echo {}'))
	@$(eval COMPILE_JOINED_THIRD_LIB_PATHS := $(call join-with,:,$(COMPILE_THIRD_LIB_PATHS)))
	_GLIBCXX_USE_CXX11_ABI=0 PKG_CONFIG_PATH="${PKG_CONFIG_ENV}" LD_LIBRARY_PATH="$(COMPILE_JOINED_THIRD_LIB_PATHS):${LD_LIBRARY_PATH}" LIBRARY_PATH="$(COMPILE_JOINED_THIRD_LIB_PATHS):${LIBRARY_PATH}" ${GO} build ${GO_BUILD_FLAGS} -o $@ .

.PHONY: go-get
go-get: pkg-config
	@echo "  >  Checking if there is any missing dependencies..."
	@$(eval COMPILE_THIRD_LIB_PATHS := $(shell find -L third_path/ -maxdepth 3 -mindepth 2 -type d \( -iname "lib*" -o -iname "stubs" \) -print0 |xargs -0 -I {} sh -c 'echo {}'))
	@$(eval COMPILE_JOINED_THIRD_LIB_PATHS := $(call join-with,:,$(COMPILE_THIRD_LIB_PATHS)))
	@PKG_CONFIG_PATH="${PKG_CONFIG_ENV}" LD_LIBRARY_PATH="$(COMPILE_JOINED_THIRD_LIB_PATHS):${LD_LIBRARY_PATH}" LIBRARY_PATH="$(COMPILE_JOINED_THIRD_LIB_PATHS):${LIBRARY_PATH}"  ${GO} mod tidy

.PHONY: install
install: build install-only-go install-only-cgo
endif

PREFIX :=
ifeq ($(DESTDIR),)
PREFIX=.
endif

.PHONY: conf.yaml
conf.yaml:
	@echo "  >  downloading ${CONF_NAME} as $(target_name).yaml"
	@if [[ (( "${CONF_NAME}"x != ""x )) && (( "${CONF_NAME}"x != "$(target_name).yaml"x )) ]]; then echo "  >  switch ${CONF_NAME} to $(target_name).yaml"; cd "${MAKEFILE_PATH}/../../conf";ln -sf "${CONF_NAME}" $(target_name).yaml; fi

.PHONY: install-only-go
install-only-go: conf.yaml
	@echo "  >  installing target to $(DESTDIR)/$(PREFIX)"
	@mkdir -p "$(DESTDIR)/$(PREFIX)/bin"
	@cp -dRv "$(target_name)" "$(DESTDIR)/$(PREFIX)/bin/$(target_name)"
	@mkdir -p "$(DESTDIR)/$(PREFIX)/conf"
	@if [[ -f "${MAKEFILE_PATH}/../../conf/$(target_name).yaml" ]]; then cp -dRv -H "${MAKEFILE_PATH}/../../conf/$(target_name).yaml" "$(DESTDIR)/$(PREFIX)/conf/"; fi
	@mkdir -p "$(DESTDIR)/$(PREFIX)/test"
	@if [[ -d "${MAKEFILE_PATH}/../../test/$(target_name)" ]]; then cp -Rdv "${MAKEFILE_PATH}/../../test/$(target_name)/" "$(DESTDIR)/$(PREFIX)/test/"; fi
	@if [[ -f "${MAKEFILE_PATH}/Dockerfile" ]]; then cp -vL "${MAKEFILE_PATH}/Dockerfile" "$(DESTDIR)/$(PREFIX)/"; fi

.PHONY: install-only-cgo
install-only-cgo:
	@echo "  >  installing libs to $(DESTDIR)/$(PREFIX)/lib"
	@mkdir -p "$(DESTDIR)/$(PREFIX)/lib/"
	@$(eval LINK_THIRD_LIB_PATHS := $(shell find -L third_path/ -maxdepth 3 -mindepth 2 -type d -iname "lib*" -print0 |xargs -0 -I {} sh -c 'echo {}'|grep -v "stubs"))
	@$(eval JOINED_LINK_THIRD_LIB_PATHS := $(call join-with,:,$(LINK_THIRD_LIB_PATHS)))
	@LD_LIBRARY_PATH="$(JOINED_LINK_THIRD_LIB_PATHS):${LD_LIBRARY_PATH}" ldd "$(target_name)" | awk '{if (match($$3,"/")){ print $$3  }}' |grep "third_path" | grep -v "^/lib64" | grep -v "^/lib" | xargs -I {} sh -c 'cp -v -d -L {} $(DESTDIR)/$(PREFIX)/lib/'
	@echo "  >  installing third_path to DESTDIR=$(DESTDIR)"
	@find -L third_path/ -maxdepth 3 -type d \( -iname "model" -o -iname "sdk_data" -o -iname "config" \) -print0 | xargs -0 -I {} sh -c 'mkdir -p $(DESTDIR)/$(PREFIX)/{}; cp -r -v -d {}/* $(DESTDIR)/$(PREFIX)/{}'

.PHONY: uninstall
uninstall:
	@echo "  >  uninstalling $(DESTDIR)/$(PREFIX)/"
	@if [[ -f "$(DESTDIR)/$(PREFIX)/bin/$(target_name)" ]]; then rm -Rv "$(DESTDIR)/$(PREFIX)/bin/$(target_name)"; fi
	@if [[ -d "$(DESTDIR)/$(PREFIX)/conf" ]]; then rm -Rv "$(DESTDIR)/$(PREFIX)/conf"; fi
	@if [[ -d "$(DESTDIR)/$(PREFIX)/lib" ]]; then rm -Rv "$(DESTDIR)/$(PREFIX)/lib"; fi
	@if [[ -d "$(DESTDIR)/$(PREFIX)/third_path" ]]; then rm -Rv "$(DESTDIR)/$(PREFIX)/third_path"; fi

.PHONY: git-info
git-info:
	@echo "  >  git_tag $(git_tag) ..."
	@echo "  >  git_commit $(git_commit) ..."
	@echo "  >  git_build_time $(git_build_time) ..."

.PHONY: deps.yaml
deps.yaml:
	@echo "  >  downloading ${DEPS_NAME} as deps.yaml"
	@if [[ (( "${DEPS_NAME}"x != ""x )) && (( "${DEPS_NAME}"x != "deps.yaml"x )) ]]; then echo "  >  switch ${DEPS_NAME} to deps.yaml"; ln -sf "${DEPS_NAME}" deps.yaml; fi
	@cd ${PROJECT_ROOT_PATH}; git clone -b v5.0 --single-branch http://git.woa.com/youtu_sdk/yt-sdk-qci-scripts.git || true
	@if [[ (( "$(DO_PARSE_DEPS)"x = "ON"x )) && (( -f ./deps.yaml )) ]]; then python ${PROJECT_ROOT_PATH}/yt-sdk-qci-scripts/deps/parse_deps.py ./deps.yaml; fi

.PHONY: pkg-config
pkg-config: deps.yaml
	@echo "  >  copying pkg-config ..."
	@mkdir -p "${PKG_CONFIG_PATH}"
	@echo "PKG: ${PKG_CONFIG_PATH}"
	@echo "MAKEFILE_PATH: ${MAKEFILE_PATH}"
	if [[ (( -d "${MAKEFILE_PATH}"/third_path/yt-sdk-go/pkgconfig )) ]]; then cp -df "${MAKEFILE_PATH}"/third_path/yt-sdk-go/pkgconfig/* "${PKG_CONFIG_PATH}/"; fi
	@echo "  >  switching opencv version by deps.yaml"
	@if [[ -L ${PKG_CONFIG_PATH}/opencv.pc ]]; then echo "  >  delete link: opencv.pc"; rm ${PKG_CONFIG_PATH}/opencv.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/opencv4.pc )) && (( -d "${MAKEFILE_PATH}/third_path/opencv/include/opencv2/imgcodecs" )) ]]; then echo "  >  switch opencv version to opencv4 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf opencv4.pc opencv.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/opencv2.pc )) && (( -d "${MAKEFILE_PATH}/third_path/opencv/include/opencv2/contrib" )) ]]; then echo "  >  switch opencv version to opencv2 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf opencv2.pc opencv.pc; fi
	@if [[ ! -f ${PKG_CONFIG_PATH}/opencv.pc ]]; then echo "  >  switch opencv version to empty by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf empty.pc opencv.pc; fi
	@echo "  >  switching python version"
	@if [[ -L ${PKG_CONFIG_PATH}/python3-embed.pc ]]; then echo "  >  delete link: python3-embed.pc"; rm ${PKG_CONFIG_PATH}/python3-embed.pc; fi
	@if [[ (( -f "/opt/rh/rh-python38/root/usr/lib64/pkgconfig/python3-embed.pc" )) && (( -d "${MAKEFILE_PATH}/third_path/pybind11" )) ]]; then echo "  >  switch python version to python-rh"; cd ${PKG_CONFIG_PATH}/; ln -sf "/opt/rh/rh-python38/root/usr/lib64/pkgconfig/python3-embed.pc" python3-embed.pc; fi
	@if [[ (( -f "/usr/lib64/pkgconfig/python-3.8-embed.pc" )) && (( -d "${MAKEFILE_PATH}/third_path/pybind11" )) ]]; then echo "  >  switch python version to python-system"; cd ${PKG_CONFIG_PATH}/; ln -sf "/usr/lib64/pkgconfig/python-3.8-embed.pc" python3-embed.pc; fi
	@if [[ (( -f "/usr/lib64/pkgconfig/python3-embed.pc" )) && (( -d "${MAKEFILE_PATH}/third_path/pybind11" )) ]]; then echo "  >  switch python version to python-system"; cd ${PKG_CONFIG_PATH}/; ln -sf "/usr/lib64/pkgconfig/python3-embed.pc" python3-embed.pc; fi
	@if [[ (( -f "/usr/local/lib/pkgconfig/python-3.10-embed.pc" )) && (( -d "${MAKEFILE_PATH}/third_path/pybind11" )) ]]; then echo "  >  switch python version to python-3.10.12"; cd ${PKG_CONFIG_PATH}/; ln -sf "/usr/local/lib/pkgconfig/python-3.10-embed.pc" python3-embed.pc; fi
	@if [[ ! -f ${PKG_CONFIG_PATH}/python3-embed.pc ]]; then echo "  >  switch python version to empty as not installed, yum install python3-devel please"; cd ${PKG_CONFIG_PATH}/; ln -sf empty.pc python3-embed.pc; fi
	@echo "  >  switching ImageMagick version by deps.yaml"
	@if [[ -L ${PKG_CONFIG_PATH}/ImageMagick.pc ]]; then echo "  >  delete link: ImageMagick.pc"; rm ${PKG_CONFIG_PATH}/ImageMagick.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/ImageMagick-7.Q16.pc )) && (( -f "${MAKEFILE_PATH}/third_path/ImageMagick/lib/libMagickCore-7.Q16.so" )) ]]; then echo "  >  switch ImageMagick version to MagickCore-7.Q16 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf ImageMagick-7.Q16.pc ImageMagick.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/ImageMagick-7.Q16HDRI.pc )) && (( -f "${MAKEFILE_PATH}/third_path/ImageMagick/lib/libMagickCore-7.Q16HDRI.so" )) ]]; then echo "  >  switch ImageMagick version to MagickCore-7.Q16HDRI by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf ImageMagick-7.Q16HDRI.pc ImageMagick.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/ImageMagick-6.Q16.pc )) && (( -f "${MAKEFILE_PATH}/third_path/ImageMagick/lib/libMagickCore-6.Q16.so" )) ]]; then echo "  >  switch ImageMagick version to MagickCore-6.Q16 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf ImageMagick-6.Q16.pc ImageMagick.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/ImageMagick-6.Q16HDRI.pc )) && (( -f "${MAKEFILE_PATH}/third_path/ImageMagick/lib/libMagickCore-6.Q16HDRI.so" )) ]]; then echo "  >  switch ImageMagick version to MagickCore-6.Q16HDRI by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf ImageMagick-6.Q16HDRI.pc ImageMagick.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/GraphicsMagick.pc )) && (( -f "${MAKEFILE_PATH}/third_path/GraphicsMagick/lib/libGraphicsMagick.so" )) ]]; then echo "  >  switch ImageMagick version to GraphicsMagick by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf GraphicsMagick.pc ImageMagick.pc; fi
	@if [[ ! -f ${PKG_CONFIG_PATH}/ImageMagick.pc ]]; then echo "  >  switch ImageMagick version to empty by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf empty.pc ImageMagick.pc; fi
	@echo "  >  switching cuda version by deps.yaml"
	@if [[ -L ${PKG_CONFIG_PATH}/cuda.pc ]]; then echo "  >  delete link: cuda.pc"; rm ${PKG_CONFIG_PATH}/cuda.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/cuda10_0.pc )) && (( -d "${MAKEFILE_PATH}/third_path/cuda10_0/include" )) ]]; then echo "  >  switch cuda version to cuda10_0 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf cuda10_0.pc cuda.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/cuda10_2.pc )) && (( -d "${MAKEFILE_PATH}/third_path/cuda10_2/include" )) ]]; then echo "  >  switch cuda version to cuda10_2 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf cuda10_2.pc cuda.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/cuda11_1.pc )) && (( -d "${MAKEFILE_PATH}/third_path/cuda11_1/include" )) ]]; then echo "  >  switch cuda version to cuda11_1 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf cuda11_1.pc cuda.pc; fi
	@if [[ ! -f ${PKG_CONFIG_PATH}/cuda.pc ]]; then echo "  >  switch cuda version to empty by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf empty.pc cuda.pc; fi
	@echo "  >  switching cudnn version by deps.yaml"
	@if [[ -L ${PKG_CONFIG_PATH}/cudnn.pc ]]; then echo "  >  delete link: cudnn.pc"; rm ${PKG_CONFIG_PATH}/cudnn.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/cudnn_cuda10_0.pc )) && (( -d "${MAKEFILE_PATH}/third_path/cudnn_cuda10_0/include" )) ]]; then echo "  >  switch cudnn version to cudnn_cuda10_0 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf cudnn_cuda10_0.pc cudnn.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/cudnn_cuda10_2.pc )) && (( -d "${MAKEFILE_PATH}/third_path/cudnn_cuda10_2/include" )) ]]; then echo "  >  switch cudnn version to cudnn_cuda10_2 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf cudnn_cuda10_2.pc cudnn.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/cudnn_cuda11_1.pc )) && (( -d "${MAKEFILE_PATH}/third_path/cudnn_cuda11_1/include" )) ]]; then echo "  >  switch cudnn version to cudnn_cuda11_1 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf cudnn_cuda11_1.pc cudnn.pc; fi
	@if [[ ! -f ${PKG_CONFIG_PATH}/cudnn.pc ]]; then echo "  >  switch cudnn version to empty by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf empty.pc cudnn.pc; fi
	@echo "  >  switching rapidnet_gpu version by deps.yaml"
	@if [[ -L ${PKG_CONFIG_PATH}/rapidnet_gpu.pc ]]; then echo "  >  delete link: rapidnet_gpu.pc"; rm ${PKG_CONFIG_PATH}/rapidnet_gpu.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/rapidnet_gpu_cuda10_0.pc )) && (( -d "${MAKEFILE_PATH}/third_path/rapidnet_gpu_cuda10_0/include" )) ]]; then echo "  >  switch rapidnet_gpu version to rapidnet_gpu_cuda10_0 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf rapidnet_gpu_cuda10_0.pc rapidnet_gpu.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/rapidnet_gpu_cuda10_2.pc )) && (( -d "${MAKEFILE_PATH}/third_path/rapidnet_gpu_cuda10_2/include" )) ]]; then echo "  >  switch rapidnet_gpu version to rapidnet_gpu_cuda10_2 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf rapidnet_gpu_cuda10_2.pc rapidnet_gpu.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/rapidnet_gpu_cuda11_1.pc )) && (( -d "${MAKEFILE_PATH}/third_path/rapidnet_gpu_cuda11_1/include" )) ]]; then echo "  >  switch rapidnet_gpu version to rapidnet_gpu_cuda11_1 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf rapidnet_gpu_cuda11_1.pc rapidnet_gpu.pc; fi
	@if [[ ! -f ${PKG_CONFIG_PATH}/rapidnet_gpu.pc ]]; then echo "  >  switch rapidnet_gpu version to empty by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf empty.pc rapidnet_gpu.pc; fi
	@echo "  >  switching tensorRT version by deps.yaml"
	@if [[ -L ${PKG_CONFIG_PATH}/tensorRT.pc ]]; then echo "  >  delete link: tensorRT.pc"; rm ${PKG_CONFIG_PATH}/tensorRT.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/tensorRT_cuda10_0.pc )) && (( -d "${MAKEFILE_PATH}/third_path/tensorRT_cuda10_0/include" )) ]]; then echo "  >  switch tensorRT version to tensorRT_cuda10_0 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf tensorRT_cuda10_0.pc tensorRT.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/tensorRT_cuda10_2.pc )) && (( -d "${MAKEFILE_PATH}/third_path/tensorRT_cuda10_2/include" )) ]]; then echo "  >  switch tensorRT version to tensorRT_cuda10_2 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf tensorRT_cuda10_2.pc tensorRT.pc; fi
	@if [[ (( -f ${PKG_CONFIG_PATH}/tensorRT_cuda11_1.pc )) && (( -d "${MAKEFILE_PATH}/third_path/tensorRT_cuda11_1/include" )) ]]; then echo "  >  switch tensorRT version to tensorRT_cuda11_1 by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf tensorRT_cuda11_1.pc tensorRT.pc; fi
	@if [[ ! -f ${PKG_CONFIG_PATH}/tensorRT.pc ]]; then echo "  >  switch tensorRT version to empty by deps.yaml"; cd ${PKG_CONFIG_PATH}/; ln -sf empty.pc tensorRT.pc; fi


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
