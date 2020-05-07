// Package nodes2wayblotdb stores `fromNodeID,toNodeID -> wayID` mapping in blotdb.
package nodes2wayblotdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/serho/osrm-backend/integration/util/waysnodes"
	"github.com/golang/glog"

	bolt "go.etcd.io/bbolt"
)

// DB stores `fromNodeID,toNodeID -> wayID` in blotdb.
type DB struct {
	db *bolt.DB
}

const (
	nodes2WayBucket = "nodes2way_bucket"
)

var (
	errEmptyDB     = errors.New("empty db")
	errKeyNotFound = errors.New("key not found")
)

// Open creates/opens a DB structure to store the nodes2way mapping.
func Open(dbFilePath string, readOnly bool) (*DB, error) {
	var db DB

	options := bolt.Options{

		// Default Bolt Options
		Timeout:      0,
		NoGrowSync:   false,
		FreelistType: bolt.FreelistArrayType,

		// set readonly
		ReadOnly: readOnly,
	}

	var err error
	db.db, err = bolt.Open(dbFilePath, 0666, &options)
	if err != nil {
		return nil, err
	}

	if readOnly {
		return &db, nil
	}

	// to improve write performance, but will manually sync before close
	db.db.NoSync = true
	db.db.NoFreelistSync = true

	// for write, make sure bucket available
	err = db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(nodes2WayBucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// Close releases all database resources.
func (db *DB) Close() error {
	if db.db == nil {
		return errEmptyDB
	}

	if !db.db.IsReadOnly() {
		startTime := time.Now()
		if err := db.db.Sync(); err != nil {
			glog.Error(err)
			//return err	// don't return since we still hope the Close can be called
		}
		glog.V(1).Infof("Flush DB to disk file %s takes %f seconds.", db.db.Path(), time.Now().Sub(startTime).Seconds())
	}

	return db.db.Close()
}

// Statistics returns internal statistics information.
func (db *DB) Statistics() string {
	if db.db == nil {
		return ""
	}

	s := struct {
		DBStat     string `json:"dbstat"`
		BucketStat string `json:"bucketstat"`
	}{}

	s.DBStat = fmt.Sprintf("%+v", db.db.Stats())
	if err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodes2WayBucket))

		s.BucketStat = fmt.Sprintf("%+v", b.Stats()) // bucket stats
		return nil
	}); err != nil {
		glog.Error(err)
		return err.Error()
	}

	b, err := json.Marshal(&s)
	if err != nil {
		glog.Error(err)
		return err.Error()
	}

	return string(b)
}

// Write writes wayID and its nodeIDs into cache or storage.
// wayID: is undirected when input, so will always be positive.
func (db *DB) Write(wayID int64, nodeIDs []int64) error {
	if db.db == nil {
		return errEmptyDB
	}
	if wayID < 0 {
		return fmt.Errorf("expect undirected wayID")
	}
	if len(nodeIDs) < 2 {
		return fmt.Errorf("at least 2 nodes for a way")
	}

	err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodes2WayBucket))

		for i := 0; i < len(nodeIDs)-1; i++ {
			if err := b.Put(key(nodeIDs[i], nodeIDs[i+1]), value(wayID)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// BatchWrite writes wayID and its nodeIDs into cache or storage by batching.
func (db *DB) BatchWrite(wayNodes []waysnodes.WayNodes) error {
	if db.db == nil {
		return errEmptyDB
	}

	// Create several keys in a transaction.
	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	b := tx.Bucket([]byte(nodes2WayBucket))

	for _, wn := range wayNodes {
		for i := 0; i < len(wn.NodeIDs)-1; i++ {
			if err := b.Put(key(wn.NodeIDs[i], wn.NodeIDs[i+1]), value(wn.WayID)); err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// QueryWay queries directed wayID by fromNodeID,toNodeID pair.
// returned wayID: positive means travel forward following the fromNodeID,toNodeID sequence, negative means backward
func (db *DB) QueryWay(fromNodeID, toNodeID int64) (int64, error) {
	if db.db == nil {
		return 0, errEmptyDB
	}

	var wayID int64
	if err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodes2WayBucket))

		v := b.Get(key(fromNodeID, toNodeID))
		if v != nil {
			wayID = parseValue(v)
			return nil
		}

		// try again on backward
		v = b.Get(key(toNodeID, fromNodeID))
		if v == nil {
			return errKeyNotFound // both forward and backward not found
		}
		wayID = parseValue(v)
		wayID = -wayID
		return nil
	}); err != nil {
		return 0, err
	}

	return wayID, nil
}

// QueryWays queries directed wayIDs by nodeIDs.
// `len(wayIDs) == len(nodeIDs)-1` since each way have to be decided by traveling from one node to another.
// returned wayIDs: positive means travel forward following the nodeIDs sequence, negative means backward
func (db *DB) QueryWays(nodeIDs []int64) ([]int64, error) {
	if db.db == nil {
		return nil, errEmptyDB
	}

	var wayIDs []int64
	if err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodes2WayBucket))

		for i := 0; i < len(nodeIDs)-1; i++ {
			v := b.Get(key(nodeIDs[i], nodeIDs[i+1]))
			if v != nil {
				wayIDs = append(wayIDs, parseValue(v))
				continue
			}

			// try again on backward
			v = b.Get(key(nodeIDs[i+1], nodeIDs[i]))
			if v == nil {
				return fmt.Errorf("%v: %d,%d", errKeyNotFound, nodeIDs[i], nodeIDs[i+1])
			}
			wayID := parseValue(v)
			wayIDs = append(wayIDs, -wayID)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return wayIDs, nil
}
