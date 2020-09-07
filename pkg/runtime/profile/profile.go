// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package profile

import (
	"os"

	"github.com/pkg/profile"
)

type noop struct{}

// Stop is a noop.
func (p *noop) Stop() {}

// Profile parses the PROFILING environment variable and executes the proper profiling task.
func Profile() interface {
	Stop()
} {
	switch os.Getenv("PROFILING") {
	case "cpu":
		return profile.Start(profile.CPUProfile, profile.NoShutdownHook)
	case "mem":
		return profile.Start(profile.MemProfile, profile.NoShutdownHook)
	case "mutex":
		return profile.Start(profile.MutexProfile, profile.NoShutdownHook)
	case "block":
		return profile.Start(profile.BlockProfile, profile.NoShutdownHook)
	case "goroutine":
		return profile.Start(profile.GoroutineProfile, profile.NoShutdownHook)
	case "thread_creation":
		return profile.Start(profile.ThreadcreationProfile, profile.NoShutdownHook)
	case "trace":
		return profile.Start(profile.TraceProfile, profile.NoShutdownHook)
	}
	return new(noop)
}

// HelpMessage returns a string explaining how profiling works.
func HelpMessage() string {
	return `- PROFILING: Set "PROFILING=cpu" to enable cpu profiling and "PROFILING=mem" to enable memory profiling.
	It is not possible to do both at the same time. Profiling is disabled per default.
	Enables profiling if set. For more details on profiling, head over to: https://blog.golang.org/profiling-go-programs
	Support cpu|mem|mutex|block|goroutine|thread_creation|trace

	Example: PROFILING=cpu

	Set this value using environment variables on
	- Linux/macOS:
		$ export PROFILING=<value>
	- Windows Command Line (CMD):
		> set PROFILING=<value>
`
}
