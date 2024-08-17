// Package env provides functionality for loading environment variables from a file.

package env

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv loads environment variables from a specified file.
// The file format should be "key=value" on each line.
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			err := os.Setenv(key, value)
			if err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
