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
	light := store.Open(nil)
	defer light.Close()
}
```
This creates a new database with default config

## TODO
Indexing
Distributed (Consensus, Failure detection)
Saving data on disk
More rich API
Documentation
Tests

