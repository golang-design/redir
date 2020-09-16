// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	conf.parse()

	// Test if all fields are filled.
	v := reflect.ValueOf(conf)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			continue
		}
		vv, ok := v.Field(i).Interface().(string)
		if ok && len(vv) > 0 {
			continue
		}
		t.Fatalf("read empty from config, field: %v", v.Type().Field(i).Name)
	}
}
