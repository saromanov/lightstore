# lightstore
[![Go Report Card](https://goreportcard.com/badge/github.com/saromanov/lightstore)](https://goreportcard.com/report/github.com/saromanov/lightstore)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d49a6728569744c08db82a534b28821f)](https://www.codacy.com/app/saromanov/lightstore?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=saromanov/lightstore&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/saromanov/lightstore.svg?branch=master)](https://travis-ci.org/saromanov/lightstore)
[![Coverage Status](https://coveralls.io/repos/github/saromanov/lightstore/badge.svg?branch=master)](https://coveralls.io/github/saromanov/lightstore?branch=master)

Key-Value store
in Progress

## Table of Contents
* [Getting Started](#getting-started)
    + [Installing](#installing)
    + [Create database](#create-database)
    + [Write transactions](#write-transactions)
    + [Read from DB](#read-from-db)
    + [Iterator](#iterator)
    + [Stapshots](#snapshots)
* [Features](#features)

## Getting Started

### Installing

```sh
$ go get github.com/saromanov/lightstore/...
```

### Create Database
Easy steps for create a new database
```go
package main

import "github.com/saromanov/lightstore/store"

func main() {
    light, err := store.Open(nil)
    if err != nil {
        panic(err)
    }
	defer light.Close()
}
```
This creates a new database with default config

Or create database only with path to db file
```go
light, err := store.OpenStrict("db.db")
if err != nil {
	panic(err)
}
defer light.Close()
```

### Write transactions

For make new transaction, need to open new transaction and then commit changes
```go
err = light.Write(func(txn *store.Txn) error {
		for i := 0; i < 20; i++ {
			err := txn.Set([]byte(fmt.Sprintf("%d", i)), []byte(fmt.Sprintf("bar+%d", i)))
			if err != nil {
				return err
			}
		}
		return txn.Commit()
	})
	if err != nil {
		log.Fatalf("unable to write data: %v", err)
    }
```

### Iterator

```go
light.View(func(txn *store.Txn) error {
		it, _ := txn.NewIterator(store.IteratorOptions{})
		for it.First(); it.Valid(); it.Next() {
			itm := it.Item()
			fmt.Println(string(itm.Key()), string(itm.Value()))
		}
		return nil
})
```

