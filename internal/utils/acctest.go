package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func AtLeastOne(value string) error {
	if value == "0" {
		return fmt.Errorf("列表应该至少有一个元素")
	}
	return nil
}

func GetSubdirectories(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var subdirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			subdirs = append(subdirs, entry.Name())
		}
	}
	return subdirs, nil
}

func LoadTestCase(filename string, parameters ...interface{}) string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	fullPath := filepath.Join(pwd, "testdata", filename)
	f, err := os.ReadFile(fullPath)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(string(f), parameters...)
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func GenerateRandomString() string {
	length := 10
	builder := strings.Builder{}
	builder.Grow(length)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		builder.WriteByte(charset[randomIndex])
	}
	return builder.String()
}
func GenerateRandomPort(min int, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成一个 min 到 max 之间的随机数
	randomNum := rand.Intn(max-min) + min // Intn(n) 返回一个范围是 [0, n) 的随机数
	return randomNum
}
