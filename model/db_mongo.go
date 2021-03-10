// Copyright 2021 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package model

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	RedirModel
	cli *mongo.Client
}

func newDB(uri string) (*MongoDB, error) {
	// initialize database connection
	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	return &MongoDB{nil, db}, nil
}

func (db *MongoDB) Close() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = db.cli.Disconnect(ctx)
	if err != nil {
		err = fmt.Errorf("failed to close database: %w", err)
	}
	return
}

// StoreAlias stores a given short alias with the given link if not exists
func (db *MongoDB) StoreAlias(ctx context.Context, r *Redirect) (err error) {
	col := db.cli.Database(dbname).Collection(collink)

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"alias": r.Alias, "kind": r.Kind}

	_, err = col.UpdateOne(ctx, filter, bson.M{"$set": r}, opts)
	if err != nil {
		err = fmt.Errorf("failed to insert given Redirect: %w", err)
		return
	}
	return
}

// UpdateAlias updates the link of a given alias
func (db *MongoDB) UpdateAlias(ctx context.Context, a, l string) (*Redirect, error) {
	col := db.cli.Database(dbname).Collection(collink)

	var r Redirect
	err := col.FindOneAndUpdate(ctx,
		bson.M{"alias": a},
		bson.M{"$set": bson.M{"url": l}},
	).Decode(&r)
	if err != nil {
		err = fmt.Errorf("failed to update alias %s: %v", a, err)
		return nil, err
	}
	r.URL = l
	return &r, nil
}

// Delete deletes a given short alias if exists
func (db *MongoDB) DeleteAlias(ctx context.Context, a string) (err error) {
	col := db.cli.Database(dbname).Collection(collink)

	_, err = col.DeleteMany(ctx, bson.M{"alias": a})
	if err != nil {
		err = fmt.Errorf("delete alias %s failed: %w", a, err)
		return
	}
	return
}

// FetchAlias reads a given alias and returns the associated link
func (db *MongoDB) FetchAlias(ctx context.Context, a string) (*Redirect, error) {
	col := db.cli.Database(dbname).Collection(collink)

	var r Redirect
	err := col.FindOne(ctx, bson.M{"alias": a}).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("cannot find alias %s: %v", a, err)
	}
	return &r, nil
}

// CountReferer fetches and counts all referers of a given alias
func (db *MongoDB) CountReferer(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Refstat, error) {
	col := db.cli.Database(dbname).Collection(collink)
	opts := options.Aggregate().SetMaxTime(10 * time.Second)
	cur, err := col.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			primitive.E{Key: "$match", Value: bson.M{
				"kind": k, "alias": a,
			}},
		},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.M{
				"from": colvisit,
				"as":   "Visit",
				"pipeline": mongo.Pipeline{bson.D{
					primitive.E{Key: "$match", Value: bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$eq": []string{a, "$alias"}},
								{"$gte": []interface{}{"$time", start}},
								{"$lt": []interface{}{"$time", end}},
							},
						},
					}},
				}},
			}},
		},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.M{
				"path":                       "$Visit",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		bson.D{
			primitive.E{Key: "$group", Value: bson.M{
				"_id": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []string{"", "$Visit.referer"},
						},
						"then": "unknown",
						"else": "$Visit.referer",
					},
				},
				"referer": bson.M{"$first": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []string{"", "$Visit.referer"},
						},
						"then": "unknown",
						"else": "$Visit.referer",
					},
				}},
				"count": bson.M{"$sum": 1},
			}},
		},
		bson.D{
			primitive.E{Key: "$sort", Value: bson.M{"count": -1}},
		},
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to count referer: %w", err)
	}
	defer cur.Close(ctx)

	var results []Refstat
	if err := cur.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to fetch referer results: %w", err)
	}

	return results, nil
}

