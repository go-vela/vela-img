// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestImg_Plugin_Exec(t *testing.T) {
	// TODO Write test
}

func TestImg_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			BuildArgs: []string{"FOO"},
			CacheFrom: []string{"index.docker.io/target/vela-img"},
			Directory: ".",
			File:      "Dockerfile",
			Labels:    []string{"sha"},
			NoCache:   true,
			NoConsole: true,
			Output:    "type=tar,dest=build.tar",
			Platforms: []string{"linux/amd64"},
			Tags:      []string{"latest"},
			Target:    "foo",
		},
		Config: &Config{
			Password: "superSecretPassword",
			URL:      "index.docker.io",
			Username: "octocat",
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}
