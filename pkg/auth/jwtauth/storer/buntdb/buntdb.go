package buntdb

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/tidwall/buntdb"
)

type Store struct {
	db *buntdb.DB
}

func NewStore(path string) (*Store, error) {
	if path != ":memory" {
		os.MkdirAll(filepath.Dir(path), 0777)
	}
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	return &Store{
		db: db,
	}, nil
}

func (a *Store) Set(ctx context.Context, tokenString string, expireation time.Duration) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		var opts *buntdb.SetOptions
		if expireation > 0 {
			opts = &buntdb.SetOptions{Expires: true, TTL: expireation}
		}
		_, _, err := tx.Set(tokenString, "1", opts)
		return err
	})
}

func (a *Store) Delete(ctx context.Context, tokenString string) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(tokenString)
		if err != nil && err != buntdb.ErrIndexExists {
			return err
		}
		return nil
	})
}

// Check ...
func (a *Store) Check(ctx context.Context, tokenString string) (bool, error) {
	var exists bool
	err := a.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(tokenString)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		exists = val == "1"
		return nil
	})
	return exists, err
}

// Close ...
func (a *Store) Close() error {
	return a.db.Close()
}