func (db *MongoDB) CountUA(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Uastat, error) {
	col := db.cli.Database(dbname).Collection(collink)
	opts := options.Aggregate().SetMaxTime(10 * time.Second)
	cur, err := col.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			primitive.E{Key: "$match", Value: bson.M{
				"kind": k, "alias": a,
			}},
		},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.M{
				"from": colvisit,
				"as":   "Visit",
				"pipeline": mongo.Pipeline{bson.D{
					primitive.E{Key: "$match", Value: bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$eq": []string{a, "$alias"}},
								{"$gte": []interface{}{"$time", start}},
								{"$lt": []interface{}{"$time", end}},
							},
						},
					}},
				}},
			}},
		},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.M{
				"path":                       "$Visit",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		bson.D{
			primitive.E{Key: "$group", Value: bson.M{
				"_id": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []string{"", "$Visit.ua"},
						},
						"then": "unknown",
						"else": "$Visit.ua",
					},
				},
				"ua": bson.M{"$first": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$eq": []string{"", "$Visit.ua"},
						},
						"then": "unknown",
						"else": "$Visit.ua",
					},
				}},
				"count": bson.M{"$sum": 1},
			}},
		},
		bson.D{
			primitive.E{Key: "$sort", Value: bson.M{"count": -1}},
		},
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to count ua: %w", err)
	}
	defer cur.Close(ctx)

	var results []Uastat
	if err := cur.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to fetch ua results: %w", err)
	}

	return results, nil
}

// CountLocation counts the recorded IPs from Visit history.
// FIXME: IP can be changed overtime, it might be a good idea to just store
// the parse geo location (latitude, and longitude, and accuracy).
// Q: Any APIs can convert IP to geo location?
func (db *MongoDB) CountLocation(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]string, error) {
	col := db.cli.Database(dbname).Collection(collink)
	opts := options.Aggregate().SetMaxTime(10 * time.Second)

	// db.links.aggregate([ {$match: {kind: 0, alias: 'gp-1-intro'}}, {'$lookup': {from: 'Visit', localField: 'alias', foreignField: 'alias', as: 'Visit'}}, {'$group': {_id: '$alias', ip: {'$first': '$Visit.ip'}}}, ])
	cur, err := col.Aggregate(ctx, mongo.Pipeline{
		bson.D{primitive.E{
			Key: "$match", Value: bson.M{
				"kind": k, "alias": a,
			},
		}},
		bson.D{primitive.E{
			Key: "$lookup", Value: bson.M{
				"from": colvisit,
				"as":   "Visit",
				"pipeline": mongo.Pipeline{bson.D{
					primitive.E{Key: "$match", Value: bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$eq": []string{a, "$alias"}},
								{"$gte": []interface{}{"$time", start}},
								{"$lt": []interface{}{"$time", end}},
							},
						},
					}},
				}},
			},
		}},
		bson.D{primitive.E{
			Key: "$group",
			Value: bson.M{
				"_id":  "$alias",
				"locs": bson.M{"$first": "$Visit.ip"},
			},
		}},
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to count ua: %w", err)
	}
	defer cur.Close(ctx)

	var results []Locstat
	if err := cur.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to fetch ua results: %w", err)
	}

	// Is it possible that we don't have any result or have mutiple entries?
	if len(results) != 1 {
		return []string{}, nil
	}

	return results[0].Locations, nil
}

