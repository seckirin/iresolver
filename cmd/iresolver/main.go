// cmd/iresolver/main.go
package main

import (
	"fmt"
	"github.com/yuukisec/iresolver/pkg/options"
	"github.com/yuukisec/iresolver/pkg/output"
	"github.com/yuukisec/iresolver/pkg/resolve"
)

func main() {
	opts := options.ParseOptions()

	// 读取文件内容
	tLines, err := options.GetTargetServers(opts.Target)
	if err != nil {
		output.PrintError(fmt.Errorf("Error getting target servers: %w", err))
		return
	}

	// 打印文件行数
	if !opts.Silent {
		output.PrintInfo(fmt.Sprintf("File/Server %s contains %d lines", opts.Target, len(tLines)))
	}

	// 解析并比较 DNS 结果
	resolve.ResolveAndCompare(opts)
}
