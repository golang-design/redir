// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"gopkg.in/yaml.v2"
)

// build info, assign by compile time or runtime.
var (
	Version   string
	BuildTime string
	GoVersion = runtime.Version()
)

type config struct {
	Host  string `json:"host"`
	Addr  string `json:"addr"`
	Store string `json:"store"`
	Log   string `json:"log"`
	S     struct {
		Prefix string `json:"prefix"`
	} `json:"s"`
	X struct {
		Prefix     string `json:"prefix"`
		VCS        string `json:"vcs"`
		ImportPath string `json:"import_path"`
		RepoPath   string `json:"repo_path"`
	} `json:"x"`
}

func (c *config) parse() {
	f := os.Getenv("REDIR_CONF")
	d, err := ioutil.ReadFile(f)
	if err != nil {
		// Just try again with default setting.
		d, err = ioutil.ReadFile("./config.yml")
		if err != nil {
			log.Fatalf("cannot read configuration, err: %v\n", err)
		}
	}
	err = yaml.Unmarshal(d, c)
	if err != nil {
		log.Fatalf("cannot parse configuration, err: %v\n", err)
	}
}

var conf config

func init() {
	conf.parse()
}
