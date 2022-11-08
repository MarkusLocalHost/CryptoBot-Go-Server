package main

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("./tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte("shiba inu")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				log.Printf("key=%s, value=%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
