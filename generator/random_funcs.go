package generator

import (
	"fmt"
	"github.com/winjeg/go-commons/str"
	"math"
	"math/rand"
	"time"
)

// using a rand.Seed(n) function slows down all random functions
// so in every function,  we don't use a seed

const (
	// time
	TimeFormat    = "2006-01-02 15:04:05"
	timeStartDate = "1970-01-01 00:00:01"
	timeEndDate   = "2038-01-19 03:14:07"

	// cjk section of unicode charset
	cjkStart = 0x4E00
	cjkStop  = 0x9FBF
)

var (
	// time related vars, max, min time and duration
	timeStart, _ = time.Parse(TimeFormat, timeStartDate)
	timeEnd, _   = time.Parse(TimeFormat, timeEndDate)
	maxDuration  = uint64(timeEnd.Sub(timeStart).Nanoseconds())
)

// random int, returns possibly negative value or positive value
// including the limit, and negative values, the arg limit should be positive
func RandomInt(limit int64) int64 {
	if RandomBool() {
		return -rand.Int63n(limit) - 2
	}
	return rand.Int63n(limit) + 1
}

// random unsigned int, positive value only
func RandomUInt(limit uint64) uint64 {
	r := rand.Uint64()
	if r > limit {
		return r % limit
	}
	return r
}

// random bool
func RandomBool() bool {
	return rand.Intn(2) == 1
}

// return a random double, limit is the max value
// part return 0 for illegal values
func RandomFloat(n, f int) string {
	if n <= 0 {
		n = 0
	}
	if f <= 0 {
		f = 0
	}
	partN := make([]byte, n)
	partF := make([]byte, f)
	for i := 0; i < n; i++ {
		partN[i] = byte(rand.Intn(10) + 48)
	}
	for i := 0; i < f; i++ {
		partF[i] = byte(rand.Intn(10) + 48)
	}
	if n == 0 && f == 0 {
		return "0.0"
	}
	if n <= 0 {
		return "0." + string(partF)
	}
	if f <= 0 {
		return string(partN) + ".0"
	}
	return string(partN) + "." + string(partF)
}

// return a random legal time both for mysql type datetime and timestamp
func RandomTime() time.Time {
	return timeStart.Add(time.Duration(RandomUInt(maxDuration)))
}

// random bits,  n should at most be 64
func RandomBits(n int) string {
	num := RandomUInt(math.MaxUint64) + math.MaxUint64>>2
	bits := fmt.Sprintf("%b", num)
	return bits[:n]
}

// random CJK in UTF8
// for utf-8 is a mutable length charset
// so the result length may not actually equals to the size specified
func RandomCJK(size int) string {
	result := make([]rune, size)
	for i := range result {
		result[i] = rune(RandIntSection(cjkStart, cjkStop))
	}
	return string(result)
}

// random ascii
func RandomASCII(size int) string {
	result := make([]byte, size)
	for i := range result {
		result[i] = byte(RandIntSection(0, 128))
	}
	return string(result)
}

// random readable string
func RandomReadable(size int) string {
	return str.RandomStrWithSpecialChars(size)
}

// random int section
func RandIntSection(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}
