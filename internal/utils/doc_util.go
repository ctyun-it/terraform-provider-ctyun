package utils

import (
	"fmt"
)

func FormatDesc(purpose, subcategory, url string) string {
	return fmt.Sprintf(`-> %s
%s
%s`, purpose, subcategory, url)
}
