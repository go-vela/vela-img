// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const buildAction = "build"

// Build represents the plugin configuration for build information.
type Build struct {
	// BuildArg should set build time variables
	BuildArgs []string
	// CacheFrom should be images to consider as cache sources
	CacheFrom []string
	// directory should be a path to the context you want img to run
	Directory string
	// File should be name and path to the Dockerfile
	File string
	// Labels should be set metadata for an image
	Labels []string
	// NoCache should be do not use cache when building the image
	NoCache bool
	// NoConole should be non-console progress UI
	NoConsole bool
	// Output BuildKit output specification (e.g. type=tar,dest=build.tar)
	Output string
	// Platform should be platforms for which the image should be built
	Platforms []string
	// Tag should be name and optionally a tag in the 'name:tag' format
	Tags []string
	// Target should be the target build stage to build
	Target string
}

// buildFlags represents for config settings on the cli.
var buildFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:     "build.build-args",
		Usage:    "should be images to consider as cache sources",
		EnvVars:  []string{"PARAMETER_BUILD_ARGS", "BUILD_BUILD_ARGS"},
		FilePath: string("/vela/parameters/img/build/build_args,/vela/secrets/img/build/build_args"),
	},
	&cli.StringSliceFlag{
		Name:     "build.cache-from",
		Usage:    "should set build time variables",
		EnvVars:  []string{"PARAMETER_CACHE_FROM", "BUILD_CACHE_FROM"},
		FilePath: string("/vela/parameters/img/build/cache_from,/vela/secrets/img/build/cache_from"),
	},
	&cli.StringFlag{
		Name:     "build.directory",
		Usage:    "should be a path to the context you want img to run",
		EnvVars:  []string{"PARAMETER_DIRECTORY", "BUILD_DIRECTORY"},
		FilePath: string("/vela/parameters/img/build/directory,/vela/secrets/img/build/directory"),
		Value:    ".",
	},
	&cli.StringFlag{
		Name:     "build.file",
		Usage:    "should be name and path to the Dockerfile",
		EnvVars:  []string{"PARAMETER_FILE", "BUILD_FILE"},
		FilePath: string("/vela/parameters/img/build/file,/vela/secrets/img/build/file"),
	},
	&cli.StringSliceFlag{
		Name:     "build.labels",
		Usage:    "should be set metadata for an image",
		EnvVars:  []string{"PARAMETER_LABELS", "BUILD_LABELS"},
		FilePath: string("/vela/parameters/img/build/labels,/vela/secrets/img/build/labels"),
	},
	&cli.BoolFlag{
		Name:     "build.no-cache",
		Usage:    "should be do not use cache when building the image",
		EnvVars:  []string{"PARAMETER_NO_CACHE", "BUILD_NO_CACHE"},
		FilePath: string("/vela/parameters/img/build/no_cache,/vela/secrets/img/build/no_cache"),
	},
	&cli.BoolFlag{
		Name:     "build.no-console",
		Usage:    "should be non-console progress UI",
		EnvVars:  []string{"PARAMETER_NO_CONSOLE", "BUILD_NO_CONSOLE"},
		FilePath: string("/vela/parameters/img/build/no_console,/vela/secrets/img/build/no_console"),
	},
	&cli.StringFlag{
		Name:     "build.output",
		Usage:    "BuildKit output specification",
		EnvVars:  []string{"PARAMETER_OUTPUT", "BUILD_OUTPUT"},
		FilePath: string("/vela/parameters/img/build/output,/vela/secrets/img/build/output"),
	},
	&cli.StringSliceFlag{
		Name:     "build.platforms",
		Usage:    "should be platforms for which the image should be built",
		EnvVars:  []string{"PARAMETER_PLATFORMS", "BUILD_PLATFORMS"},
		FilePath: string("/vela/parameters/img/build/platform,/vela/secrets/img/build/platform"),
	},
	&cli.StringSliceFlag{
		Name:     "build.tags",
		Usage:    "should be name and optionally a tag in the 'name:tag' format",
		EnvVars:  []string{"PARAMETER_TAGS", "BUILD_TAGS"},
		FilePath: string("/vela/parameters/img/build/tags,/vela/secrets/img/build/tags"),
	},
	&cli.StringFlag{
		Name:     "build.target",
		Usage:    "should be the target build stage to build",
		EnvVars:  []string{"PARAMETER_TARGET", "BUILD_TARGET"},
		FilePath: string("/vela/parameters/img/build/target,/vela/secrets/img/build/target"),
	},
}

