// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
)

// op is a short link operator
type op string

const (
	// opCreate represents a create operation for short link
	opCreate op = "create"
	// opDelete represents a delete operation for short link
	opDelete = "delete"
	// opUpdate represents a update operation for short link
	opUpdate = "update"
	// opFetch represents a fetch operation for short link
	opFetch = "fetch"
)

func (o op) valid() bool {
	switch o {
	case opCreate, opDelete, opUpdate, opFetch:
		return true
	default:
		return false
	}
}

func shortFile(fname string) {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("cannot read import file, err: %v\n", err)
	}

	d := map[string]string{}
	err = yaml.Unmarshal(b, &d)
	if err != nil {
		log.Fatalf("cannot unmarshal the imported file, err: %v\n", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	for alias, link := range d {
		err = shortCmd(ctx, opUpdate, alias, link)
		if err != nil {
			err = shortCmd(ctx, opCreate, alias, link)
			if err != nil {
				log.Fatalf("cannot import alias %v, err: %v\n", alias, err)
			}
		}
	}
	return
}

// shortCmd processes the given alias and link with a specified op.
func shortCmd(ctx context.Context, operate op, alias, link string) (retErr error) {
	s, err := newStore(conf.Store)
	if err != nil {
		return fmt.Errorf("cannot create a new alias, err: %w", err)
	}
	defer s.Close()

	errf := func(o op, err error) {
		retErr = fmt.Errorf("cannot %v alias to data store, err: %w", o, err)
	}
	switch operate {
	case opCreate:
		err = s.StoreAlias(ctx, alias, link)
		if err != nil {
			errf(opCreate, err)
			return
		}
		log.Printf("alias %v has been created.\n", alias)
	case opUpdate:
		err = s.UpdateAlias(ctx, alias, link)
		if err != nil {
			errf(opUpdate, err)
			return
		}
		log.Printf("alias %v has been updated.\n", alias)
	case opDelete:
		err = s.DeleteAlias(ctx, alias)
		if err != nil {
			errf(opDelete, err)
			return
		}
		log.Printf("alias %v has been deleted.\n", alias)
	case opFetch:
		r, err := s.FetchAlias(ctx, alias)
		if err != nil {
			errf(opFetch, err)
			return
		}
		log.Println(r)
	}
	return
}

// sHandler redirects ...
func (s *server) sHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		if err != nil {
			// Just tell the user we could not find the record rather than
			// throw 50x. The server should be able to identify the issue.
			log.Printf("stats err: %v\n", err)
			// Use 307 redirect to 404 page
			http.Redirect(w, r, conf.Host+"/404.html", http.StatusTemporaryRedirect)
		}
	}()

	alias := strings.TrimLeft(r.URL.Path, conf.S.Prefix)
	if alias == "" {
		err = s.stats(ctx, w)
		return
	}

	checkdb := func(url string) (string, error) {
		raw, err := s.db.FetchAlias(ctx, alias)
		if err != nil {
			return "", err
		}
		c := arecord{}
		err = json.Unmarshal(StringToBytes(raw), &c)
		if err != nil {
			return "", err
		}
		if url != c.URL {
			s.cache.Put(alias, c.URL)
			url = c.URL
		}
		return url, nil
	}

	url, ok := s.cache.Get(alias)
	if !ok {
		url, err = checkdb(url)
		if err != nil {
			return
		}
	}

	// redirect the user immediate, but run pv/uv count in background
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	// count visit in another goroutine so it won't block the redirect.
	go func() { s.visitCh <- visit{s.readIP(r), alias} }()
}

// readIP implements a best effort approach to return the real client IP,
// it parses X-Real-IP and X-Forwarded-For in order to work properly with
// reverse-proxies such us: nginx or haproxy. Use X-Forwarded-For before
// X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
//
// This implementation is derived from gin-gonic/gin.
func (s *server) readIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}
	if addr := r.Header.Get("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}
	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return "unknown" // use unknown to guarantee non empty string
	}
	return ip
}

type arecords struct {
	Records []arecord
}

func (s *server) stats(ctx context.Context, w http.ResponseWriter) (retErr error) {
	aliases, retErr := s.db.Keys(ctx, prefixalias+"*")
	if retErr != nil {
		return
	}

	ars := arecords{Records: make([]arecord, len(aliases))}
	for i, a := range aliases {
		raw, err := s.db.Fetch(ctx, a)
		if err != nil {
			retErr = err
			return
		}
		err = json.Unmarshal(StringToBytes(raw), &ars.Records[i])
		if err != nil {
			retErr = err
			return
		}
	}

	var buf bytes.Buffer
	retErr = statsTmpl.Execute(&buf, ars)
	if retErr != nil {
		return
	}
	w.Write(buf.Bytes())
	return
}

func (s *server) counting(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case visit := <-s.visitCh:
			_, err := s.db.FetchIP(context.Background(), visit.ip)
			if errors.Is(err, redis.Nil) {
				s.db.StoreIP(ctx, visit.ip, visit.alias) // new ip
			} else if err != nil {
				log.Printf("cannot fetch data store for ip processing, err: %v\n", err)
				continue
			} else {
				s.db.UpdateIP(ctx, visit.ip, visit.alias) // old ip
			}
		}
	}
}
