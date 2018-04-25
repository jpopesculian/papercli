package store

import (
	bolt "github.com/coreos/bbolt"
	"log"
)

var UPSTREAM_FOLDER_NAME_B = []byte("upstream_folder_name")
var UPSTREAM_FOLDER_TREE_B = []byte("upstream_folder_tree")

type Id string

type Store struct {
	db *bolt.DB
}

type Document struct {
	Id       Id
	Title    string
	Revision int
	Folder   Id
	Content  []byte
}

type Folder struct {
	Id     Id
	Name   string
	Parent Id
}

func NewStore() *Store {
	db, err := bolt.Open("paper.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	store := &Store{
		db: db,
	}
	err = store.createBuckets()
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func (store *Store) createBuckets() error {
	buckets := [][]byte{
		UPSTREAM_FOLDER_NAME_B,
		UPSTREAM_FOLDER_TREE_B,
	}
	err := store.db.Batch(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (folder *Folder) saveUpstreamTree(tx *bolt.Tx) error {
	b := tx.Bucket(UPSTREAM_FOLDER_TREE_B)
	err := b.Put([]byte(folder.Id), []byte(folder.Parent))
	if err != nil {
		return err
	}
	return nil
}

func (folder *Folder) saveUpstreamName(tx *bolt.Tx) error {
	b := tx.Bucket(UPSTREAM_FOLDER_NAME_B)
	err := b.Put([]byte(folder.Id), []byte(folder.Name))
	if err != nil {
		return err
	}
	return nil
}

func (folder *Folder) saveUpstream(tx *bolt.Tx) error {
	err := folder.saveUpstreamName(tx)
	if err != nil {
		return err
	}
	return folder.saveUpstreamTree(tx)
}

func (store *Store) SaveUpstreamFolders(folders []Folder) error {
	err := store.db.Batch(func(tx *bolt.Tx) error {
		for _, folder := range folders {
			err := folder.saveUpstream(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
