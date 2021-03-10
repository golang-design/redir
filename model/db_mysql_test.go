// Copyright 2021 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package model

import (
	"context"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

var dsn = "root:dc@tcp(127.0.0.1:3306)/" + dbname + "?charset=utf8mb4&parseTime=true"

func TestStoreAlias(t *testing.T) {
	sqlxDB, err := initMySQL(dsn)
	require.NoError(t, err)
	ctx := context.Background()
	_, err = sqlxDB.ExecContext(ctx, `truncate table collink`)
	require.NoError(t, err)
	alias := "T1"
	redirURL := "https://golang.design/s/talkgo"
	red := &Redirect{
		Alias:   alias,
		Kind:    0,
		URL:     redirURL,
		Private: false,
	}

	mySQLDB := MySQLDB{sqlxDB: sqlxDB}
	err = mySQLDB.StoreAlias(ctx, red)
	require.NoError(t, err)
	ret, err := mySQLDB.FetchAlias(ctx, alias)
	require.NoError(t, err)
	require.EqualValues(t, redirURL, ret.URL)
	ret.URL = "https://golang.design/s/xlff"
	ret.Kind = 1
	err = mySQLDB.UpdateAlias(ctx, ret)
	require.NoError(t, err)
	ret, err = mySQLDB.FetchAlias(ctx, alias)
	require.NoError(t, err)
	require.NotEqualValues(t, redirURL, ret.URL)
	err = mySQLDB.DeleteAlias(ctx, alias)
	require.NoError(t, err)
	ret, err = mySQLDB.FetchAlias(ctx, alias)
	require.NoError(t, err)
	require.Nil(t, ret)
}

func TestVisit(t *testing.T) {
	sqlxDB, err := initMySQL(dsn)
	require.NoError(t, err)
	ctx := context.Background()
	_, err = sqlxDB.ExecContext(ctx, `truncate table visit`)
	require.NoError(t, err)

	now := time.Now()
	visits := []*Visit{
		{"t1", 1, "192.168.0.1", "ua1", "chrome", now.Add(-12 * time.Second)},
		{"t1", 1, "192.168.0.2", "ua2", "chrome", now.Add(-10 * time.Second)},
		{"t2", 1, "192.168.0.3", "ua2", "safari", now},
		{"t3", 0, "192.168.0.2", "ua3", "firefox", now},
		{"t3", 1, "192.168.0.3", "ua4", "firefox", now},
		{"t3", 1, "192.168.0.4", "ua5", "chrome", now},
	}

	mySQLDB := MySQLDB{sqlxDB: sqlxDB}
	for _, v := range visits {
		err = mySQLDB.RecordVisit(ctx, v)
		require.NoError(t, err)
	}
	ret, err := mySQLDB.CountReferer(ctx, "t2", 1, now.Add(-1*time.Second), now.Add(1*time.Second))
	require.NoError(t, err)
	require.Len(t, ret, 1)

	uast, err := mySQLDB.CountUA(ctx, "t2", 1, now.Add(-1*time.Second), now.Add(1*time.Second))
	require.NoError(t, err)
	require.Len(t, uast, 1)

	locst, err := mySQLDB.CountLocation(ctx, "t2", 1, now.Add(-1*time.Second), now.Add(1*time.Second))
	require.NoError(t, err)
	require.Len(t, locst, 1)

	cvhst, err := mySQLDB.CountVisitHist(ctx, "t2", 1, now.Add(-10*time.Second), now.Add(1*time.Second))
	require.NoError(t, err)
	require.Len(t, cvhst, 1)
	require.EqualValues(t, 1, cvhst[0].Count)

	cvhst, err = mySQLDB.CountVisitHist(ctx, "t3", 1, now.Add(-10*time.Second), now.Add(1*time.Second))
	require.NoError(t, err)
	require.Len(t, cvhst, 1)
	require.EqualValues(t, 2, cvhst[0].Count)

	rs, err := mySQLDB.CountVisit(ctx, 1)
	require.NoError(t, err)
	require.Len(t, rs, 3)
}
