package store

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/winjeg/db-filler/config"

	"database/sql"
	"errors"
	"log"
	"sync"
)

var (
	conf = config.GetConf()
	myDb *sql.DB
	once sync.Once
)

func getDb() *sql.DB {
	once.Do(func() {
		dbAddr := conf.DbSettings.GetConnectionStr()
		db, err := sql.Open(conf.DbSettings.DbType, dbAddr)
		checkErr(err)
		db.SetMaxIdleConns(int(conf.DbSettings.IdleConn))
		db.SetMaxOpenConns(int(conf.DbSettings.MaxConn))
		pingErr := db.Ping()
		if pingErr != nil {
			checkErr(err)
		}
		myDb = db
	})
	return myDb
}

func GetDb() *sql.DB {
	return getDb()
}

// panic and ends the program here, when can't connect to base db
func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// get generated id and rows affected and possible error
func GetFromResult(result sql.Result, passError error) (int64, int64, error) {
	if passError != nil || result == nil {
		return 0, 0, errors.New("result should not be null")
	}
	var finalError error = nil
	id, err := result.LastInsertId()
	if err != nil {
		finalError = err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		finalError = err
	}
	return id, affected, finalError
}

/////////////////////////////////////////////////////////////
// the following functions are  for nullable data insertion
////////////////////////////////////////////////////////////

// nullable string
func NewNullString(i interface{}) sql.NullString {
	if i == nil {
		return sql.NullString{}
	}
	if _, ok := i.(*string); !ok {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *i.(*string),
		Valid:  true,
	}
}

// nullable int
func NewNullInt(i interface{}) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	if _, ok := i.(*int64); !ok {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: *i.(*int64),
		Valid: true,
	}
}

// nullable float
func NewNullFloat(i interface{}) sql.NullFloat64 {
	if i == nil {
		return sql.NullFloat64{}
	}
	if _, ok := i.(*float64); !ok {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: *i.(*float64),
		Valid:   true,
	}
}

// nullable float
func NewNullBool(i interface{}) sql.NullBool {
	if i == nil {
		return sql.NullBool{}
	}
	if _, ok := i.(*bool); !ok {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  *i.(*bool),
		Valid: true,
	}
}
