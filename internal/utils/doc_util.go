package utils

import "fmt"

func FormatDesc(subcategory, url string) string {
	return fmt.Sprintf(`%s
-> 详细说明请见文档：%s`, subcategory, url)
}
