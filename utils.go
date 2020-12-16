// Copyright 2020 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"time"
	"unsafe"
)

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
}

const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomString generates a random string
func RandomString(n int) string {
	var str string
	length := len(alphanum)
	for i := 0; i < n; i++ {
		a := alphanum[r.Intn(len(alphanum))%length]
		str += string(a)
	}
	return str
}
