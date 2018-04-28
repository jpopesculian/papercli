package store

import (
	bolt "github.com/coreos/bbolt"
)

type Folder struct {
	Id     Id
	Name   string
	Parent Id
}

func (folder *Folder) FolderId() Id {
	return folder.Parent
}

func (folder *Folder) InFolder() bool {
	return len(folder.Parent) > 0
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

func (folder *Folder) saveLocalTree(tx *bolt.Tx) error {
	b := tx.Bucket(LOCAL_FOLDER_TREE_B)
	err := b.Put([]byte(folder.Id), []byte(folder.Parent))
	if err != nil {
		return err
	}
	return nil
}

func (folder *Folder) saveLocalName(tx *bolt.Tx) error {
	b := tx.Bucket(LOCAL_FOLDER_NAME_B)
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

func (folder *Folder) saveLocal(tx *bolt.Tx) error {
	if err := folder.saveLocalName(tx); err != nil {
		return err
	}
	return folder.saveLocalTree(tx)
}

func (folder *Folder) getUpstreamName(store *Store) chan string {
	result := make(chan string, 1)
	go func() {
		store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(UPSTREAM_FOLDER_NAME_B)
			name := b.Get([]byte(folder.Id))
			result <- string((name))
			return nil
		})
	}()
	return result
}

func (folder *Folder) getUpstreamParent(store *Store) chan Id {
	result := make(chan Id, 1)
	go func() {
		store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(UPSTREAM_FOLDER_TREE_B)
			parent := b.Get([]byte(folder.Id))
			result <- Id(string((parent)))
			return nil
		})
	}()
	return result
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

func (store *Store) SaveFolderPath(id Id, path string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(FOLDER_PATH_B)
		err := b.Put([]byte(path), []byte(id))
		if err != nil {
			return err
		}
		return nil
	})
}

func (store *Store) SaveLocalFolders(folders []Folder) error {
	err := store.db.Batch(func(tx *bolt.Tx) error {
		for _, folder := range folders {
			err := folder.saveLocal(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (store *Store) UpstreamFolderById(id Id) *Folder {
	folder := &Folder{
		Id: id,
	}
	name := folder.getUpstreamName(store)
	parent := folder.getUpstreamParent(store)
	folder.Name = <-name
	folder.Parent = <-parent
	return folder
}
