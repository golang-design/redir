// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Originally written by Changkun Ou <changkun.de> at
// changkun.de/s/redir, adopted by Mai Yang <maiyang.me>.

package main

import (
	"strings"
	"testing"
)

func TestRandomString(t *testing.T) {
	s1 := randstr(12)
	s2 := randstr(12)
	if len(s1) != 12 || len(s2) != 12 {
		t.Fatalf("want 12 chars, got: %v, %v", len(s1), len(s2))
	}
	if strings.Compare(s1, s2) == 0 {
		t.Fatalf("want two different string, got: %v, %v", s1, s2)
	}
	t.Log(s1, s2)
}