// Command formats and outputs the Build command from
// the provided configuration to build a Docker image.
func (b *Build) Command() *exec.Cmd {
	logrus.Trace("creating img build command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if BuildArgs is provided
	if len(b.BuildArgs) > 0 {
		var args string
		for _, arg := range b.BuildArgs {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for BuildArgs from provided build command
		flags = append(flags, fmt.Sprintf("--build-arg \"%s\"", strings.TrimPrefix(args, " ")))
	}

	// check if CacheFrom is provided
	if len(b.CacheFrom) > 0 {
		var caches string
		for _, cache := range b.CacheFrom {
			caches += fmt.Sprintf(" %s", cache)
		}
		// add flag for CacheFrom from provided build command
		flags = append(flags, fmt.Sprintf("--cache-from \"%s\"", strings.TrimPrefix(caches, " ")))
	}

	// check if File is provided
	if len(b.File) > 0 {
		// add flag for File from provided build command
		flags = append(flags, fmt.Sprintf("-f=%s", b.File))
	}

	// check if Labels is provided
	if len(b.Labels) > 0 {
		var labels string
		for _, label := range b.Labels {
			labels += fmt.Sprintf(" %s", label)
		}
		// add flag for Labels from provided build command
		flags = append(flags, fmt.Sprintf("--label \"%s\"", strings.TrimPrefix(labels, " ")))
	}

	// check if NoCache is provided
	if b.NoCache {
		// add flag for NoCache from provided build command
		flags = append(flags, "--no-cache")
	}

	// check if NoConsole is provided
	if b.NoConsole {
		// add flag for NoConsole from provided build command
		flags = append(flags, "--no-console")
	}

	// check if Output is provided
	if len(b.Output) > 0 {
		// add flag for Output from provided build command
		flags = append(flags, fmt.Sprintf("--output %s", b.Output))
	}

	// check if Platforms is provided
	if len(b.Platforms) > 0 {
		var platforms string
		for _, platform := range b.Platforms {
			platforms += fmt.Sprintf(" %s", platform)
		}
		// add flag for Labels from provided build command
		flags = append(flags, fmt.Sprintf("--platform \"%s\"", strings.TrimPrefix(platforms, " ")))
	}

	// check if Tags is provided
	if len(b.Tags) > 0 {
		var tags string
		for _, tag := range b.Tags {
			tags += fmt.Sprintf(" %s", tag)
		}
		// add flag for Labels from provided build command
		flags = append(flags, fmt.Sprintf("-t=%s", strings.TrimPrefix(tags, " ")))
	}

	// check if Target is provided
	if len(b.Target) > 0 {
		// add flag for Target from provided build command
		flags = append(flags, fmt.Sprintf("--target %s", b.Target))
	}

	// add the required directory param
	flags = append(flags, b.Directory)

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.Command(_img, append([]string{buildAction}, flags...)...)
}

// Exec formats and runs the commands for building a Docker image.
func (b *Build) Exec() error {
	logrus.Trace("running build with provided configuration")

	// create the build command for the file
	cmd := b.Command()

	// run the build command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Build is properly configured.
func (b *Build) Validate() error {
	logrus.Trace("validating build plugin configuration")

	// verify directory are provided
	if len(b.Directory) == 0 {
		return fmt.Errorf("no build directory provided")
	}

	// verify tag are provided
	if len(b.Tags) == 0 {
		return fmt.Errorf("no build tag provided")
	}

	return nil
}
