package store

import (
	"errors"
	"fmt"
)

var errFnIsNotDefined = errors.New("function is not defined")

// Lightstore defines main struct for db
type Lightstore struct {
	store           *Store
	errTransactions []string
}

// Open provides creating of lightstore object
func Open(c *Config) (*Lightstore, error) {
	store, err := newStore(c)
	if err != nil {
		return nil, err
	}
	return &Lightstore{
		store: store,
	}, nil
}

// OpenStrict provides creating of lightstore object only with db path
func OpenStrict(path string) (*Lightstore, error) {
	return Open(&Config{LoadPath: path})
}

// IsCreated retruns true if Lightstore was initialized
func (l *Lightstore) IsCreated() bool {
	return l.store.IsCreated()
}

// View creates new read-only transaction
func (l *Lightstore) View(fn func(*Txn) error) error {
	if fn == nil {
		return errFnIsNotDefined
	}
	t := l.store.NewTransaction(false)
	err := fn(t)
	if err != nil {
		l.errTransactions = append(l.errTransactions, t.ID())
		return fmt.Errorf("unable to apply transaction: %v", err)
	}
	return nil
}

// Write provides write transaction
func (l *Lightstore) Write(fn func(*Txn) error) error {
	if fn == nil {
		return errFnIsNotDefined
	}
	t := l.store.NewTransaction(true)
	err := fn(t)
	if err != nil {
		l.errTransactions = append(l.errTransactions, t.ID())
		return fmt.Errorf("unable to apply transaction: %v", err)
	}
	return nil
}

// Close provides closing of Lightstore session
func (l *Lightstore) Close() error {
	if l != nil && l.store != nil {
		l.store.Close()
	}
	return nil
}

// getStore returns store engine of lightstore
func (l *Lightstore) getStore() *Store {
	return l.store
}
