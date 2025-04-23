// Package file 文件操作辅助函数
package file

import (
	"os"
)

// Put 将数据存入文件
func Put(data []byte, to string) error {
	err := os.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Exists 判断文件是否存在
func Exists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}

// EnsureDirExists 确保目录存在，不存在则创建
func EnsureDirExists(dirPath string) error {
	if Exists(dirPath) {
		return nil
	}
	return os.MkdirAll(dirPath, 0755)
}
