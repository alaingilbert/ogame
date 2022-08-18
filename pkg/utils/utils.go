package utils

import (
	"strconv"
	"strings"
)

// ParseInt ...
func ParseInt(val string) int64 {
	val = strings.Replace(val, ".", "", -1)
	val = strings.Replace(val, ",", "", -1)
	val = strings.TrimSpace(val)
	return DoParseI64(val)
}

func ToInt(buf []byte) (n int) {
	for _, v := range buf {
		n = n*10 + int(v-'0')
	}
	return
}

// I64Ptr returns a pointer to int64
func I64Ptr(v int64) *int64 {
	return &v
}

// MinInt returns the minimum int64 value
func MinInt(vals ...int64) int64 {
	min := vals[0]
	for _, num := range vals {
		if num < min {
			min = num
		}
	}
	return min
}

// MaxInt returns the minimum int64 value
func MaxInt(vals ...int64) int64 {
	max := vals[0]
	for _, num := range vals {
		if num > max {
			max = num
		}
	}
	return max
}

// Clamp ensure the value is within a range
func Clamp(val, min, max int64) int64 {
	val = MinInt(val, max)
	val = MaxInt(val, min)
	return val
}

func ParseI64(v string) (out int64, err error) {
	return strconv.ParseInt(v, 10, 64)
}

func DoParseI64(v string) (out int64) {
	out, _ = ParseI64(v)
	return
}
