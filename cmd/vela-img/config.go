// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/base64"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

const (
	credentials = `%s:%s`

	registryFile = `{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`
)

// Config holds input parameters for the plugin.
type Config struct {
	// Password for communication with the Docker Registry
	Password string
	// full url to Docker Registry
	URL string
	// user name for communication with the Docker Registry
	Username string
}

var (
	appFS = afero.NewOsFs()

	// configFlags represents for config settings on the cli.
	configFlags = []cli.Flag{
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REGISTRY", "REGISTRY_NAME"},
			FilePath: string("/vela/parameters/img/registry/name,/vela/secrets/docker/registry/name"),
			Name:     "config.name",
			Usage:    "Docker registry name to communicate with",
			Value:    "index.docker.io",
		},
		// nolint
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "REGISTRY_USERNAME", "DOCKER_USERNAME"},
			FilePath: string("/vela/parameters/img/registry/username,/vela/secrets/img/registry/username,/vela/secrets/img/username"),
			Name:     "config.username",
			Usage:    "user name for communication with the registry",
		},
		// nolint
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "REGISTRY_PASSWORD", "DOCKER_PASSWORD"},
			FilePath: string("/vela/parameters/img/registry/password,/vela/secrets/img/registry/password,/vela/secrets/img/password"),
			Name:     "config.password",
			Usage:    "password for communication with the registry",
		},
	}
)

// Write creates a Docker config.json file for building and publishing the image.
func (c *Config) Write() error {
	logrus.Trace("writing registry configuration file")

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	// check if name, username and password are provided
	if len(c.URL) == 0 || len(c.Username) == 0 || len(c.Password) == 0 {
		return nil
	}

	// create basic authentication string for config.json file
	basicAuth := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(credentials, c.Username, c.Password)),
	)

	// create output string for config.json file
	out := fmt.Sprintf(
		registryFile,
		c.URL,
		basicAuth,
	)

	// create full path for config.json file
	path := "/root/.docker/config.json"

	return a.WriteFile(path, []byte(out), 0644)
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config plugin configuration")

	// verify password are provided
	if len(c.Password) == 0 {
		return fmt.Errorf("no config password provided")
	}

	// verify url is provided
	if len(c.URL) == 0 {
		return fmt.Errorf("no config url provided")
	}

	// verify username is provided
	if len(c.Username) == 0 {
		return fmt.Errorf("no config username provided")
	}

	return nil
}
