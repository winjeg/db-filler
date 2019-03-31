package store

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestFileStore_Append(t *testing.T) {
	fs, err := NewFileStore("a.txt", &lock)
	assert.NotNil(t, fs)
	assert.Nil(t, err)

	err = fs.Append("SELECT 1")
	assert.Nil(t, err)
	for i := 0; i < 100; i++ {
		go func(i int) {
			err := fs.Append("SELECT 1" + fmt.Sprintf("# %d", i))
			assert.Nil(t, err)
		}(i)
	}
	time.Sleep(time.Second * 2)
	err = fs.Close()
	assert.Nil(t, err)
	err = os.Remove("a.txt")
	assert.Nil(t, err)
}
