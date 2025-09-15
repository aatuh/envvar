package envvar

import (
	"fmt"
	"os"
)

var envFilePaths = []string{".env", "/env/.env"}

// MustLoadEnvVars loads environment variables from the environment file.
// It panics if the environment variables cannot be loaded.
//
// Parameters:
//   - customEnvFilePaths: The custom environment file paths to load.
func MustLoadEnvVars(customEnvFilePaths []string) {
	if err := LoadEnvVars(customEnvFilePaths); err != nil {
		panic(err)
	}
}

// LoadEnvVars loads environment variables from the environment file.
// It returns an error if the environment variables cannot be loaded.
//
// Parameters:
//   - customEnvFilePaths: The custom environment file paths to load.
//
// Returns:
//   - An error if the environment variables cannot be loaded.
func LoadEnvVars(customEnvFilePaths []string) error {
	useEnvFilePaths := envFilePaths
	if len(customEnvFilePaths) > 0 {
		useEnvFilePaths = customEnvFilePaths
	}
	var retErr error
	once.Do(func() {
		fileFound := false
		for _, path := range useEnvFilePaths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}
			fileFound = true
			retErr = LoadAndSetEnvVarsFromPath(path)
			break
		}

		// If no env file found, fallback to system environment variables.
		if !fileFound {
			retErr = nil
		}
	})

	return retErr
}

// LoadAndSetEnvVarsFromPath loads environment variables from a file and sets
// them in the environment.
//
// Parameters:
//   - path: The path to the environment file to load.
//
// Returns:
//   - An error if the environment variables cannot be loaded.
func LoadAndSetEnvVarsFromPath(path string) error {
	envVars, err := ReadEnvVarFile(path)
	if err != nil {
		return fmt.Errorf(
			"error reading env file: %v",
			err,
		)
	}

	if err := SetEnvVars(envVars); err != nil {
		return fmt.Errorf(
			"error setting env vars: %v",
			err,
		)
	}

	return nil
}
