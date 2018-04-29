package store

import (
	bolt "github.com/coreos/bbolt"
	"github.com/davecgh/go-spew/spew"
)

func (store *Store) FetchFirstId() Id {
	idC := make(chan Id, 1)
	store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UPSTREAM_TITLE_B)
		c := b.Cursor()
		k, _ := c.First()
		if k == nil {
			idC <- Id("")
		} else {
			idC <- Id(k)
		}
		return nil
	})
	return <-idC
}

func (store *Store) Stats() {
	spew.Dump(store.db.IsReadOnly())
	spew.Dump(store.db.Stats())
}
