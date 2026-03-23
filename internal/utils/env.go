package utils

import (
	"os"
	"strings"
)

const (
	EnvPrefix = "CTYUN_"
)

// AcquireEnvKey 获取环境变量中的key
// 例如：入参为ak，返回CTYUN_AK
func AcquireEnvKey(name string) string {
	return EnvPrefix + strings.ToUpper(name)
}

// AcquireEnvParam 从环境变量中获取参数
// 例如：入参为ak，返回CTYUN_AK环境变量中的值
func AcquireEnvParam(name string) string {
	key := AcquireEnvKey(name)
	value := os.Getenv(key)
	if key == "CTYUN_PROJECT_ID" && strings.TrimSpace(value) == "" {
		return "0"
	}
	return value
}
