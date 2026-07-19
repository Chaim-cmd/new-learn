package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	//定义命令行参数
	//在终端输入 -path="./test" 时，path 变量就会等于"./test"
	path := flag.String("path", "", "需要重命名的文件夹路径")
	oldStr := flag.String("old", "", "要被替换的旧字符串")
	newStr := flag.String("new", "", "替换后新的字符串")

	//解析参数
	flag.Parse()

	//参数校验
	if *path == "" || *oldStr == "" || *newStr == "" {
		fmt.Println("具体使用：go run main.go -path=文件夹 -old=旧字符串 -new=新字符串")
	}

	//指定某个文件夹下的所有文件
	entries, err := os.ReadDir(*path)
	if err != nil {
		fmt.Printf("读取文件失败：%v\n", err)
	}

	//遍历文件
	for _, v := range entries {
		//跳过文件夹
		if v.IsDir() {
			continue
		}
		oldName := v.Name()
		//判断文件中是否有重复的字符串
		if strings.Contains(oldName, *oldStr) {
			//生成新的文件名
			newName := strings.ReplaceAll(oldName, *oldStr, *newStr)

			//拼接完整的路径
			oldDir := filepath.Join(*path, oldName)
			newDir := filepath.Join(*path, newName)

			err := os.Rename(oldDir, newDir)
			if err != nil {
				fmt.Printf("重命名失败%v\n", err)
			} else {
				fmt.Printf("重命名成功，%v -> %v \n", oldName, newName)
			}

		}
	}

}
