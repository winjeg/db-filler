package worker

import (
	"github.com/winjeg/db-filler/config"
	"github.com/winjeg/db-filler/log"
	"github.com/winjeg/db-filler/schema"
	"github.com/winjeg/db-filler/store"

	"fmt"
	"sync"
)

var (
	conf   = config.GetConf()
	logger = log.GetLogger()
	db     = store.GetDb()
	wg     = sync.WaitGroup{}
)

func logError(err error) {
	if err != nil {
		logger.Error(err)
	}
}

// do the work
func Perform() {
	//  run the sql. before everything start to work
	if len(conf.ExtraSettings.Sql) > 0 {
		_, _, err := store.GetFromResult(db.Exec(conf.ExtraSettings.Sql))
		logError(err)
	}

	// only generation?
	if conf.ExtraSettings.OnlyGenerate {
		GenerateSql(conf.WorkerSettings.SqlNum)
		return
	}
	// generate sql
	sqls := GenerateSql(conf.WorkerSettings.SqlNum)

	// parallel insertion
	ParallelExecute(sqls, conf.WorkerSettings.InsertWorkerNum)

}

// generate sql  with thread num specified
func GenerateSql(num int) []string {
	// to check if table exits, if not, panic and exit the program
	_, _, err := store.GetFromResult(db.Exec(fmt.Sprintf("show create table %s", conf.WorkerSettings.TableName)))
	if err != nil {
		logger.Panic("the table you specified may not exits", err.Error())
	}
	tableDef := schema.SchemaService.GetTableDefinition(db, conf.WorkerSettings.TableName, conf.DbSettings.DbName)
	waitGroup := sync.WaitGroup{}
	threadNum := conf.WorkerSettings.GenerateWorkerNum
	waitGroup.Add(threadNum)

	genNum := num
	perNum := genNum / threadNum / conf.WorkerSettings.RowNum
	restNum := genNum % (threadNum * conf.WorkerSettings.RowNum)

	result := make([]string, 0, genNum)
	var lock sync.Mutex

	for i := 0; i < threadNum; i++ {
		go func(n int) {
			defer waitGroup.Done()
			ts := make([]string, 0, perNum)
			if n == 0 {
				// first thread takes the reset and original
				for j := 0; j < restNum+perNum; j++ {
					sql := tableDef.GenerateInsert(conf.WorkerSettings.RowNum)
					ts = append(ts, sql)
				}
			} else {
				for j := 0; j < perNum; j++ {
					sql := tableDef.GenerateInsert(conf.WorkerSettings.RowNum)
					ts = append(ts, sql)
				}
			}
			// put it to the result slice safely
			if len(ts) > 0 {
				lock.Lock()
				result = append(result, ts...)
				lock.Unlock()
			}
		}(i)
	}
	waitGroup.Wait()
	// if the sql file is set,  append the sql to this file
	if store.SqlFileStore != nil {
		for i := range result {
			err := store.SqlFileStore.Append(result[i])
			logError(err)
		}
		err = store.SqlFileStore.Close()
		logError(err)
	}
	return result
}

// parallel insert sql generation and insertion
//func ParallelizeGenAndInsert() {
//	// TODO generate and insert, using multi go routines
//}

// parallel execute sql in n routines
func ParallelExecute(sqls []string, workerNum int) {
	var waitGroup sync.WaitGroup
	sqlLen := len(sqls)
	perNum := sqlLen / workerNum
	waitGroup.Add(workerNum)
	errorSqls := make([]string, 0, sqlLen)
	var lock sync.Mutex
	for i := 0; i < workerNum; i++ {
		go func(n int) {
			defer waitGroup.Done()
			terrorSqls := make([]string, 0, perNum)
			tsqls := sqls[n*perNum : n*perNum+perNum]
			if n == workerNum-1 {
				tsqls = sqls[n*perNum:]
			}
			for _, v := range tsqls {
				_, _, err := store.GetFromResult(db.Exec(v))
				if err != nil {
					logError(err)
					terrorSqls = append(terrorSqls, v)
				}
			}
			if len(terrorSqls) > 0 {
				lock.Lock()
				errorSqls = append(errorSqls, terrorSqls...)
				lock.Unlock()
			}
		}(i)
	}
	waitGroup.Wait()
	if store.ErrorFileStore != nil {
		for i := range errorSqls {
			err := store.ErrorFileStore.Append(errorSqls[i])
			logError(err)
		}
		err := store.ErrorFileStore.Close()
		logError(err)
	}
}
