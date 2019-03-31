package store

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestTableName      = "test_1"
	TestCreateTableSQL = `
	 CREATE TABLE %s (
		  id int(11) NOT NULL COMMENT 'the primary key',
		  name varchar(64) NOT NULL DEFAULT 'unknown' COMMENT 'name',
		  age tinyint(4) NOT NULL DEFAULT '0' COMMENT 'age',
		  deleted bit(1) NOT NULL DEFAULT b'0',
		  note varchar(255) DEFAULT NULL COMMENT 'desc',
		  created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
		  updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last updated',
		  PRIMARY KEY (id),
		  UNIQUE KEY idx_id (id),
		  KEY idx_name (name),
		  KEY idx_name_deleted (name) USING BTREE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='table for person information'
	`
)

func TestGetDb(t *testing.T) {
	db := GetDb()
	db.QueryRow("SELECT 1")
	i, n, err := GetFromResult(db.Exec("SELECT 1"))
	assert.Nil(t, err)
	assert.True(t, i == 0)
	assert.True(t, n == 0)
	_, _, err = GetFromResult(db.Exec("WHEN 1"))
	assert.NotNil(t, err)
	_, _, err = GetFromResult(db.Exec(fmt.Sprintf(TestCreateTableSQL, TestTableName)))
	assert.Nil(t, err)
	_, _, err = GetFromResult(db.Exec(fmt.Sprintf("DROP table %s", TestTableName)))
	assert.Nil(t, err)
}

func TestNullable(t *testing.T) {
	// bool
	b1 := NewNullBool(nil)
	assert.False(t, b1.Valid)
	v1 := false
	b2 := NewNullBool(&v1)
	assert.True(t, b2.Valid)
	b3 := NewNullBool(10)
	assert.False(t, b3.Valid)

	// string
	s1 := NewNullString(nil)
	assert.False(t, s1.Valid)
	v2 := "test"
	s2 := NewNullString(&v2)
	assert.True(t, s2.Valid)
	s3 := NewNullString(2)
	assert.False(t, s3.Valid)

	// int
	i1 := NewNullInt(nil)
	assert.False(t, i1.Valid)
	v3 := int64(1)
	i2 := NewNullInt(&v3)
	assert.True(t, i2.Valid)
	i3 := NewNullInt(1.2)
	assert.False(t, i3.Valid)

	// float
	f1 := NewNullFloat(nil)
	assert.False(t, f1.Valid)
	v4 := float64(1.3)
	f2 := NewNullFloat(&v4)
	assert.True(t, f2.Valid)
	f3 := NewNullFloat(12)
	assert.False(t, f3.Valid)
}
