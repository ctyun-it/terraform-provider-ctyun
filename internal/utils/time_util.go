package utils

import (
	"time"
)

const (
	Layout1 = "2006-01-02 15:04:05 -0700 MST"
	Layout2 = "2006-01-02T15:04:05.000-07:00"
	Layout3 = "2006-01-02T15:04:05-07:00"
	Layout4 = "2006-01-02 15:04:05.000000Z"
)

// 从RFC3339转换到本地时间格式
func FromRFC3339ToLocal(timeStr string) string {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t.In(time.Local).Format(time.DateTime)
}

// FromUnixToUTC 将Unix时间戳（10位秒级/13位毫秒级）转换为UTC时间字符串
// 输入的时间值：1717574263（10位）、1717574263000（13位）
// 返回值: 转换后的UTC时间字符串，如"2025-09-09T03:42:52Z"，无效输入返回空
func FromUnixToUTC(timestamp int64) string {
	// 调整范围校验：包含13位合理最大值（约2286年）
	if timestamp <= 0 || timestamp > 3000000000000 {
		return ""
	}

	var t time.Time
	if timestamp > 9999999999 {
		sec := timestamp / 1000
		nsec := (timestamp % 1000) * 1e6
		t = time.Unix(sec, nsec).UTC()
	} else {
		t = time.Unix(timestamp, 0).UTC()
	}

	return t.Format(time.RFC3339)
}

// ConvertToUTCZ 将时间字符串转换为UTC时间，并格式化为2006-01-02T15:04:05Z格式
// 返回值: 转换后的UTC时间字符串，如"2025-09-09T03:42:52Z"
func ConvertToUTCZ(layout, input string) string {
	t, err := time.Parse(layout, input)
	if err != nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// FromBJTimeToUTC 将北京时间（YYYY-MM-DD HH:MM:SS）转换为UTC时间的RFC3399格式（带Z）
// 入参：input - 北京时间字符串（格式：2025-12-02 13:37:32）
// 出参：UTC时间的RFC3399字符串，解析失败返回空字符串
func FromBJTimeToUTCZ(input string) string {
	layout := time.DateOnly + " " + time.TimeOnly
	shanghaiLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return ""
	}
	t, err := time.ParseInLocation(layout, input, shanghaiLoc)
	if err != nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
