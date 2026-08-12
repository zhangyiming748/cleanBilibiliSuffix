package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"cleanBilibiliSuffix/core"
)

// 以下变量由 goreleaser 在编译时通过 ldflags 注入，与触发工作流的 tag 联动
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	var keywordPath string
	var dir string

	rootCmd := &cobra.Command{
		Use:     "CBS",
		Short:   "CBS - 清理哔哩哔哩下载视频文件名中的多余前缀",
		Version: fmt.Sprintf("%s (commit %s, built %s)", version, gitCommit, buildTime),
		RunE: func(cmd *cobra.Command, args []string) error {
			// 进入 RunE 说明参数已解析正确，后续运行时错误不再打印用法
			cmd.SilenceUsage = true
			if err := core.LoadKeywords(keywordPath); err != nil {
				return err
			}
			// 边查找边处理：每找到一个视频文件就立即重命名，不先攒完整列表
			return core.WalkVideoFiles(dir, func(path string) {
				newPath, err := core.CleanFile(path)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
				if newPath != path {
					cmd.Printf("重命名成功: %s -> %s\n", filepath.Base(path), filepath.Base(newPath))
				} else {
					cmd.Println("无需修改:", filepath.Base(path))
				}
			})
		},
	}

	rootCmd.Flags().StringVarP(&keywordPath, "keyword", "k", "", "关键词文本文档的路径，每行一个关键词")
	rootCmd.Flags().StringVarP(&dir, "dir", "d", "", "查找视频文件的根目录")
	rootCmd.MarkFlagRequired("keyword")
	rootCmd.MarkFlagRequired("dir")

	// Execute 内部已打印错误信息，这里只需退出码
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
