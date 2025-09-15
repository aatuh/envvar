package envvar

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadEnvVarFileAndSetEnvVars reads an env var file, parses the contents, and
// sets the environment variables.
//
// Parameters:
//   - filePath: The path to the environment file to read.
//
// Returns:
//   - A map of environment variables.
//   - An error if the environment file cannot be read.
func ReadEnvVarFile(filePath string) (map[string]string, error) {
	file, err := openFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envVars, err := parseEnvVars(file)
	if err != nil {
		return nil, err
	}

	return envVars, nil
}

// SetEnvVars sets environment variables from a key-value map.
//
// Parameters:
//   - envVars: The environment variables to set.
//
// Returns:
//   - An error if the environment variables cannot be set.
//
// Errors:
//   - If the environment variables cannot be set.
func SetEnvVars(envVars map[string]string) error {
	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf(
				"error setting environment variable %s: %v",
				key,
				err,
			)
		}
	}
	return nil
}

// openFile opens the environment file.
func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf(
			"error while opening env var file %q: %v",
			filePath,
			err,
		)
	}
	return file, nil
}

// parseEnvVars parses the environment variables from the file.
func parseEnvVars(file *os.File) (map[string]string, error) {
	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning env var file: %v", err)
	}

	return envVars, nil
}
