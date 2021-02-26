# Go Tools

This subrepository holds the source for various packages and tools that support
the Go programming language.

All of the tools can be fetched with `go get`.

Packages include an implementation of the
Static Single Assignment form (SSA) representation for Go programs.

## Download/Install

The easiest way to install is to run `go get -u github.com/searKing/golang/tools/...`. You can
also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/tools/`.

### Tips
```bash
go get -u github.com/searKing/golang/tools/cmd/go-atomicvalue
go get -u github.com/searKing/golang/tools/cmd/go-enum
go get -u github.com/searKing/golang/tools/cmd/go-import
go get -u github.com/searKing/golang/tools/cmd/go-nulljson
go get -u github.com/searKing/golang/tools/cmd/go-option
go get -u github.com/searKing/golang/tools/cmd/go-sqlx
go get -u github.com/searKing/golang/tools/cmd/go-syncmap
go get -u github.com/searKing/golang/tools/cmd/go-syncpool
go get -u github.com/searKing/golang/tools/cmd/go-validator
go get -u github.com/searKing/golang/tools/cmd/protoc-gen-go-tag
```

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The main issue tracker for the tools repository is located at
https://github.com/searKing/golang/issues. Prefix your issue with "golang/tools/(your
subdir):" in the subject line, so it is easy to find.