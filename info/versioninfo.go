/*
 *  Copyright 2022 Ashok Pon Kumar
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package info

import (
	"runtime"
)

var (
	version       = "v0.1.0"
	buildmetadata = ""
	gitCommit     = ""
	gitTreeState  = ""
)

// GetVersion returns the semver string of the version
func GetVersion() string {
	if buildmetadata == "" {
		return version
	}
	return version + "+" + buildmetadata
}

// GetVersionInfo returns version info
func GetVersionInfo() VersionInfo {
	v := VersionInfo{
		Version:      GetVersion(),
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
	}
	return v
}

// VersionInfo describes the compile time information.
type VersionInfo struct {
	// Version is the current semver.
	Version string `yaml:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `yaml:"gitCommit,omitempty"`
	// GitTreeState is the state of the git tree.
	GitTreeState string `yaml:"gitTreeState,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `yaml:"goVersion,omitempty"`
}