func (db *MongoDB) CountVisitHist(ctx context.Context, a string, k AliasKind, start, end time.Time) ([]Timehist, error) {
	// db.links.aggregate([
	//     {$match: {kind: 0, alias: 'blog'}},
	//     {'$lookup': {from: 'Visit', localField: 'alias', foreignField: 'alias', as: 'Visit'}},
	//     {
	//         $group: {
	//             _id: "$alias", time: {'$first': '$Visit.time'}
	//         },
	//     },
	//     {'$unwind': "$time"},
	//     {
	//         $project: {
	//             "year": {$year: "$time"},
	//             "month": {$month: "$time"},
	//             "day": {$dayOfMonth: "$time"},
	//             "hour": {$hour: "$time"},
	//         },
	//     },
	//     {
	//         "$group":{
	//             "_id": {
	//                 "year":"$year","month":"$month","day":"$day","hour":"$hour",
	//             },
	//             'year': {'$first': '$year'},
	//             'month': {'$first': '$month'},
	//             'day': {'$first': '$day'},
	//             'hour': {'$first': '$hour'},
	//             'count': {$sum: 1},
	//         },
	//     },
	//     {
	//         $sort: {
	//             year: -1,
	//             month: -1,
	//             day: -1,
	//             hour: -1,
	//         },
	//     },
	// ])

	col := db.cli.Database(dbname).Collection(collink)
	opts := options.Aggregate().SetMaxTime(10 * time.Second)
	cur, err := col.Aggregate(ctx, mongo.Pipeline{
		bson.D{primitive.E{
			Key: "$match", Value: bson.M{
				"kind": k, "alias": a,
			},
		}},
		bson.D{primitive.E{
			Key: "$lookup", Value: bson.M{
				"from": colvisit,
				"as":   "Visit",
				"pipeline": mongo.Pipeline{bson.D{
					primitive.E{Key: "$match", Value: bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$eq": []string{a, "$alias"}},
								{"$gte": []interface{}{"$time", start}},
								{"$lt": []interface{}{"$time", end}},
							},
						},
					}},
				}},
			},
		}},
		bson.D{primitive.E{
			Key: "$group", Value: bson.M{
				"_id":  "$alias",
				"time": bson.M{"$first": "$Visit.time"},
			},
		}},
		bson.D{primitive.E{
			Key: "$unwind", Value: bson.M{
				"path": "$time",
			},
		}},
		bson.D{primitive.E{
			Key: "$project", Value: bson.M{
				"year":  bson.M{"$year": "$time"},
				"month": bson.M{"$month": "$time"},
				"day":   bson.M{"$dayOfMonth": "$time"},
				"hour":  bson.M{"$hour": "$time"},
			},
		}},
		bson.D{primitive.E{
			Key: "$group",
			Value: bson.M{
				"_id": bson.M{
					"$dateFromParts": bson.M{
						"year":  "$year",
						"month": "$month",
						"day":   "$day",
						"hour":  "$hour",
					},
				},
				"time": bson.M{"$first": bson.M{
					"$dateFromParts": bson.M{
						"year":  "$year",
						"month": "$month",
						"day":   "$day",
						"hour":  "$hour",
					},
				}},
				"count": bson.M{"$sum": 1},
			},
		}},
		bson.D{primitive.E{
			Key: "$sort", Value: bson.M{
				"_id": -1,
			},
		}},
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to count time hist: %w", err)
	}
	defer cur.Close(ctx)

	var results []Timehist
	if err := cur.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to fetch time hist results: %w", err)
	}

	return results, nil
}

func (db *MongoDB) RecordVisit(ctx context.Context, v *Visit) (err error) {
	col := db.cli.Database(dbname).Collection(colvisit)

	_, err = col.InsertOne(ctx, v)
	if err != nil {
		err = fmt.Errorf("failed to insert Record: %w", err)
		return
	}
	return
}

// CountVisit counts the PV/UV of aliases of a given kind
func (db *MongoDB) CountVisit(ctx context.Context, kind AliasKind) (rs []Record, err error) {
	// uv based on number of ip, this is not accurate since the number will be
	// smaller than the actual.
	// raw query:
	//
	// db.links.aggregate([
	// 	{$match: {kind: 0}},
	// 	{'$lookup': {from: 'Visit', localField: 'alias', foreignField: 'alias', as: 'Visit'}},
	// 	{'$unwind': {path: '$Visit', preserveNullAndEmptyArrays: true}},
	// 	{$group: {_id: {alias: '$alias', ip: '$Visit.ip'}, count: {$sum: 1}}},
	// 	{$group: {_id: '$_id.alias', uv: {$sum: 1}, pv: {$sum: '$count'}}},
	// 	{$sort : {pv: -1}},
	// 	{$sort : {uv: -1}},
	// ])
	col := db.cli.Database(dbname).Collection(collink)
	opts := options.Aggregate().SetMaxTime(10 * time.Second)
	cur, err := col.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			primitive.E{Key: "$match", Value: bson.M{
				"kind": kind,
			}},
		},
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.M{
				"from":         colvisit,
				"localField":   "alias",
				"foreignField": "alias",
				"as":           "Visit",
			}},
		},
		bson.D{
			primitive.E{Key: "$unwind", Value: bson.M{
				"path":                       "$Visit",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		bson.D{
			primitive.E{Key: "$group", Value: bson.M{
				"_id":   bson.M{"alias": "$alias", "ip": "$Visit.ip"},
				"count": bson.M{"$sum": 1},
			}},
		},
		bson.D{
			primitive.E{Key: "$group", Value: bson.M{
				"_id":   "$_id.alias",
				"alias": bson.M{"$first": "$_id.alias"},
				"uv":    bson.M{"$sum": 1},
				"pv":    bson.M{"$sum": "$count"},
			}},
		},
		bson.D{
			primitive.E{Key: "$sort", Value: bson.M{"pv": -1, "uv": -1}},
		},
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to count Visit: %w", err)
	}
	defer cur.Close(ctx)

	var results []Record
	if err := cur.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to fetch Visit results: %w", err)
	}

	return results, nil
}
