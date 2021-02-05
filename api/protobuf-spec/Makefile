# note: call scripts from /scripts

GO := go

.PHONY: all
all: tools generate

.PHONY: tools
tools:
	@echo "  >  installing tools"
	@go generate --tags=tools tools.go

.PHONY: generate
generate:
	@echo "  >  compiling protocol buffers"
	@go generate --tags=tools proto-gen.go
