// Copyright 2020 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type visit struct {
	ip    string
	alias string
}

type server struct {
	db      *store
	cache   *lru
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
		cache:   newLRU(true),
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
	http.HandleFunc(conf.S.Prefix, s.shortHandler(kindShort))
	http.HandleFunc(conf.R.Prefix, s.shortHandler(kindRandom))
	// repo redirector
	http.Handle(conf.X.Prefix, s.xHandler(conf.X.VCS, conf.X.ImportPath, conf.X.RepoPath))
}

// backup tries to backup the data store to local files.
func (s *server) backup(ctx context.Context) {
	if _, err := os.Stat(conf.BackupDir); os.IsNotExist(err) {
		err := os.Mkdir(conf.BackupDir, os.ModePerm)
		if err != nil {
			log.Fatalf("cannot create backup directory, err: %v\n", err)
		}
	}

	t := time.NewTicker(time.Minute * time.Duration(conf.BackupMin))
	log.Printf("internal backup is running...")
	for {
		select {
		case <-t.C:
			r, err := s.db.Keys(ctx, "*")
			if err != nil {
				log.Printf("backup failure, err: %v\n", err)
				continue
			}
			if len(r) == 0 { // no keys for backup
				continue
			}

			d := make(map[string]interface{}, len(r))
			for _, k := range r {
				v, err := s.db.Fetch(ctx, k)
				if err != nil {
					log.Printf("backup failed because of key %v, err: %v\n", k, err)
					continue
				}
				var vv interface{}
				err = json.Unmarshal(StringToBytes(v), &vv)
				if err != nil {
					log.Printf("backup failed because unmarshal of key %v, err: %v\n", k, err)
					continue
				}
				d[k] = vv
			}

			b, err := yaml.Marshal(d)
			if err != nil {
				log.Printf("backup failed when converting to yaml, err: %v\n", err)
				continue
			}

			name := fmt.Sprintf("/backup-%s.yml", time.Now().Format(time.RFC3339))
			err = ioutil.WriteFile(conf.BackupDir+name, b, os.ModePerm)
			if err != nil {
				log.Printf("backup failed when saving the file, err: %v\n", err)
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}
