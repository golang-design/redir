// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Originally written by Mai Yang <maiyang.me>.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQLDB MySQL DB struct
type MySQLDB struct {
	sqlxDB *sqlx.DB
}

// NewDB new MySQLDB from uri
func NewDB(dsn string) (*MySQLDB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return nil, err
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	return &MySQLDB{sqlxDB: db}, nil
}

func (db MySQLDB) Close() (err error) {
	err = db.sqlxDB.Close()
	if err != nil {
		err = fmt.Errorf("failed to close database: %w", err)
	}
	return
}

// StoreAlias stores a given short alias with the given link if not exists
func (db MySQLDB) StoreAlias(ctx context.Context, r *Redirect) (err error) {
	now := time.Now()
	query, args, err := sqlx.In(`
INSERT INTO collink (ALIAS, kind, url, private, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?)
`, r.Alias, r.Kind, r.URL, r.Private, now, now)
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAlias updates the link of a given alias
func (db MySQLDB) UpdateAlias(ctx context.Context, red *Redirect) error {
	query, args, err := sqlx.In(`
UPDATE collink
SET url=?
WHERE ALIAS=?
`, red.URL, red.Alias)
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAlias deletes a given short alias if exists
func (db MySQLDB) DeleteAlias(ctx context.Context, a string) (err error) {
	query, args, err := sqlx.In(`
UPDATE collink
SET is_deleted=TRUE
WHERE ALIAS=?
`, a)
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// FetchAlias reads a given alias and returns the associated link
func (db MySQLDB) FetchAlias(ctx context.Context, a string) (*Redirect, error) {
	query, args, err := sqlx.In(`
SELECT alias,
       kind,
       url,
       private
FROM collink
WHERE alias=?
  AND is_deleted=FALSE
`, a)
	red := []*Redirect{}
	err = db.sqlxDB.SelectContext(ctx, &red, query, args...)
	if err != nil {
		return nil, err
	}
	if len(red) > 0 {
		return red[0], nil
	}
	return nil, sql.ErrNoRows
}

// CountReferer fetches and counts all referers of a given alias
func (db MySQLDB) CountReferer(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Refstat, error) {
	query, args, err := sqlx.In(`
SELECT referer,
       COUNT(*) AS count
FROM visit
WHERE alias=?
  AND kind = ?
  AND is_deleted=FALSE
  AND created_at BETWEEN ? AND ?
GROUP BY referer
`, a, k, start, end)
	ref := []Refstat{}
	err = db.sqlxDB.SelectContext(ctx, &ref, query, args...)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// CountUA fetches and counts all uas of a given alias
func (db MySQLDB) CountUA(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Uastat, error) {
	query, args, err := sqlx.In(`
SELECT ua,
       COUNT(*) AS count
FROM visit
WHERE alias=?
  AND kind = ?
  AND is_deleted=FALSE
  AND created_at BETWEEN ? AND ?
GROUP BY ua
`, a, k, start, end)
	ref := []Uastat{}
	err = db.sqlxDB.SelectContext(ctx, &ref, query, args...)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// CountLocation counts the recorded IPs from Visit history.
// FIXME: IP can be changed overtime, it might be a good idea to just store
// the parse geo location (latitude, and longitude, and accuracy).
// Q: Any APIs can convert IP to geo location?
func (db MySQLDB) CountLocation(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]string, error) {
	query, args, err := sqlx.In(`
SELECT ip
FROM visit
WHERE alias=?
  AND kind = ?
  AND is_deleted=FALSE
  AND created_at BETWEEN ? AND ?
`, a, k, start, end)
	locs := []string{}
	err = db.sqlxDB.SelectContext(ctx, &locs, query, args...)
	if err != nil {
		return nil, err
	}
	return locs, nil
}

// CountVisitHist counts the recorded history every 30 minutes from Visit history.
func (db MySQLDB) CountVisitHist(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Timehist, error) {
	query, args, err := sqlx.In(`
SELECT CONVERT(DATE_FORMAT(created_at,'%Y-%m-%d-%H:30:00'),DATETIME) AS time,
       COUNT(1) AS count
FROM visit
WHERE alias=?
  AND kind = ?
  AND is_deleted=FALSE
  AND created_at BETWEEN ? AND ?
GROUP BY time
ORDER BY time
`, a, k, start, end)
	timehists := []Timehist{}
	err = db.sqlxDB.SelectContext(ctx, &timehists, query, args...)
	if err != nil {
		return nil, err
	}
	return timehists, nil
}

// RecordVisit record a given visit data
func (db MySQLDB) RecordVisit(ctx context.Context, v *Visit) (err error) {
	now := time.Now()
	query, args, err := sqlx.In(`
INSERT INTO visit (alias, kind, ip, ua, referer, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?)
`, v.Alias, v.Kind, v.IP, v.UA, v.Referer, now, now)
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// CountVisit count visit by AliasKind
func (db MySQLDB) CountVisit(ctx context.Context, kind AliasKind) (rs []Record, err error) {
	query, args, err := sqlx.In(`
SELECT alias,
       COUNT(*) pv,
       COUNT(DISTINCT ip) uv
FROM visit
WHERE kind = ?
  AND is_deleted=FALSE
GROUP BY alias
`, kind)
	err = db.sqlxDB.SelectContext(ctx, &rs, query, args...)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
