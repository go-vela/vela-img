// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-vela/vela-img/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-img"
	app.HelpName = "vela-img"
	app.Usage = "Vela img plugin for building and publishing images"
	app.Copyright = "Copyright (c) 2021 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Plugin Flags

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "VELA_LOG_LEVEL", "IMG_LOG_LEVEL"},
			FilePath: string("/vela/parameters/img/log_level,/vela/secrets/img/log_level"),
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},
	}

	// add config flags
	app.Flags = append(app.Flags, configFlags...)

	// add build flags
	app.Flags = append(app.Flags, buildFlags...)

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
	// set the log level for the plugin
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-img",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/img",
		"registry": "https://hub.docker.com/r/target/vela-img",
	}).Info("Vela Img Plugin")

	// create the plugin
	p := Plugin{
		Config: &Config{
			Password: c.String("config.password"),
			URL:      c.String("config.registry"),
			Username: c.String("config.username"),
		},
		Build: &Build{
			BuildArgs: c.StringSlice("build.build-args"),
			CacheFrom: c.StringSlice("build.cache-from"),
			Directory: c.String("build.directory"),
			File:      c.String("build.file"),
			Labels:    c.StringSlice("build.labels"),
			NoCache:   c.Bool("build.no-cache"),
			NoConsole: c.Bool("build.no-console"),
			Output:    c.String("build.output"),
			Platforms: c.StringSlice("build.platforms"),
			Tags:      c.StringSlice("build.tags"),
			Target:    c.String("build.target"),
		},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
