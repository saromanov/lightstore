//This package provides create and manage db data

package store

type DB struct {
	//db name
	dbtitle string
	//Limit of number of records in this db
	limit int
	//All data is immutable, no overwrites
	immutable bool
	//Check if this db is active
	isactive bool
	//Number of data object
	datacount int
	//Main store of db
	mainstore interface{}
}

//CreateNewDB provides new DB object
func CreateNewDB(name string) *DB {
	return &DB{name, -1, false, true, 0, NewDict()}
}

//SetLimit provides set maximum number of data in this db
func (db *DB) SetLimit(limit int) {
	db.limit = limit
}

//Stop provides stopping this db and can't receive new updates
func (db *DB) Stop() {
	db.isactive = false
}
