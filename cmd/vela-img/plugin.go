// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// build arguments loaded for the plugin
	Build *Build
	// config arguments loaded for the plugin
	Config *Config
}

// Exec formats and runs the commands for building and publishing a Docker image.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	// output img version for troubleshooting
	err := execCmd(versionCmd())
	if err != nil {
		return err
	}

	// write the config.json file with Docker credentials
	err = p.Config.Login()
	if err != nil {
		return err
	}

	// execute build action
	return p.Build.Exec()
}

// Validate verifies the Plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate config configuration
	err := p.Config.Validate()
	if err != nil {
		return err
	}

	// validate build configuration
	err = p.Build.Validate()
	if err != nil {
		return err
	}

	return nil
}
