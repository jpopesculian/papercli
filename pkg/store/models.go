package store

import (
	bolt "github.com/coreos/bbolt"
	"github.com/jpopesculian/papercli/pkg/config"
	"log"
	"path/filepath"
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

type FolderEntity interface {
	FolderId() Id
	InFolder() bool
}

func getDbPath(options *config.CliOptions) (string, error) {
	dir, err := options.ConfigDir()
	if err != nil {
		return "", nil
	}
	return filepath.Join(dir, "paper.db"), nil
}

func NewStore(options *config.CliOptions) *Store {
	path, err := getDbPath(options)
	if err != nil {
		log.Fatal(err)
	}
	db, err := bolt.Open(path, 0600, nil)
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
