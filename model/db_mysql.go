// Copyright 2021 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package model

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLDB struct {
	RedirModel
	sqlxDB *sqlx.DB
}

func NewDB(uri string) (*MySQLDB, error) {
	dsn := "root:dc@tcp(127.0.0.1:3306)/" + dbname + "?charset=utf8mb4&parseTime=true"
	fmt.Println(dsn)
	sqlxDB, err := initMySQL(uri)
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return nil, err
	}
	return &MySQLDB{sqlxDB: sqlxDB}, nil
}

// 初始化数据库
func initMySQL(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return nil, err
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	return db, nil
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
	INSERT INTO collink (alias, kind, url, private, created_at, updated_at)
	VALUES(?, ?, ?, ?, ?, ?)
`, r.Alias, r.Kind, r.URL, r.Private, now, now)
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (db MySQLDB) UpdateAlias(ctx context.Context, red *Redirect) error {
	query, args, err := sqlx.In(`UPDATE collink SET url=? WHERE alias=?`, red.URL, red.Alias)
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (db MySQLDB) DeleteAlias(ctx context.Context, a string) (err error) {
	query, args, err := sqlx.In(`UPDATE collink set is_deleted=true WHERE alias=?`, a)
	_, err = db.sqlxDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (db MySQLDB) FetchAlias(ctx context.Context, a string) (*Redirect, error) {
	query, args, err := sqlx.In(`
SELECT alias, kind, url, private FROM collink WHERE alias=? AND is_deleted=false
`, a)
	red := []*Redirect{}
	err = db.sqlxDB.SelectContext(ctx, &red, query, args...)
	if err != nil {
		return nil, err
	}
	if len(red) > 0 {
		return red[0], nil
	}
	return nil, nil
}

func (db MySQLDB) CountReferer(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Refstat, error) {
	query, args, err := sqlx.In(`
SELECT referer, count(*) as count FROM visit
WHERE alias=? AND kind = ? AND is_deleted=false
AND created_at between ? and ?
GROUP BY referer
`, a, k, start, end)
	ref := []Refstat{}
	err = db.sqlxDB.SelectContext(ctx, &ref, query, args...)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

func (db MySQLDB) CountUA(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Uastat, error) {
	query, args, err := sqlx.In(`
SELECT ua, count(*) as count FROM visit
WHERE alias=? AND kind = ? AND is_deleted=false
AND created_at between ? and ?
GROUP BY ua
`, a, k, start, end)
	ref := []Uastat{}
	err = db.sqlxDB.SelectContext(ctx, &ref, query, args...)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

func (db MySQLDB) CountLocation(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]string, error) {
	query, args, err := sqlx.In(`
SELECT ip FROM visit
WHERE alias=? AND kind = ? AND is_deleted=false
AND created_at between ? and ?
`, a, k, start, end)
	locs := []string{}
	err = db.sqlxDB.SelectContext(ctx, &locs, query, args...)
	if err != nil {
		return nil, err
	}
	return locs, nil
}

func (db MySQLDB) CountVisitHist(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Timehist, error) {
	query, args, err := sqlx.In(`
SELECT CONVERT(DATE_FORMAT(created_at,'%Y-%m-%d-%H:30:00'),DATETIME) as time, count(1) as count
FROM visit
WHERE alias=? AND kind = ? AND is_deleted=false
AND created_at between ? and ?
GROUP BY time order by time
`, a, k, start, end)
	timehists := []Timehist{}
	err = db.sqlxDB.SelectContext(ctx, &timehists, query, args...)
	if err != nil {
		return nil, err
	}
	return timehists, nil
}

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

func (db MySQLDB) CountVisit(ctx context.Context, kind AliasKind) (rs []Record, err error) {
	query, args, err := sqlx.In(`
SELECT alias, count(*) pv, count(distinct ip) uv
FROM visit
WHERE kind = ? AND is_deleted=false
GROUP BY alias
`, kind)
	err = db.sqlxDB.SelectContext(ctx, &rs, query, args...)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
