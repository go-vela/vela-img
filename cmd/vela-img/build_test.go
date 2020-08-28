// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"
)

func TestImg_Build_Command(t *testing.T) {
	// setup types
	b := &Build{
		BuildArgs: []string{"FOO"},
		CacheFrom: []string{"index.docker.io/target/vela-img"},
		File:      "Dockerfile",
		Labels:    []string{"sha"},
		NoCache:   true,
		NoConsole: true,
		Output:    "type=tar,dest=build.tar",
		Platforms: []string{"linux/amd64"},
		Tags:      []string{"image_name:tag"},
		Target:    "foo",
	}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	want := exec.Command(
		_img,
		buildAction,
		fmt.Sprintf("--build-arg \"%s\"", b.BuildArgs[0]),
		fmt.Sprintf("--cache-from \"%s\"", b.CacheFrom[0]),
		fmt.Sprintf("--file %s", b.File),
		fmt.Sprintf("--label \"%s\"", b.Labels[0]),
		"--no-cache",
		"--no-console",
		fmt.Sprintf("--output %s", b.Output),
		fmt.Sprintf("--platform \"%s\"", b.Platforms[0]),
		fmt.Sprintf("--tag \"%s\"", b.Tags[0]),
		fmt.Sprintf("--target %s", b.Target),
	)

	got := b.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestImg_Build_Exec_Error(t *testing.T) {
	// setup types
	b := &Build{}

	err := b.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestImg_Build_Validate(t *testing.T) {
	// setup types
	b := &Build{
		Tags: []string{"image_name:tag"},
	}

	err := b.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestImg_Config_Validate_NoTags(t *testing.T) {
	// setup types
	b := &Build{}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
