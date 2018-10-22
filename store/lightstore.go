package store

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
func View(fn func(*Txn) error) error {
	return nil
}
