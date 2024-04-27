// pkg/exporting/export.go
package exporting

import (
	"fmt"
	"os"
)

func ExportToFile(filename string, data []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	for _, line := range data {
		_, err := fmt.Fprintln(file, line)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}
