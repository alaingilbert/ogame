package utils

import (
	"cmp"
	"github.com/PuerkitoBio/goquery"
	"iter"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// Noop ...
func Noop(_ ...any) {}

// Clamp ensure the value is within a range
func Clamp[T cmp.Ordered](val, pMin, pMax T) T {
	val = min(val, pMax)
	val = max(val, pMin)
	return val
}

func ParseI64(v string) (out int64, err error) {
	return strconv.ParseInt(v, 10, 64)
}

func DoParseI64(v string) (out int64) {
	out, _ = ParseI64(v)
	return
}

type Ints interface {
	~int64 | ~int
}

// FI64 formats any int types to string
func FI64[T Ints](v T) string {
	return strconv.FormatInt(int64(v), 10)
}

func DoCastF64(v any) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

func DoCastStr(v any) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func GetNbr(doc *goquery.Document, name string) int64 {
	div := doc.Find("div." + name)
	level := div.Find("span.level")
	level.Children().Remove()
	return ParseInt(level.Text())
}

func GetNbrShips(doc *goquery.Document, name string) int64 {
	div := doc.Find("div." + name)
	title := div.AttrOr("title", "")
	if title == "" {
		title = div.Find("a").AttrOr("title", "")
	}
	m := regexp.MustCompile(`.+\(([\d.,]+)\)`).FindStringSubmatch(title)
	if len(m) != 2 {
		return 0
	}
	return ParseInt(m[1])
}

type Equalable[T any] interface {
	Equal(other T) bool
}

func InArray[T Equalable[T]](needle T, haystack []T) bool {
	for _, el := range haystack {
		if needle.Equal(el) {
			return true
		}
	}
	return false
}

func InArr[T comparable](needle T, haystack []T) bool {
	for _, el := range haystack {
		if needle == el {
			return true
		}
	}
	return false
}

// RandChoice returns a random element from an array
func RandChoice[T any](arr []T) T {
	if len(arr) == 0 {
		panic("empty array")
	}
	return arr[rand.Intn(len(arr))]
}

// Random generates a number between min and max inclusively
func Random(min, max int64) int64 {
	if min == max {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return rand.Int63n(max-min+1) + min
}

// RandDuration generates random duration
func RandDuration(min, max time.Duration) time.Duration {
	n := Random(min.Nanoseconds(), max.Nanoseconds())
	return time.Duration(n) * time.Nanosecond
}

func randDur(min, max int64, dur time.Duration) time.Duration {
	return RandDuration(time.Duration(min)*dur, time.Duration(max)*dur)
}

// RandMs generates random duration in milliseconds
func RandMs(min, max int64) time.Duration {
	return randDur(min, max, time.Millisecond)
}

func RandFloat(min, max float64) float64 {
	if min == max {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return rand.Float64()*(max-min) + min
}

// Count2 counts element in an iter.Seq2
func Count2[K, V any](it iter.Seq2[K, V]) (out int) {
	for range it {
		out++
	}
	return
}

// Any2 return true if calling clb with any item in Seq2 return true
func Any2[K, V any](it iter.Seq2[K, V], clb func(K, V) bool) bool {
	for k, v := range it {
		if clb(k, v) {
			return true
		}
	}
	return false
}

// All return true if calling clb with all item in Seq return true
func All[V any](it iter.Seq[V], clb func(V) bool) bool {
	for v := range it {
		if !clb(v) {
			return false
		}
	}
	return true
}

func Ternary[T any](predicate bool, a, b T) T {
	if predicate {
		return a
	}
	return b
}

func TernaryOrZero[T any](predicate bool, a T) T {
	var zero T
	return Ternary(predicate, a, zero)
}

// Or return "a" if it is non-zero otherwise "b"
func Or[T comparable](a, b T) (zero T) {
	return Ternary(a != zero, a, b)
}

// RoundThousandth round value to the nearest thousandth
func RoundThousandth(n float64) float64 {
	return math.Floor(n*1000) / 1000
}

// Ptr return a pointer to v
func Ptr[T any](v T) *T { return &v }

func First[T any](a T, _ ...any) T { return a }

func Second[T any](_ any, a T, _ ...any) T { return a }

// Find looks through each value in the list, returning the first one that passes a truth test (predicate),
// or nil if no value passes the test.
// The function returns as soon as it finds an acceptable element, and doesn't traverse the entire list
func Find[T any](arr []T, predicate func(T) bool) (out *T) {
	return First(FindIdx(arr, predicate))
}

func FindIdx[T any](arr []T, predicate func(T) bool) (*T, int) {
	for i, el := range arr {
		if predicate(el) {
			return &el, i
		}
	}
	return nil, -1
}

// Deref generic deref return the zero value if v is nil
func Deref[T any](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}

// Default ...
func Default[T any](v *T, d T) T {
	if v == nil {
		return d
	}
	return *v
}
