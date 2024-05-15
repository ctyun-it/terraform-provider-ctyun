package utils

// IsLetter 判断是否为英文的字符
func IsLetter(target rune) bool {
	return IsUpper(target) || IsLower(target)
}

// IsUpper 是否为英文大写字母
func IsUpper(target rune) bool {
	return target >= 'A' && target <= 'Z'
}

// IsLower 是否为英文小写字母
func IsLower(target rune) bool {
	return target >= 'a' && target <= 'z'
}

// IsDigit 是否为数字
func IsDigit(target rune) bool {
	return target >= '0' && target <= '9'
}
