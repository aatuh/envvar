package lazy

import (
	"sync"
	"time"

	"github.com/aatuh/envvar/v2/getters"
)

// onceVal is a struct that contains a once and a value.
type onceVal[T any] struct {
	once sync.Once
	val  T
}

// LazyString returns a function that returns the value of the environment
// variable with the given key.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key.
func LazyString(key string) func() string {
	var o onceVal[string]
	return func() string {
		o.once.Do(func() { o.val = getters.MustGet(key) })
		return o.val
	}
}

// LazyBool returns a function that returns the value of the environment
// variable with the given key as a boolean.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a boolean.
func LazyBool(key string) func() bool {
	var o onceVal[bool]
	return func() bool {
		o.once.Do(func() { o.val = getters.MustGetBool(key) })
		return o.val
	}
}

// LazyInt returns a function that returns the value of the environment
// variable with the given key as an integer.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as an integer.
func LazyInt(key string) func() int {
	var o onceVal[int]
	return func() int {
		o.once.Do(func() { o.val = getters.MustGetInt(key) })
		return o.val
	}
}

// LazyInt64 returns a function that returns the value of the environment
// variable with the given key as an int64.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as an int64.
func LazyInt64(key string) func() int64 {
	var o onceVal[int64]
	return func() int64 {
		o.once.Do(func() { o.val = getters.MustGetInt64(key) })
		return o.val
	}
}

// LazyFloat64 returns a function that returns the value of the environment
// variable with the given key as a float64.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a float64.
func LazyFloat64(key string) func() float64 {
	var o onceVal[float64]
	return func() float64 {
		o.once.Do(func() { o.val = getters.MustGetFloat64(key) })
		return o.val
	}
}

// LazyDuration returns a function that returns the value of the environment
// variable with the given key as a duration.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a duration.
func LazyDuration(key string) func() time.Duration {
	var o onceVal[time.Duration]
	return func() time.Duration {
		o.once.Do(func() { o.val = getters.MustGetDuration(key) })
		return o.val
	}
}

// LazyStringSlice returns a function that returns the value of the environment
// variable with the given key as a slice of strings.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a slice of strings.
func LazyStringSlice(key string) func() []string {
	var o onceVal[[]string]
	return func() []string {
		o.once.Do(func() { o.val = getters.MustGetStringSlice(key) })
		return o.val
	}
}

// LazyTyped returns a function that returns the value of the environment
// variable with the given key as a typed value.
//
// Parameters:
//   - key: The key to get.
//   - conv: The converter function.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a typed value.
func LazyTyped[T any](key string, conv func(string) (T, error)) func() T {
	var o onceVal[T]
	return func() T {
		o.once.Do(func() { o.val = getters.MustGetTyped(key, conv) })
		return o.val
	}
}
