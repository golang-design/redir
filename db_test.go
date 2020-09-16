// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
)

func TestReadAndUpdate(t *testing.T) {
	s, err := newStore("redis://localhost:6379/8")
	if err != nil {
		t.Fatalf("cannot connect to data store")
	}

	ctx := context.Background()

	err = s.StoreAlias(ctx, "alias", "link")
	if err != nil {
		t.Fatalf("cannot store alias to data store, err: %v\n", err)
	}
	t.Cleanup(func() {
		err := s.DeleteAlias(ctx, "alias")
		if err != nil {
			t.Fatalf("DeleteAlias failure, err: %v", err)
		}
		s.Close()
	})

	// create a number of concurrent updater
	// check data is still consistent
	concurrent := 1000
	wg := sync.WaitGroup{}
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()
			err := s.countVisit(ctx, "alias", 1, 1)
			if err != nil {
				t.Fatalf("countVisit failure, err: %v\n", err)
			}
		}()
	}
	wg.Wait()

	r, err := s.FetchAlias(ctx, "alias")
	if err != nil {
		t.Fatalf("FetchAlias failure, err: %v\n", err)
	}

	rr := arecord{}
	err = json.Unmarshal([]byte(r), &rr)
	if err != nil {
		t.Fatalf("Unmarshal failure, err: %v\n", err)
	}

	if rr.PV != uint64(concurrent) || rr.UV != uint64(concurrent) {
		t.Fatalf("Incorrect atomic readAndUpdate implementaiton: pv:%v, uv:%v\n", rr.PV, rr.UV)
	}
}
