// Copyright 2020 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package model

import (
	"context"
	"testing"
)

const kalias = "alias"

func prepare(ctx context.Context, t *testing.T) *MongoDB {
	s, err := newDB("mongodb://0.0.0.0:27017")
	if err != nil {
		t.Fatalf("cannot connect to data store")
	}

	err = s.StoreAlias(ctx, &Redirect{
		Alias:   kalias,
		Kind:    KindShort,
		URL:     "link",
		Private: false,
	})
	if err != nil {
		t.Fatalf("cannot store alias to data store: %v\n", err)
	}
	return s
}

func TestUpdateAlias(t *testing.T) {
	want := "link2"

	ctx := context.Background()
	s := prepare(ctx, t)

	r, err := s.UpdateAlias(ctx, kalias, want)
	if err != nil {
		t.Fatalf("UpdateAlias failed with err: %v", err)
	}

	if r.URL != want {
		t.Fatalf("Incorrect UpdateAlias implementaiton, want %v, got %v", want, r.URL)
	}
}
