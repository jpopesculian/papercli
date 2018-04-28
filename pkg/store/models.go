package store

import (
	bolt "github.com/coreos/bbolt"
	"github.com/jpopesculian/papercli/pkg/config"
	"log"
	"path/filepath"
)

var UPSTREAM_FOLDER_NAME_B = []byte("upstream_folder_name")
var LOCAL_FOLDER_NAME_B = []byte("local_folder_name")
var UPSTREAM_FOLDER_TREE_B = []byte("upstream_folder_tree")
var LOCAL_FOLDER_TREE_B = []byte("local_folder_tree")
var UPSTREAM_DOC_FOLDER_B = []byte("upstream_doc_folder")
var LOCAL_DOC_FOLDER_B = []byte("local_doc_folder")
var UPSTREAM_TITLE_B = []byte("upstream_title")
var LOCAL_TITLE_B = []byte("local_title")
var UPSTREAM_REVISION_B = []byte("upstream_revision")
var LOCAL_REVISION_B = []byte("local_revision")
var LAST_FETCH_B = []byte("last_fetch")
var LAST_PUSH_B = []byte("last_push")
var DOC_PATH_B = []byte("doc_path")
var FOLDER_PATH_B = []byte("folder_path")

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
		LOCAL_FOLDER_NAME_B,
		LOCAL_FOLDER_TREE_B,
		LOCAL_DOC_FOLDER_B,
		LOCAL_TITLE_B,
		LOCAL_REVISION_B,
		LAST_PUSH_B,
		DOC_PATH_B,
		FOLDER_PATH_B,
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

func (store *Store) Close() error {
	return store.db.Close()
}
