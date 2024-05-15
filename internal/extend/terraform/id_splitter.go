package terraform

import (
	"errors"
	"regexp"
	"strings"
)

const (
	defaultSortSplitExpression = "\\s*,\\s*"
)

var defaultSortSplitRegex = regexp.MustCompile(defaultSortSplitExpression)

// Split 分割字符串
func Split(id string, element ...*string) error {
	str := defaultSortSplitRegex.Split(strings.TrimSpace(id), -1)
	if len(element) != len(str) {
		return errors.New("分割后的字符串数量与目标数量不一致")
	}
	for i, s := range str {
		*element[i] = s
	}
	return nil
}
