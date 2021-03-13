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

// Store is persistent storage that provides a group of operations
// to interact with the underlying database.
type Store struct {
	sqlxDB *sqlx.DB
}

// NewDB parses a given DSN and returns a DB instance for
// further operations. It returns an error if the database
// instance is not able to connect.
func NewDB(dsn string) (*Store, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return nil, err
	}
	db.SetMaxOpenConns(120)
	db.SetMaxIdleConns(50)
	return &Store{sqlxDB: db}, nil
}

func (db Store) Close() (err error) {
	err = db.sqlxDB.Close()
	if err != nil {
		err = fmt.Errorf("failed to close database: %w", err)
	}
	return
}

// StoreAlias stores a given short alias with the given link if not exists
func (db Store) StoreAlias(ctx context.Context, r *Redirect) error {
	now := time.Now().UTC()
	query, args, err := sqlx.In(`
INSERT INTO collink (alias, kind, url, private, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?)
`, r.Alias, r.Kind, r.URL, r.Private, now, now)
	if err != nil {
		return err
	}
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAlias updates the link of a given alias
func (db Store) UpdateAlias(ctx context.Context, red *Redirect) error {
	query, args, err := sqlx.In(`
UPDATE collink
SET url=?
WHERE alias=?
`, red.URL, red.Alias)
	if err != nil {
		return err
	}
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAlias deletes a given short alias if exists
func (db Store) DeleteAlias(ctx context.Context, a string) error {
	tx, err := db.sqlxDB.Begin()
	if err != nil {
		return err
	}
	query, args, err := sqlx.In(`
DELETE FROM collink WHERE alias=?
`, a)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(`
DELETE FROM visit WHERE alias=?
`, a)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// FetchAlias reads a given alias and returns the associated link
func (db Store) FetchAlias(ctx context.Context, a string) (*Redirect, error) {
	query, args, err := sqlx.In(`
SELECT alias,
       kind,
       url,
       private
FROM collink
WHERE alias=?
`, a)
	if err != nil {
		return nil, err
	}
	red := []*Redirect{}
	err = db.sqlxDB.SelectContext(ctx, &red, query, args...)
	if err != nil {
		return nil, err
	}
	if len(red) == 0 {
		return nil, sql.ErrNoRows
	}
	return red[0], nil
}

// CountReferer fetches and counts all referers of a given alias
func (db Store) CountReferer(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Refstat, error) {
	query, args, err := sqlx.In(`
SELECT referer,
       COUNT(*) AS count
FROM visit
WHERE alias=?
  AND kind = ?
  AND created_at BETWEEN ? AND ?
GROUP BY referer
`, a, k, start, end)
	if err != nil {
		return nil, err
	}
	ref := []Refstat{}
	err = db.sqlxDB.SelectContext(ctx, &ref, query, args...)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// CountUA fetches and counts all uas of a given alias
func (db Store) CountUA(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]UAstat, error) {
	query, args, err := sqlx.In(`
SELECT ua,
       COUNT(*) AS count
FROM visit
WHERE alias=?
  AND kind = ?
  AND created_at BETWEEN ? AND ?
GROUP BY ua
`, a, k, start, end)
	if err != nil {
		return nil, err
	}
	ref := []UAstat{}
	err = db.sqlxDB.SelectContext(ctx, &ref, query, args...)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// CountVisitHist counts the recorded history every an hour from Visit history.
func (db Store) CountVisitHist(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Timehist, error) {
	query, args, err := sqlx.In(`
SELECT CONVERT(DATE_FORMAT(created_at,'%Y-%m-%d-00:00:00'),DATETIME) AS time,
       COUNT(1) AS count
FROM visit
WHERE alias=?
  AND kind = ?
  AND created_at BETWEEN ? AND ?
GROUP BY time
ORDER BY time
`, a, k, start, end)
	if err != nil {
		return nil, err
	}
	timehists := []Timehist{}
	err = db.sqlxDB.SelectContext(ctx, &timehists, query, args...)
	if err != nil {
		return nil, err
	}
	return timehists, nil
}

// RecordVisit record a given visit data
func (db Store) RecordVisit(ctx context.Context, v *Visit) error {
	query, args, err := sqlx.In(`
INSERT INTO visit (alias, kind, ip, ua, referer, created_at)
VALUES(?, ?, ?, ?, ?, ?)
`, v.Alias, v.Kind, v.IP, v.UA, v.Referer, v.Time)
	if err != nil {
		return err
	}
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// CountVisit count visit by AliasKind
func (db Store) CountVisit(ctx context.Context, kind AliasKind) ([]Record, error) {
	query, args, err := sqlx.In(`
SELECT alias,
       COUNT(*) pv,
       COUNT(DISTINCT ip) uv
FROM visit
WHERE kind = ?
GROUP BY alias
`, kind)
	if err != nil {
		return nil, err
	}
	rs := []Record{}
	err = db.sqlxDB.SelectContext(ctx, &rs, query, args...)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
