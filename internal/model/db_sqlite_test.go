// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Originally written by Mai Yang <maiyang.me>.

package model

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dsn = "../../data/redir.db"

func TestStoreAlias(t *testing.T) {
	db, err := NewDB(dsn)
	if err != nil {
		t.Fatalf("NewDB with err: %v", err)
	}
	ctx := context.Background()
	alias := "T1"
	redirURL := "https://golang.design/s/talkgo"
	red := &Redirect{
		Alias:   alias,
		Kind:    0,
		URL:     redirURL,
		Private: false,
	}

	err = db.StoreAlias(ctx, red)
	if err != nil {
		t.Fatalf("StoreAlias with err: %v", err)
	}
	ret, err := db.FetchAlias(ctx, alias)
	if err != nil {
		t.Fatalf("FetchAlias with err: %v", err)
	}
	if ret.URL != redirURL {
		t.Fatalf("wrong alias URL, want %s, got %v", redirURL, ret.URL)
	}
	ret.URL = "https://golang.design/s/go-questions"
	ret.Kind = 1
	err = db.UpdateAlias(ctx, ret)
	if err != nil {
		t.Fatalf("UpdateAlias with err: %v", err)
	}
	ret, err = db.FetchAlias(ctx, alias)
	if err != nil {
		t.Fatalf("FetchAlias with err: %v", err)
	}
	if ret.URL == redirURL {
		t.Fatalf("FetchAlias URL must be not equal, but got %v", ret.URL)
	}
	err = db.DeleteAlias(ctx, alias)
	if err != nil {
		t.Fatalf("DeleteAlias with err: %v", err)
	}
	ret, err = db.FetchAlias(ctx, alias)
	if err == nil {
		t.Fatalf("FetchAlias with err: %v", err)
	}
	if ret != nil {
		t.Fatalf("FetchAlias return alias must be nil, got %+v", ret)
	}
}

func TestVisit(t *testing.T) {
	db, err := NewDB(dsn)
	if err != nil {
		t.Fatalf("NewDB with err: %v", err)
	}
	ctx := context.Background()

	now := time.Now()
	visits := []*Visit{
		{"t1", 1, "192.168.0.1", "ua1", "https://example1.com", now.Add(-12 * time.Second)},
		{"t1", 1, "192.168.0.2", "ua2", "https://example2.com", now.Add(-10 * time.Second)},
		{"t2", 1, "192.168.0.3", "ua2", "https://example3.com", now},
		{"t3", 0, "192.168.0.2", "ua3", "https://example4.com", now},
		{"t3", 1, "192.168.0.3", "ua4", "https://example5.com", now},
		{"t3", 1, "192.168.0.4", "ua5", "https://example6.com", now},
	}

	for _, v := range visits {
		err = db.RecordVisit(ctx, v)
		if err != nil {
			t.Fatalf("RecordVisit with err: %v", err)
		}
	}
	ret, err := db.CountReferer(ctx, "t2", 1, now.Add(-1*time.Second), now.Add(1*time.Second))
	if err != nil {
		t.Fatalf("CountReferer with err: %v", err)
	}
	if len(ret) != 1 {
		t.Fatalf("CountReferer len is 1, but got: %d", len(ret))
	}

	uast, err := db.CountUA(ctx, "t2", 1, now.Add(-1*time.Second), now.Add(1*time.Second))
	if err != nil {
		t.Fatalf("CountUA with err: %v", err)
	}
	if len(uast) != 1 {
		t.Fatalf("CountUA result is not equal 1, got: %v", len(uast))
	}

	cvhst, err := db.CountVisitHist(ctx, "t2", 1, now.Add(-10*time.Second), now.Add(1*time.Second))
	if err != nil {
		t.Fatalf("CountVisitHist with err: %v", err)
	}
	if len(cvhst) != 1 {
		t.Fatalf("CountUA result is not equal 1, got: %v", len(cvhst))
	}
	if cvhst[0].Count != 1 {
		t.Fatalf("CountUA result want 1, got: %v", cvhst[0].Count)
	}

	cvhst, err = db.CountVisitHist(ctx, "t3", 1, now.Add(-10*time.Second), now.Add(1*time.Second))
	if err != nil {
		t.Fatalf("CountVisitHist with err: %v", err)
	}
	if len(cvhst) != 1 {
		t.Fatalf("CountVisitHist result is not equal 1, got: %v", len(cvhst))
	}
	if cvhst[0].Count != 2 {
		t.Fatalf("CountVisitHist result want 1, got: %v", cvhst[0].Count)
	}

	rs, err := db.CountVisit(ctx, 1)
	if err != nil {
		t.Fatalf("CountVisit with err: %v", err)
	}
	if len(rs) != 3 {
		t.Fatalf("CountVisit result is not equal 3, got: %v", len(rs))
	}
}
