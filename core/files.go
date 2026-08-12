package core

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
)

// FindVideoFiles 递归查找指定目录下任意层级的视频文件，返回它们的绝对路径列表。
func FindVideoFiles(root string) ([]string, error) {
	var files []string
	err := WalkVideoFiles(root, func(path string) {
		files = append(files, path)
	})
	return files, err
}

// WalkVideoFiles 递归遍历指定目录，每找到一个视频文件就立即调用 fn 处理，
// 不先生成完整文件列表，适合边查找边重命名的场景。
func WalkVideoFiles(root string, fn func(path string)) error {
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("访问目录失败: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("不是一个目录: %s", root)
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() && isVideo(path) {
			abs, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			fn(abs)
		}
		return nil
	})
}

// isVideo 通过文件头魔数检测是否为视频；对于 filetype 无法识别的
// rmvb/rm/vob/flv/ts/m2ts 格式，退而用扩展名判断。
func isVideo(fp string) bool {
	file, err := os.Open(fp)
	if err != nil {
		return false
	}
	defer file.Close()
	head := make([]byte, 261)
	if _, err := file.Read(head); err != nil && err != io.EOF {
		return false
	}
	if filetype.IsVideo(head) {
		return true
	}
	switch strings.ToLower(filepath.Ext(fp)) {
	case ".rmvb", ".rm", ".vob", ".flv", ".ts", ".m2ts":
		return true
	}
	return false
}
