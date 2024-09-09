# note: call scripts from /scripts
MAKEFILE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PROJECT_ROOT_PATH := $(shell git rev-parse --show-toplevel)
PKG_CONFIG_ENV := ${MAKEFILE_PATH}/pkgconfig:${PKG_CONFIG_PATH}
PKG_CONFIG_PATH := ${MAKEFILE_PATH}/pkgconfig
target_name := $(shell basename ${MAKEFILE_PATH})

include ${PROJECT_ROOT_PATH}/build/scripts/build_service.mk