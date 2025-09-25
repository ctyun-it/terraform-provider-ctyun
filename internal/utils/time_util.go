package utils

import (
	"time"
)

// 从RFC3339转换到本地时间格式
func FromRFC3339ToLocal(timeStr string) string {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t.In(time.Local).Format(time.DateTime)
}
