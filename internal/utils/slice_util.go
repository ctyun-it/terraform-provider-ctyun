package utils

// 泛型函数：判断 val 是否等于切片 s 中的任意一个变量（支持所有可比较类型）
func Contain[T comparable](s []T, val T) bool {
	for _, item := range s {
		if item == val {
			return true
		}
	}
	return false
}
