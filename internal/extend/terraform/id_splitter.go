package terraform

import (
	"errors"
	"regexp"
	"strings"
)

const (
	defaultSortSplitExpression = "\\s*,\\s*"
	SortSplitExpression        = "\\s*:\\s*"
)

var defaultSortSplitRegex = regexp.MustCompile(defaultSortSplitExpression)
var SortSplitRegex = regexp.MustCompile(SortSplitExpression)

// Split 分割字符串
func Split(id string, element ...*string) error {
	str := defaultSortSplitRegex.Split(strings.TrimSpace(id), -1)
	if len(element) != len(str) {
		return errors.New("分割后的字符串数量与目标数量不一致，请按照导入命令重新导入")
	}
	for i, s := range str {
		*element[i] = s
	}
	return nil
}

func SplitComma(id string, element ...*string) error {
	str := SortSplitRegex.Split(strings.TrimSpace(id), -1)
	if len(element) != len(str) {
		return errors.New("分割后的字符串数量与目标数量不一致，请按照导入命令重新导入")
	}
	for i, s := range str {
		*element[i] = s
	}
	return nil
}
