package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// keywords 是包级关键词切片，程序启动时由 LoadKeywords 读取一次，
// 之后每次重命名直接查这个切片，无需重复读文件。
var keywords []string

// LoadKeywords 从文本文档中读取前缀关键词，每行一个，忽略空行，
// 结果存入包级 keywords 切片。只需在程序启动时调用一次。
func LoadKeywords(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打开关键词文件失败: %w", err)
	}
	defer f.Close()

	keywords = keywords[:0]
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			keywords = append(keywords, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取关键词文件失败: %w", err)
	}
	if len(keywords) == 0 {
		return fmt.Errorf("关键词文件 %s 中没有任何关键词", path)
	}
	return nil
}

// CleanFile 接收一个文件的绝对路径，拆分出路径名和文件名，
// 用包级 keywords 切片删除文件名开头的关键词前缀以及随之出现的开头空格，然后执行重命名。
// 返回处理后的新绝对路径；若文件名无需修改，原样返回。
func CleanFile(absPath string) (string, error) {
	if len(keywords) == 0 {
		return "", fmt.Errorf("尚未加载关键词，请先调用 LoadKeywords")
	}
	if !filepath.IsAbs(absPath) {
		return "", fmt.Errorf("必须提供绝对路径: %s", absPath)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("访问文件失败: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("目标是一个目录: %s", absPath)
	}

	dir := filepath.Dir(absPath)
	name := filepath.Base(absPath)

	newName := cleanName(name)
	if newName == name {
		return absPath, nil
	}
	if newName == "" || newName == filepath.Ext(name) {
		return "", fmt.Errorf("清理后文件名为空，跳过重命名: %s", name)
	}

	newPath := filepath.Join(dir, newName)
	if _, err := os.Stat(newPath); err == nil {
		return "", fmt.Errorf("目标文件已存在，无法重命名: %s", newPath)
	}
	if err := os.Rename(absPath, newPath); err != nil {
		return "", fmt.Errorf("重命名失败: %w", err)
	}
	return newPath, nil
}

// cleanName 反复剥离文件名开头的关键词前缀，并在每次剥离后删除开头的空格，
// 以处理「4K修复 1080P高清修复 视频名」这类多重前缀。
func cleanName(name string) string {
	for {
		changed := false
		for _, kw := range keywords {
			if strings.HasPrefix(name, kw) {
				name = strings.TrimPrefix(name, kw)
				changed = true
				break
			}
		}
		// 删除开头空格（含全角空格）
		trimmed := strings.TrimLeft(name, " \u3000")
		if trimmed != name {
			name = trimmed
			changed = true
		}
		if !changed {
			return name
		}
	}
}
