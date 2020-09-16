// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"context"
	"testing"
)

func TestOpValid(t *testing.T) {
	tests := []struct {
		o     string
		want  op
		valid bool
	}{
		{o: "", want: op(""), valid: false},
		{o: "create", want: opCreate, valid: true},
		{o: "delete", want: opDelete, valid: true},
		{o: "update", want: opUpdate, valid: true},
		{o: "fetch", want: opFetch, valid: true},
	}
	for _, tt := range tests {
		o := op(tt.o)
		if o != tt.want || o.valid() != tt.valid {
			t.Fatalf("want %v got %v, valid %v but %v", tt.want, o, tt.valid, o.valid())
		}
	}
}

func TestShortCmd(t *testing.T) {
	k, v := "alias", "link"
	ctx := context.Background()

	tests := []struct {
		o       op
		k, v    string
		wantNil bool
	}{
		{o: opCreate, k: "alias", v: "link", wantNil: true},
		{o: opUpdate, k: "alias", v: "link", wantNil: true},
		{o: opFetch, k: "alias", v: "link", wantNil: true},
		{o: opDelete, k: "alias", v: "link", wantNil: true},
		{o: opFetch, k: "alias2", v: "link", wantNil: false},
		{o: opFetch, k: "alias2", v: "link", wantNil: false},
		{o: opDelete, k: "alias2", v: "link", wantNil: true},
	}

	for _, tt := range tests {
		err := shortCmd(ctx, tt.o, tt.k, tt.v)
		if tt.wantNil {
			if err != nil {
				t.Fatalf("shortCmd with err: %v", err)
			}
		} else {
			if err == nil {
				t.Fatalf("shortCmd without error: %v, %v, %v", tt.o, k, v)
			}
		}
	}
}
