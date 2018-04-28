package store

import (
	bolt "github.com/coreos/bbolt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/utils"
)

type Document struct {
	Id       Id
	Title    string
	Revision int
	Folder   Id
	Content  []byte
}

func (document *Document) FolderId() Id {
	return document.Folder
}

func (document *Document) InFolder() bool {
	return len(document.Folder) > 0
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
	b := tx.Bucket(UPSTREAM_REVISION_B)
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

func (document *Document) getUpstreamTitle(store *Store) chan string {
	result := make(chan string, 1)
	go func() {
		store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(UPSTREAM_TITLE_B)
			title := b.Get([]byte(document.Id))
			result <- string(title)
			return nil
		})
	}()
	return result
}

func (document *Document) getUpstreamFolder(store *Store) chan Id {
	result := make(chan Id, 1)
	go func() {
		store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(UPSTREAM_DOC_FOLDER_B)
			folder := b.Get([]byte(document.Id))
			result <- Id(string((folder)))
			return nil
		})
	}()
	return result
}

func (document *Document) getUpstreamContent(store *Store) chan []byte {
	result := make(chan []byte, 1)
	go func() {
		store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(LAST_FETCH_B)
			content := b.Get([]byte(document.Id))
			result <- content
			return nil
		})
	}()
	return result
}

func (document *Document) getUpstreamRevision(store *Store) chan int {
	result := make(chan int, 1)
	go func() {
		store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(UPSTREAM_REVISION_B)
			revision := b.Get([]byte(document.Id))
			result <- utils.BToI(revision)
			return nil
		})
	}()
	return result
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

func (store *Store) UpstreamDocumentById(id Id) *Document {
	document := &Document{
		Id: id,
	}
	title := document.getUpstreamTitle(store)
	revision := document.getUpstreamRevision(store)
	content := document.getUpstreamContent(store)
	folder := document.getUpstreamFolder(store)
	document.Title = <-title
	document.Revision = <-revision
	document.Content = <-content
	document.Folder = <-folder
	return document
}

func (store *Store) UpstreamDocuments() (documents chan *Document, cont chan bool) {
	documents = make(chan *Document)
	cont = make(chan bool)
	go func() {
		store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(UPSTREAM_TITLE_B)
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				cont <- true
				go func(id Id) {
					documents <- store.UpstreamDocumentById(id)
				}(Id(string(k)))
			}
			cont <- false
			return nil
		})
	}()
	return documents, cont
}
