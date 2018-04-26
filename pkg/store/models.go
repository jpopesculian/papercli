package store

import (
	bolt "github.com/coreos/bbolt"
	"github.com/jpopesculian/papercli/pkg/utils"
	"log"
)

var UPSTREAM_FOLDER_NAME_B = []byte("upstream_folder_name")
var UPSTREAM_FOLDER_TREE_B = []byte("upstream_folder_tree")
var UPSTREAM_DOC_FOLDER_B = []byte("upstream_doc_folder")
var UPSTREAM_TITLE_B = []byte("upstream_title")
var UPSTREAM_REVISION_B = []byte("upstream_revision")
var LAST_FETCH_B = []byte("last_fetch")

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
		UPSTREAM_DOC_FOLDER_B,
		UPSTREAM_TITLE_B,
		UPSTREAM_REVISION_B,
		LAST_FETCH_B,
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

func (document *Document) saveUpstreamTitle(tx *bolt.Tx) error {
	b := tx.Bucket(UPSTREAM_TITLE_B)
	err := b.Put([]byte(document.Id), []byte(document.Title))
	if err != nil {
		return err
	}
	return nil
}

func (document *Document) saveUpstreamFolder(tx *bolt.Tx) error {
	b := tx.Bucket(UPSTREAM_DOC_FOLDER_B)
	err := b.Put([]byte(document.Id), []byte(document.Folder))
	if err != nil {
		return err
	}
	return nil
}

func (document *Document) saveUpstreamRevision(tx *bolt.Tx) error {
	b := tx.Bucket(UPSTREAM_DOC_FOLDER_B)
	err := b.Put([]byte(document.Id), utils.IToB(document.Revision))
	if err != nil {
		return err
	}
	return nil
}

func (document *Document) saveLastFetch(tx *bolt.Tx) error {
	b := tx.Bucket(LAST_FETCH_B)
	err := b.Put([]byte(document.Id), document.Content)
	if err != nil {
		return err
	}
	return nil
}

func (document *Document) saveUpstream(tx *bolt.Tx) error {
	if err := document.saveUpstreamTitle(tx); err != nil {
		return err
	}
	if err := document.saveUpstreamFolder(tx); err != nil {
		return err
	}
	if err := document.saveUpstreamRevision(tx); err != nil {
		return err
	}
	return document.saveLastFetch(tx)
}

func (store *Store) SaveUpstreamDocuments(documents []Document) error {
	err := store.db.Batch(func(tx *bolt.Tx) error {
		for _, document := range documents {
			err := document.saveUpstream(tx)
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
	if err := folder.saveUpstreamName(tx); err != nil {
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
