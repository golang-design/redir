// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type visit struct {
	ip    string
	alias string
}

type server struct {
	db      *store // *store
	visitCh chan visit
}

var (
	xTmpl     *template.Template
	statsTmpl *template.Template
)

func newServer(ctx context.Context) *server {
	xTmpl = template.Must(template.ParseFiles("public/x.html"))
	statsTmpl = template.Must(template.ParseFiles("public/stats.html"))

	db, err := newStore(conf.Store)
	if err != nil {
		log.Fatalf("cannot establish connection to database, err: %v", err)
	}
	s := &server{
		db:      db,
		visitCh: make(chan visit, 100),
	}
	go s.counting(ctx)
	go s.backup(ctx)

	return s
}

func (s *server) close() {
	log.Println(s.db.Close())
}

func (s *server) registerHandler() {
	http.HandleFunc("/.info", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(struct {
			Version   string `json:"version"`
			GoVersion string `json:"go_version"`
			BuildTime string `json:"build_time"`
		}{
			Version:   Version,
			GoVersion: GoVersion,
			BuildTime: BuildTime,
		})
		w.Write(b)
	})

	// short redirector
	http.HandleFunc(conf.S.Prefix, s.sHandler)
	// repo redirector
	http.Handle(conf.X.Prefix, s.xHandler(conf.X.VCS, conf.X.ImportPath, conf.X.RepoPath))
}

// backup tries to backup the data store to local files every week.
// it will keeps the latest 10 backups of the data read from data store.
func (s *server) backup(ctx context.Context) {
	// TODO: do self-backups
}
