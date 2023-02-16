package main

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

// CRUD with badger DB

type Book struct{
	title string
	author string
	db *badger.DB
}


func main(){
	book := create("Exodus", "Moses")
	fmt.Println(book.author, book.title)
	book.read()
	book.update("Joshua")
	book.read()
	book.delete(book.title)
	book.read()
}


func createDB() (*badger.DB, error){

	opts := badger.DefaultOptions("./book")
	db, err := badger.Open(opts)
	return db, err
}


func create(title, author string)*Book{
	newDB, err := createDB()
	if err !=  nil{
		log.Fatal(err)
	}
	err = newDB.Update(func(txn *badger.Txn) error{
		err = txn.Set([]byte(title), []byte(author))
		if err != nil{
			return err
		}
		fmt.Println("Record successfully created")
		return nil
	})

	book  := &Book{title,author,newDB}
	return book
}
func (b *Book)read(){
	var value string 	
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(b.title))
		if err != nil{
			return err
		}
		err = item.Value(func(val []byte) error{
			value = string(val)
			fmt.Println("Successful read the value")
			fmt.Println(value)
			return nil
		})
		return nil
	})
	if err != nil{
		log.Fatal(err)
	}
}
func (b *Book)update(val string){
	err := b.db.Update(func(txn *badger.Txn) error{
		err := txn.Set([]byte(b.title), []byte(val))
		if err != nil{
			return err
		}
		return nil
	})
	if err != nil{
		log.Fatal(err)
	}
}
func (b *Book) delete(key string){
	err := b.db.Update(func (txn *badger.Txn) error {
		if err := txn.Delete([]byte(key)); err == badger.ErrKeyNotFound{
			fmt.Println("Oops key not found")
		}
		return nil
	})
	if err != nil{
		log.Fatal(err)
	}
}