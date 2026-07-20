package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	// 1. 定义命令行参数
	path := flag.String("path", "", "需要重命名的文件夹路径")
	oldStr := flag.String("old", "", "要被替换的旧字符串")
	newStr := flag.String("new", "", "替换后的新字符串")
	flag.Parse()

	// 2. 参数校验
	if *path == "" || *oldStr == "" || *newStr == "" {
		fmt.Println("用法: go run day6_renamer_advanced.go -path=文件夹 -old=旧字符串 -new=新字符串")
		return
	}
	err := filepath.WalkDir(*path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		oldName := d.Name()
		if strings.Contains(oldName, *oldStr) {
			newName := strings.ReplaceAll(oldName, *oldStr, *newStr)

			//拼接
			oldPath := path
			newPath := filepath.Join(filepath.Dir(path), newName)

			err := os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Printf("重命名失败[%s]: %v \n", oldName, err)
			} else {
				fmt.Printf("重命名成功 %v -> %v ", oldName, newName)
			}
		}
		return nil

	})
	if err != nil {
		fmt.Printf("遍历文件夹出错: %v\n", err)
	} else {
		fmt.Println("🎉 递归批量重命名处理完毕！")
	}

}
