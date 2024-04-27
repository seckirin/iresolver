// Package input pkg/input/input.go
package input

import (
	"github.com/yuukisec/iresolver/pkg/utils"
	"log"
)

func ReadFileLines(filename string) ([]string, error) {
	lines, err := utils.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %s\n", filename)
	}
	return lines, nil
}
