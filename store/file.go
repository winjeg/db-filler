package store

import (
	"github.com/winjeg/db-filler/generator"

	"fmt"
	"os"
	"sync"
	"time"
)

var (
	lock sync.Mutex
	// use two lock to write file
	errorFileLock sync.Mutex
	sqlFileLock   sync.Mutex

	// none nil values should prove this is enabled
	ErrorFileStore, _ = NewFileStore(conf.ExtraSettings.ErrorFile, errorFileLock)
	SqlFileStore, _   = NewFileStore(conf.ExtraSettings.SqlFile, sqlFileLock)
)

type fileStore struct {
	File *os.File
}

// new file store, we only need to append to current file
func NewFileStore(file string, lock sync.Mutex) (*fileStore, error) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	return &fileStore{
		File: f,
	}, nil
}

// for appending sql only
// should be thread safe when doing the write work
func (fs *fileStore) Append(c string) error {
	timeInfo := fmt.Sprintf("-- %s\n", time.Now().Format(generator.TimeFormat))
	toWrite := timeInfo + c + "\n"
	lock.Lock()
	_, err := fs.File.Write([]byte(toWrite))
	lock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

// close the file
func (fs *fileStore) Close() error {
	err := fs.File.Close()
	if err != nil {
		return err
	}
	return nil
}
