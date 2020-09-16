// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import "testing"

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

// TODO: integration test
