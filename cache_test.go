// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Originally written by Changkun Ou <changkun.de> at
// changkun.de/s/redir, adopted by Mai Yang <maiyang.me>.

package main

import (
	"math/rand"
	"testing"
)

func TestLRU(t *testing.T) {
	l := newLRU(false)
	l.cap = 2 // for testing

	if _, ok := l.Get("a"); ok {
		t.Fatalf("Get value from empty LRU")
	}
	if l.Len() != 0 {
		t.Fatalf("wrong size, want 0, got %v", l.Len())
	}

	l.Put("a", "1") // a
	v, ok := l.Get("a")
	if !ok { // a -> b
		t.Fatalf("Get value from LRU found nothing")
	}
	if l.Len() != 1 {
		t.Fatalf("wrong size, want 1, got %v", l.Len())
	}
	t.Log(v)

	l.Put("b", "2") // b -> a
	v, ok = l.Get("a")
	if !ok { // a -> b
		t.Fatalf("Get value after Put from LRU found nothing")
	}
	if v != "1" {
		t.Fatalf("Get value from LRU want 1 got %v", v)
	}
	if l.Len() != 2 {
		t.Fatalf("wrong size, want 2, got %v", l.Len())
	}
	t.Log(v)

	l.Put("c", "3") // c -> a
	_, ok = l.Get("b")
	if ok {
		t.Fatalf("Get value success meaning LRU incorrect")
	}
	v, ok = l.Get("c")
	if !ok {
		t.Fatalf("Get value fail meaning LRU incorrect")
	}
	if v != "3" {
		t.Fatalf("Get value from LRU want 3 got %v", v)
	}
	if l.Len() != 2 {
		t.Fatalf("wrong size, want 2, got %v", l.Len())
	}
	t.Log(v)

	l.flush()
	if l.Len() != 0 {
		t.Fatalf("wrong size, want 0, got %v", l.Len())
	}
	l.Put("a", "1")
	l.Put("b", "1")
	l.Put("c", "1")
	l.Put("a", "1")
	v, ok = l.Get("a")
	if !ok { // a
		t.Fatalf("Get value from LRU found nothing")
	}
	if v != "1" {
		t.Fatalf("Get value from LRU want 1 got %v", v)
	}
	if l.Len() != 2 {
		t.Fatalf("wrong size, want 2, got %v", l.Len())
	}
	t.Log(v)
}

func rands() string {
	var alphabet = "qazwsxedcrfvtgbyhnujmikolpQAZWSXEDCRFVTGBYHNUJMIKOLP"
	ret := make([]byte, 5)
	for i := 0; i < 5; i++ {
		ret[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(ret)
}

func BenchmarkLRU(b *testing.B) {
	l := newLRU(false)
	l.Put("a", "1")
	b.Run("Get", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Get("a")
			}
		})
	})
	b.Run("Put-Same", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			// each goroutine put its own k/v
			k, v := rands(), rands()
			for pb.Next() {
				l.Put(k, v)
			}
		})
	})

	// This is a very naive bench test, especially it
	// mostly measures the rands().
	b.Run("Put-Different", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// each put has a different k/v
				l.Put(rands(), rands())
			}
		})
	})
}
