package worker

import (
	"github.com/stretchr/testify/assert"
	"github.com/winjeg/db-filler/schema"
	"testing"
)

func TestPerform(t *testing.T) {
	Perform()
}

func TestParallelExecute(t *testing.T) {
	tsc := schema.SchemaService.GetTableDefinition(db, conf.WorkerSettings.TableName, conf.DbSettings.DbName)
	assert.NotNil(t, tsc)
	alterSqls := tsc.GenerateAlter()
	assert.NotEmpty(t, alterSqls)
	ParallelExecute(alterSqls, 2)
}
