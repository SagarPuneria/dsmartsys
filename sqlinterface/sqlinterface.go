package SqlInterface

import (
	"database/sql"
	"sync"

	// Register mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// MySqldb contains db and dbRWMutex
type MySqldb struct {
	db        *sql.DB
	dbRWMutex sync.RWMutex
}

// CreateDataBase login to mysql db, create db, table if not exists
func CreateDataBase(DNS string, quries ...string) (*MySqldb, error) {
	var dbConn = new(MySqldb)
	dbConn.dbRWMutex.Lock()
	defer dbConn.dbRWMutex.Unlock()
	var err error
	dbConn.db, err = sql.Open("mysql", DNS)
	if err != nil {
		return nil, err
	}
	for _, query := range quries {
		_, err := dbConn.db.Exec(query)
		if err != nil {
			dbConn.Close()
			return dbConn, err
		}
	}
	return dbConn, nil
}

// Close the db
func (DBObject *MySqldb) Close() {
	DBObject.db.Close()
}

// ExecuteQuery excute the given query
func (DBObject *MySqldb) ExecuteQuery(strQuery string) error {
	DBObject.dbRWMutex.Lock()
	defer DBObject.dbRWMutex.Unlock()

	_, err := DBObject.db.Exec(strQuery)
	return err
}

// SelectQuery execute select query
func (DBObject *MySqldb) SelectQuery(strQuery string) (*sql.Rows, error) {
	DBObject.dbRWMutex.Lock()
	defer DBObject.dbRWMutex.Unlock()

	rows, err := DBObject.db.Query(strQuery)

	return rows, err
}
