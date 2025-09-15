package envvar

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var once sync.Once

// EnvVar represents an env var with a name and value.
type EnvVar struct {
	Key   string
	Value any
}

// MustGet returns the value of an environment variable.
// It panics if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - The value of the environment variable.
func MustGet(name string) string {
	MustLoadEnvVars(nil)

	value := os.Getenv(name)
	if len(value) == 0 {
		panic(fmt.Sprintf("no env var: %q", name))
	}

	return value
}

// GetOr returns the value of an environment variable.
// It returns the default value if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//   - defaultValue: The default value to return if the environment variable is not set.
//
// Returns:
//   - The value of the environment variable.
func GetOr(name string, defaultValue string) string {
	value, err := GetOrErr(name)
	if err != nil {
		return defaultValue
	}
	return value
}

// Get returns the value of an environment variable.
// It returns an empty string if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - The value of the environment variable.
func Get(name string) string {
	value, err := GetOrErr(name)
	if err != nil {
		return ""
	}
	return value
}

// GetOrErr returns the value of an environment variable.
// It returns an error if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - The value of the environment variable.
//   - An error if the environment variable is not set.
func GetOrErr(name string) (string, error) {
	MustLoadEnvVars(nil)

	value := os.Getenv(name)
	if len(value) == 0 {
		return "", fmt.Errorf("no env var %q", name)
	}
	return value, nil
}

// MustGetInt returns the value of an environment variable as an int.
// It panics if the environment variable is not set or cannot be converted to an
// int.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - The value of the environment variable as an int.
func MustGetInt(name string) int {
	converted, err := strconv.Atoi(MustGet(name))

	if err != nil {
		panic(fmt.Sprintf(
			"could not convert env var %q value %q into int",
			name,
			MustGet(name),
		))
	}

	return converted
}

// GetInt returns the value of an environment variable as an int.
// It returns an error if the environment variable is not set or cannot be
// converted to an int.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - The value of the environment variable as an int.
//   - An error if the environment variable is not set or cannot be converted to an int.
func GetInt(name string) (int, error) {
	value, err := GetOrErr(name)
	if err != nil {
		return 0, err
	}

	converted, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return converted, nil
}

// MustGetStringSlice returns the value of an environment variable as a string
// slice. It panics if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - The value of the environment variable as a string slice.
func MustGetStringSlice(name string) []string {
	MustLoadEnvVars(nil)

	value := os.Getenv(name)
	if len(value) == 0 {
		panic(fmt.Sprintf("no env var: %q", name))
	}

	return strings.Split(
		value,
		",",
	)
}

// GetStringSlice returns the value of an environment variable as a string
// slice. It returns an error if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - The value of the environment variable as a string slice.
//   - An error if the environment variable is not set.
func GetStringSlice(name string) ([]string, error) {
	value, err := GetOrErr(name)
	if err != nil {
		return []string{}, err
	}

	values := strings.Split(
		value,
		",",
	)

	return values, nil
}

// MustGetBool returns the value of an environment variable as a boolean.
// It panics if the environment variable is not set or cannot be converted to a
// boolean.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - The value of the environment variable as a boolean.
func MustGetBool(name string) bool {
	value := MustGet(name)
	booleanValue, err := strconv.ParseBool(value)

	if err != nil {
		panic(fmt.Sprintf(
			"invalid boolean env var %q value %q",
			name,
			value,
		))
	}

	return booleanValue
}

// GetBool returns the value of an environment variable as a boolean.
// It returns an error if the environment variable is not set or cannot be
// converted to a boolean.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - The value of the environment variable as a boolean.
//   - An error if the environment variable is not set or cannot be converted to a boolean.
func GetBool(name string) (bool, error) {
	value, err := GetOrErr(name)
	if err != nil {
		return false, err
	}
	booleanValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return booleanValue, nil
}

// GetBoolOr returns the value of an environment variable as a boolean.
// It returns the default value if the environment variable is not set or cannot
// be converted to a boolean.
//
// Parameters:
//   - name: The name of the environment variable to get.
//   - defaultValue: The default value to return if the environment variable is not set or cannot be converted to a boolean.
//
// Returns:
//   - The value of the environment variable as a boolean.
func GetBoolOr(name string, defaultValue bool) bool {
	value, err := GetBool(name)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetTyped returns the value of an environment variable as a typed value.
// It returns an error if the environment variable is not set or cannot be
// converted to the expected type.
//
// Parameters:
//   - name: The name of the environment variable to get.
//   - converterFn: A function that converts a string to the expected type.
//
// Returns:
//   - The value of the environment variable as a typed value.
//   - An error if the environment variable is not set or cannot be converted to
//     the expected type.
func GetTyped[T any](
	name string, converterFn func(string) (T, error),
) (T, error) {
	var zero T
	value, err := GetOrErr(name)
	if err != nil {
		return zero, err
	}

	return converterFn(value)
}

// LazyInt returns a function that loads a string from the environment.
// It panics if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - A function that loads a string from the environment.
func LazyString(name string) func() string {
	var once sync.Once
	var value string

	return func() string {
		once.Do(func() {
			value = MustGet(name)
		})

		return value
	}
}

// LazyInt returns a function that loads an int from the environment.
// It panics if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - A function that loads an int from the environment.
func LazyInt(name string) func() int {
	var once sync.Once
	var value int

	return func() int {
		once.Do(func() {
			value = MustGetInt(name)
		})

		return value
	}
}

// LazyInt returns a function that loads a bool from the environment.
// It panics if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - A function that loads a bool from the environment.
func LazyBool(name string) func() bool {
	var once sync.Once
	var value bool

	return func() bool {
		once.Do(func() {
			value = MustGetBool(name)
		})

		return value
	}
}

// LazyInt returns a function that loads a string slice from the environment.
// It panics if the environment variable is not set.
//
// Parameters:
//   - name: The name of the environment variable to get.
//
// Returns:
//   - A function that loads a string slice from the environment.
func LazyStringSlice(name string) func() []string {
	var once sync.Once
	var value []string

	return func() []string {
		once.Do(func() {
			value = MustGetStringSlice(name)
		})

		return value
	}
}

// LazyTyped returns a function that loads a typed value from the environment.
// It panics if the environment conversion fails.
//
// Parameters:
//   - name: The name of the environment variable to get.
//   - converterFn: A function that converts a string to the expected type.
//
// Returns:
//   - A function that loads a typed value from the environment.
func LazyTyped[T any](
	name string, converterFn func(string) (T, error),
) func() T {
	var once sync.Once
	var value T
	var err error

	return func() T {
		once.Do(func() {
			value, err = GetTyped(name, converterFn)
			if err != nil {
				panic(err)
			}
		})

		return value
	}
}
