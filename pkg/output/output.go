// pkg/output/output.go
package output

import (
	"fmt"
	"time"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
	SuccessColor = "\033[1;32m%s\033[0m"
)

func PrintInfo(info string) {
	timestamp := time.Now().Format("2006-01-02 15:04")
	fmt.Printf("[INFO] [%s] "+InfoColor+"\n", timestamp, info)
}

func PrintError(err error) {
	timestamp := time.Now().Format("2006-01-02 15:04")
	fmt.Printf("[ERROR] [%s] "+ErrorColor+"\n", timestamp, err)
}

func PrintSuccess(info string) {
	timestamp := time.Now().Format("2006-01-02 15:04")
	fmt.Printf("[SUCCESS] [%s] "+SuccessColor+"\n", timestamp, info)
}

func PrintSilent(server string) {
	fmt.Println(server)
}
