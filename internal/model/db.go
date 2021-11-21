// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Originally written by Changkun Ou <changkun.de> at
// changkun.de/s/redir, adopted by Mai Yang <maiyang.me>.

package model

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrExistedAlias indicates an error where an alias is already existed in the data store.
	ErrExistedAlias = errors.New("alias is existed")
)

// Redirect records alias and its correlated link.
type Redirect struct {
	Alias string `json:"alias"   db:"alias"`
	URL   string `json:"url"     db:"url"`
}

// Visit indicates an Record of Visit pattern.
type Visit struct {
	Alias   string    `json:"alias"   db:"alias"`
	IP      string    `json:"ip"      db:"ip"`
	UA      string    `json:"ua"      db:"ua"`
	Referer string    `json:"referer" db:"referer"`
	Time    time.Time `json:"time"    db:"time"`
}

// Refstat counts the occurrence of a referer.
type Refstat struct {
	Referer string `json:"referer"`
	Count   int64  `json:"count"`
}

// UAstat counts the occurrence of a user agent.
type UAstat struct {
	UA    string `json:"ua"`
	Count int64  `json:"count"`
}

// Timehist counts the occurrence of visit events during a time period (an hour).
type Timehist struct {
	Time  time.Time `json:"time"`
	Count int       `json:"count"`
}

// Record contains a record of alias's UV/PV
type Record struct {
	Alias string `json:"alias"`
	UV    int64  `json:"uv"`
	PV    int64  `json:"pv"`
}

type RedirAliasDataModel interface {
	StoreAlias(context.Context, *Redirect) error
	UpdateAlias(ctx context.Context, red *Redirect) error
	DeleteAlias(ctx context.Context, alias string) error
	FetchAlias(ctx context.Context, alias string) (*Redirect, error)
}

type RedirVisitDataModel interface {
	RecordVisit(context.Context, *Visit) error
}

type RedirStatModel interface {
	CountReferer(ctx context.Context, alias string, start, end time.Time) ([]Refstat, error)
	CountUA(ctx context.Context, alias string, start, end time.Time) ([]UAstat, error)
	CountVisitHist(ctx context.Context, alias string, start, end time.Time) ([]Timehist, error)
	CountVisit(context.Context) (rs []Record, err error)
}
