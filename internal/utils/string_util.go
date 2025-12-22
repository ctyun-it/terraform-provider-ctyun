package utils

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
)

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

// SecString *string转string
func SecString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

// SecStringValue *string转types.String
func SecStringValue(str *string) types.String {
	if str == nil {
		return types.StringValue("")
	}
	return types.StringValue(*str)
}

// SecLowerStringValue *string转types.String全小写
func SecLowerStringValue(str *string) types.String {
	if str == nil {
		return types.StringValue("")
	}
	return types.StringValue(strings.ToLower(*str))
}

// SecUpperStringValue *string转types.String全大写
func SecUpperStringValue(str *string) types.String {
	if str == nil {
		return types.StringValue("")
	}
	return types.StringValue(strings.ToLower(*str))
}

// StrPointerArrayToStrArray []*string转[]string
func StrPointerArrayToStrArray(array []*string) []string {
	ret := []string{}
	for _, str := range array {
		if str != nil {
			ret = append(ret, *str)
		} else {
			ret = append(ret, "")
		}
	}
	return ret
}

// StrArrayToStrPointerArray []string转[]*string
func StrArrayToStrPointerArray(array []string) []*string {
	ret := []*string{}
	for _, str := range array {
		s := str
		ret = append(ret, &s)
	}
	return ret
}

// DifferenceStrArray 获取两个字符串数组的差集并去重
func DifferenceStrArray(a, b []string) (diffA []string, diffB []string) {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	diffAMap := make(map[string]struct{})
	// 找出 a 中不在 b 里的元素
	for _, x := range a {
		if _, found := mb[x]; !found {
			if _, exists := diffAMap[x]; !exists {
				diffA = append(diffA, x)
				diffAMap[x] = struct{}{}
			}
		}
	}

	ma := make(map[string]struct{}, len(a))
	for _, x := range a {
		ma[x] = struct{}{}
	}
	diffBMap := make(map[string]struct{})
	// 找出 b 中不在 a 里的元素
	for _, x := range b {
		if _, found := ma[x]; !found {
			if _, exists := diffBMap[x]; !exists {
				diffB = append(diffB, x)
				diffBMap[x] = struct{}{}
			}
		}
	}
	return
}

// AreStringSlicesEqual 对比字符串数组是否相同
func AreStringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	count := make(map[string]int, len(a))
	for _, s := range a {
		count[s]++
	}
	for _, s := range b {
		if count[s]--; count[s] < 0 {
			return false
		}
	}
	return true
}

// StringToInt32Must 字符串转int32
func StringToInt32Must(s string) int32 {
	num, _ := strconv.ParseInt(s, 10, 64)
	return int32(num)
}

func JsonString(obj interface{}) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

func StringArrayToValueArray(array []string) []attr.Value {
	var ret []attr.Value
	for _, str := range array {
		ret = append(ret, types.StringValue(str))
	}
	return ret
}

func ParseYamlValue(input string, key string) (value interface{}, err error) {
	var parsed map[string]interface{}
	decoder := yaml.NewDecoder(strings.NewReader(input))
	err = decoder.Decode(&parsed)
	if err != nil {
		return
	}
	value, err = GetValue(key, parsed)
	if err != nil {
		return
	}

	return
}

func GetValue(key string, m map[string]interface{}) (interface{}, error) {
	if strings.Contains(key, ".") {
		keys := strings.SplitN(key, ".", 2)
		if _, ok := m[keys[0]]; !ok {
			return nil, fmt.Errorf("必须存在key: %s", keys[0])
		}
		next, ok := m[keys[0]].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s 没有下一层级", keys[0])
		}
		return GetValue(keys[1], next)
	} else {
		if _, ok := m[key]; !ok {
			return nil, fmt.Errorf("必须存在key: %s", key)
		}
		return m[key], nil
	}
}
