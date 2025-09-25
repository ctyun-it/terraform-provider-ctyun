package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type ChangeLog struct {
	Entries map[string][]string
}

const changelogDir = ".changelog"

func main() {
	// 从命令行参数获取版本号
	if len(os.Args) < 2 {
		panic("请提供版本号作为命令行参数，例如: go run main.go v1.0.0")
	}
	version := os.Args[1]

	// 获取今天的日期并格式化为 "January 11, 2022" 形式
	date := time.Now().Format("January 2, 2006")

	cl := &ChangeLog{
		Entries: make(map[string][]string),
	}

	// 读取所有PR文件
	files, err := os.ReadDir(changelogDir)
	if err != nil {
		panic(fmt.Sprintf("Error reading changelog directory: %v", err))
	}

	// 解析每个文件内容
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".txt" {
			processFile(cl, filepath.Join(changelogDir, f.Name()))
		}
	}

	// 生成最终CHANGELOG.md，传入版本号和日期
	generateMarkdown(cl, version, date)
}

func processFile(cl *ChangeLog, path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Error reading file %s: %v", path, err))
	}

	lines := strings.Split(string(content), "\n")
	var currentType string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "release-note:") {
			currentType = strings.TrimPrefix(line, "release-note:")
			continue
		}

		if line != "" && currentType != "" {
			// 去重处理
			if !contains(cl.Entries[currentType], line) {
				cl.Entries[currentType] = append(cl.Entries[currentType], line)
			}
		}
	}
}

func generateMarkdown(cl *ChangeLog, version string, date string) {
	// 定义输出顺序和标题
	categories := []struct {
		Key   string
		Title string
	}{
		{"new-resource", "New Resources"},
		{"new-data-source", "New Data Sources"},
		{"enhancement", "Enhancements"},
		{"bug", "Bug Fixes"},
		{"deprecation", "Deprecations"},
	}

	// 创建模板，增加版本和日期信息
	tmpl := `# Changelog

## {{.Version}} - {{.Date}}
{{range .Categories}}
### {{.Title}}
{{range $index, $entry := (index $.Entries .Key)}}
* {{$entry}}{{end}}
{{end}}`

	// 准备模板数据，包含版本和日期
	data := struct {
		Version    string
		Date       string
		Entries    map[string][]string
		Categories []struct {
			Key   string
			Title string
		}
	}{
		Version:    version,
		Date:       date,
		Entries:    cl.Entries,
		Categories: categories,
	}

	// 执行模板渲染
	t := template.Must(template.New("changelog").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		panic(fmt.Sprintf("Error executing template: %v", err))
	}

	// 写入文件
	if err := os.WriteFile("CHANGELOG.md", buf.Bytes(), 0644); err != nil {
		panic(fmt.Sprintf("Error writing CHANGELOG.md: %v", err))
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
