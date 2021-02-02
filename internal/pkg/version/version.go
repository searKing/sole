package version

import (
	"fmt"
	"runtime"
)

var (
	// NOTE: The $Format strings are replaced during 'git archive' thanks to the
	// companion .gitattributes file containing 'export-subst' in this same
	// directory.  See also https://git-scm.com/docs/gitattributes
	Version         = "v0.0.0-master+$Format:%h$" // git describe --long --tags --dirty --tags --always
	BuildTime       = "1970-01-01T00:00:00Z"      // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	GitHash         = "$Format:%H$"               // sha1 from git, output of $(git rev-parse HEAD)
	GoVersion       = runtime.Version()
	Compiler        = runtime.Compiler
	Platform        = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

const (
	ServiceName        = "sole"
	ServiceDescription = "sole is a cloud native high throughput service manager server, " +
		"allowing you to manage all services."
)
