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
		Use:          "CBS",
		Short:        "CBS - 清理哔哩哔哩下载视频文件名中的多余前缀",
		Version:      fmt.Sprintf("%s (commit %s, built %s)", version, gitCommit, buildTime),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := core.LoadKeywords(keywordPath); err != nil {
				return err
			}
			// 边查找边处理：每找到一个视频文件就立即重命名，不先攒完整列表
			return core.WalkVideoFiles(dir, func(path string) {
				newPath, err := core.CleanFile(path)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
				if newPath != path {
					fmt.Printf("重命名成功: %s -> %s\n", filepath.Base(path), filepath.Base(newPath))
				} else {
					fmt.Println("无需修改:", filepath.Base(path))
				}
			})
		},
	}

	rootCmd.Flags().StringVarP(&keywordPath, "keyword", "k", "", "关键词文本文档的路径，每行一个关键词")
	rootCmd.Flags().StringVarP(&dir, "dir", "d", "", "查找视频文件的根目录")
	rootCmd.MarkFlagRequired("keyword")
	rootCmd.MarkFlagRequired("dir")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
