package model

import (
	"context"
	"errors"
	"time"
)

var (
	ErrExistedAlias = errors.New("alias is existed")
)

type AliasKind int

const (
	KindShort AliasKind = iota
	KindRandom
)

const (
	dbname   = "redir"
	collink  = "links"
	colvisit = "Visit"
)

// Redirect records a kind of alias and its correlated link.
type Redirect struct {
	Alias   string    `json:"alias"   bson:"alias"`
	Kind    AliasKind `json:"kind"    bson:"kind"`
	URL     string    `json:"url"     bson:"url"`
	Private bool      `json:"private" bson:"private"`
}

// Visit indicates an Record of Visit pattern.
type Visit struct {
	Alias   string    `json:"alias"   bson:"alias"`
	Kind    AliasKind `json:"kind"    bson:"kind"`
	IP      string    `json:"ip"      bson:"ip"`
	UA      string    `json:"ua"      bson:"ua"`
	Referer string    `json:"referer" bson:"referer"`
	Time    time.Time `json:"time"    bson:"time"`
}

type Refstat struct {
	Referer string `json:"referer" bson:"referer"`
	Count   int64  `json:"count"   bson:"count"`
}

type Uastat struct {
	UA    string `json:"ua"    bson:"ua"`
	Count int64  `json:"count" bson:"count"`
}

type Locstat struct {
	Locations []string `bson:"locs" json:"locs"`
}

type Timehist struct {
	Time  time.Time `bson:"time"  json:"time"`
	Count int       `bson:"count" json:"count"`
}

type Record struct {
	Alias string `bson:"alias"`
	UV    int64  `bson:"uv"`
	PV    int64  `bson:"pv"`
}

type RedirModel interface {
	NewDB(url string) (RedirModel, error)
	Close() (err error)
	StoreAlias(ctx context.Context, r *Redirect) (err error)
	UpdateAlias(ctx context.Context, a, l string) (*Redirect, error)
	DeleteAlias(ctx context.Context, a string) (err error)
	FetchAlias(ctx context.Context, a string) (*Redirect, error)
	CountReferer(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Refstat, error)
	CountUA(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Uastat, error)
	CountLocation(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]string, error)
	CountVisitHist(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Timehist, error)
	RecordVisit(ctx context.Context, v *Visit) (err error)
	CountVisit(ctx context.Context, kind AliasKind) (rs []Record, err error)
}
