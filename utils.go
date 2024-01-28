package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func IsRegexMatch(matchStr string, patternArr []string) bool {
	for i := 0; i < len(patternArr); i++ {
		match, _ := regexp.MatchString(patternArr[i], matchStr)
		if match {
			return true
		}
	}
	return false
}

func IsSubDomain(domain string, targets []string) bool {
	for i := 0; i < len(targets); i++ {
		domain = "." + strings.TrimLeft(domain, ".")
		target := "." + strings.TrimLeft(targets[i], ".")
		match := strings.HasSuffix(domain, target)
		if match {
			return true
		}
	}
	return false
}

func ReadUrlsFromFile(filename string) ([]string, error) {
	// 定义一个切片来存储 URLs
	var urls []string

	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 使用 bufio.Scanner 逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 将每一行的内容添加到切片中
		urls = append(urls, scanner.Text())
	}

	// 检查是否有读取文件时出现的错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件时出现错误: %v", err)
	}

	return urls, nil
}

func ShowLogo() {
	fmt.Println("\033[32m") // 设置文本颜色为绿色

	logo := `
	__  __     ______     ______   __     _____     ______     ______    
	/\_\_\_\   /\  ___\   /\  == \ /\ \   /\  __-.  /\  ___\   /\  == \   
	\/_/\_\/_  \ \___  \  \ \  _-/ \ \ \  \ \ \/\ \ \ \  __\   \ \  __<   
	  /\_\/\_\  \/\_____\  \ \_\    \ \_\  \ \____-  \ \_____\  \ \_\ \_\ 
	  \/_/\/_/   \/_____/   \/_/     \/_/   \/____/   \/_____/   \/_/ /_/ 																		
`
	fmt.Print(logo)
	fmt.Println("\033[0m") // 恢复文本颜色
}
