package main

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"gopkg.in/vmihailenco/msgpack.v2"
)

/**
*  Created by Galileo on 27/5/17.
 */

const (
	_RW_MODE = 0600
	_BUCKET  = "geymsla"
)

// Geymsla represents the k-v store
type Geymsla struct {
	store *bolt.DB
}

var bucketName = []byte(_BUCKET)

// Open a k-v store. Path is the fully qualified path to the BoltDB file
func Open(path string) (*Geymsla, error) {
	store, err := bolt.Open(path, _RW_MODE, &bolt.Options{
		Timeout: time.Second * time.Duration(1),
	})
	if err != nil {
		return nil, err
	}
	err = store.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Geymsla{store: store}, nil
}

// Set a key value pair in Geymsla
func (g *Geymsla) Set(key string, value interface{}) error {
	if value == nil {
		return errors.New("Bad value: Value cannot be nil")
	}
	_bytes, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	return g.store.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), _bytes)
	})
}

// Get a value from the key from Geymsla
func (g *Geymsla) Get(key string) (interface{}, error) {
	var value interface{}
	err := g.store.View(func(tx *bolt.Tx) error {
		k, v := tx.Bucket(bucketName).Cursor().Seek([]byte(key))
		if k == nil || string(k) != key {
			return errors.New("Bad key: key not found.")
		} else {
			if err := msgpack.Unmarshal([]byte(v), &value); err != nil {
				return err
			}
			return nil
		}
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Delete the give key value pair from Geymsla
func (g *Geymsla) Delete(key string) error {
	return g.store.Update(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(bucketName).Cursor()
		if k, _ := cursor.Seek([]byte(key)); k != nil || string(k) != key {
			return errors.New("Bad key: key not found.")
		} else {
			return cursor.Delete()
		}
	})
}

// Close the store
func (g *Geymsla) Close() error {
	return g.store.Close()
}
