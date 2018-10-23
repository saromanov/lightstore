package store

import "fmt"

// Lightstore defines main struct for db
type Lightstore struct {
	store *Store
}

// Open provides creating of lightstore object
func Open(c *Config) *Lightstore {
	return &Lightstore{
		store: newStore(c),
	}
}

// View creates new read-only transaction
func (l *Lightstore) View(fn func(*Txn) error) error {
	t := l.store.NewTransaction(false)
	err := fn(t)
	if err != nil {
		return fmt.Errorf("unable to apply transaction: %v", err)
	}
	return nil
}

// Write provides write transaction
func (l *Lightstore) Write(fn func(*Txn) error) error {
	t := l.store.NewTransaction(true)
	err := fn(t)
	if err != nil {
		return fmt.Errorf("unable to apply transaction: %v", err)
	}
	return nil
}
