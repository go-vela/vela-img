// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-vela/types/constants"
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
	// password for communication with the Docker Registry
	Password string
	// config path the docker json file exists for authentication
	Path string
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
			Name:     "config.registry",
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
		// nolint
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PATH", "REGISTRY_PATH", "DOCKER_CONFIG_PATH", "DOCKER_CONFIG"},
			FilePath: string("/vela/parameters/img/registry/path,/vela/secrets/img/registry/path,/vela/secrets/img/path"),
			Name:     "config.path",
			Usage:    "password for communication with the registry",
			Value:    "~/.docker/config.json",
		},
	}
)

// Write creates a Docker config.json file for building and publishing the image.
func (c *Config) Login() error {
	logrus.Trace("logging in registry information")

	// variable to store flags for command
	var flags []string

	// check if name, username and password are provided
	if len(c.URL) == 0 || len(c.Username) == 0 || len(c.Password) == 0 {
		return nil
	}

	flags = append(flags, fmt.Sprintf("-p=%s", c.Password))
	flags = append(flags, fmt.Sprintf("-u=%s", c.Username))
	flags = append(flags, c.URL)

	//nolint
	e := exec.Command(_img, append([]string{"login"}, flags...)...)

	// set command stdout to OS stdout
	e.Stdout = os.Stdout
	// set command stderr to OS stderr
	e.Stderr = os.Stderr

	cmd := strings.ReplaceAll(strings.Join(e.Args, " "), c.Password, constants.SecretMask)

	fmt.Println("$", cmd)

	return e.Run()
}

// Write creates a Docker config.json file for building and publishing the image.
//
// This function is not being used but keeping around incase the file approach
// becomes more viable.
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

	return a.WriteFile(c.Path, []byte(out), 0644)
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
