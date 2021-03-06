// Code generated by "go-option -type compileQuery"; DO NOT EDIT.

package sql

var _default_compileQuery_value = func() (val compileQuery) { return }()

// A CompileQueryOption sets options.
type CompileQueryOption interface {
	apply(*compileQuery)
}

// EmptyCompileQueryOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyCompileQueryOption struct{}

func (EmptyCompileQueryOption) apply(*compileQuery) {}

// CompileQueryOptionFunc wraps a function that modifies compileQuery into an
// implementation of the CompileQueryOption interface.
type CompileQueryOptionFunc func(*compileQuery)

func (f CompileQueryOptionFunc) apply(do *compileQuery) {
	f(do)
}

// sample code for option, default for nothing to change
func _CompileQueryOptionWithDefault() CompileQueryOption {
	return CompileQueryOptionFunc(func(*compileQuery) {
		// nothing to change
	})
}

func (o *compileQuery) ApplyOptions(options ...CompileQueryOption) *compileQuery {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
