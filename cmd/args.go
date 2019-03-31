package cmd

import (
	"flag"
	"fmt"
)

// dealing with command line args
// to make this tool usable from command line
// the args here has higher priority over the args in conf.yaml

const (
	defaultNum  = 10
	defaultPort = 3306
)

var (
	host               string // -h host
	port               int    // -p port
	username           string // -u user
	auth               string // -a password
	dbName             string // -d db the database name
	tableName          string // table name for generating sql and to insert generated sql
	sqlNum             int    // the number of sql to be generated
	rowNumPerSql       int    // how many rows in one sql
	generateWorkerNum  int    // how many go routines  to use to generate sql, default 1
	insertionWorkerNum int    // how many go routines used to run insert, default 1
	ddlWorkerNum       int    // how many go routines to use to do ddl, default 1
	sql                string // just like mysql client -e, use it to execute a sql, to do some initialization
	errorFile          string // -ef save error sql to a file
	sqlFile            string // save generated sql to a file
	onlyGenerate       bool   // only generate the sql, don't do insertion

)

// dealing with command line args
func init() {

	// basic connection information
	flag.StringVar(&host, "H", "", "address of a mysql server")
	flag.IntVar(&port, "p", defaultPort, "port of the mysql server")
	flag.StringVar(&username, "u", "", "username of the mysql server")
	flag.StringVar(&auth, "a", "", "password/auth of the mysql server")
	flag.StringVar(&dbName, "d", "", "database name  of the mysql server to operate on")

	// must specify a table name
	flag.StringVar(&tableName, "t", "", "the table name to be used to generate sql and to be inserted")

	// args with default values
	flag.IntVar(&sqlNum, "n", defaultNum, "the number of sql to be generated and inserted")
	flag.IntVar(&rowNumPerSql, "rn", 1, "the number of rows in one sql")
	flag.IntVar(&generateWorkerNum, "gwn", defaultNum, "the number of routines to generate sql ")
	flag.IntVar(&insertionWorkerNum, "iwn", defaultNum, "the number of routines to perform insertion")
	flag.IntVar(&ddlWorkerNum, "dwn", defaultNum, "the number of routines to perform ddl")

	// optional args
	flag.StringVar(&sql, "e", "", "the sql to be executed before everything starts")
	flag.StringVar(&errorFile, "ef", "", "the file where to put the error encountered when executing sql, if nothing specified, the error will be only in stdout ")
	flag.StringVar(&sqlFile, "sf", "", "the file to save generated sql, if none specified, no sql will be saved ")
	flag.BoolVar(&onlyGenerate, "og", false, "only generate the sql, don't do insertion")

	// override usage to show user defined content
	flag.Usage = usage

	flag.Parse()
}

func usage() {
	fmt.Println("db-filler usage:")
	fmt.Println("db-filler -H localhost -p 3306 -a 123456 -d test -t person -n 1000")
	fmt.Println("======================================================")
	flag.PrintDefaults()
}

// get database name
func GetDbName() string {
	return dbName
}

// get table name
func GetTableName() string {
	return tableName
}

// get sql to be executed before generation and insertion
func GetSql() string {
	return sql
}

// get file to store error sql
func GetErrorFile() string {
	return errorFile
}

// get file to store sql
func GetSqlFile() string {
	return sqlFile
}

// get the number of sql to be generated
func GetGenSqlNum() int {
	return sqlNum
}

// only generate sql don't do insert
func OnlyGenerate() bool {
	return onlyGenerate
}

// get the number of routines to do sql generation
func GetGenerateWorkerNum() int {
	return generateWorkerNum
}

// get the number of routines to do insertion
func GetInsertionWorkerNum() int {
	return insertionWorkerNum
}

// get the number of routines to run ddl
func GetDDLWorkerNum() int {
	return ddlWorkerNum
}

// get mysql port
func GetPort() int {
	return port
}

// get mysql username
func GetUserName() string {
	return username
}

// get mysql password
func GetPassword() string {
	return auth
}

//  get mysql host addr
func GetHost() string {
	return host
}

// get row num of a sql
func GetRowNumPerSql() int {
	return rowNumPerSql
}
