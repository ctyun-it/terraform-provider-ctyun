package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
	"strconv"
)

func ParseAcl(grants []*s3.Grant) string {
	// 1. 检查grants是否为空，为空默认私有权限
	if len(grants) == 0 {
		return "private"
	}

	// 2. 检查grants是否包含公共读写权限（AllUsers + FULL_CONTROL）
	for _, g := range grants {
		if g != nil && g.Grantee != nil &&
			SecString(g.Grantee.URI) == "http://acs.amazonaws.com/groups/global/AllUsers" &&
			SecString(g.Permission) == "FULL_CONTROL" {
			return "public-read-write"
		}
	}

	// 3. 检查grants是否包含公共读权限（AllUsers + READ）
	hasAllUsersRead := false
	for _, g := range grants {
		if g != nil && g.Grantee != nil &&
			SecString(g.Grantee.URI) == "http://acs.amazonaws.com/groups/global/AllUsers" &&
			SecString(g.Permission) == "READ" {
			hasAllUsersRead = true
			break
		}
	}

	if hasAllUsersRead {
		// 4. 有公共读时再检查是否包含写权限，即公共读写
		for _, g := range grants {
			if g != nil && g.Grantee != nil &&
				SecString(g.Grantee.URI) == "http://acs.amazonaws.com/groups/global/AllUsers" &&
				(SecString(g.Permission) == "FULL_CONTROL" || SecString(g.Permission) == "WRITE") {
				return "public-read-write"
			}
		}
		return "public-read"
	}
	return "private"
}

func GetFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("获取镜像源文件失败")
	}
	defer resp.Body.Close() // 必须关闭

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("获取镜像源文件失败")
	}

	// 3. 从响应头获取 Content-Length
	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("获取镜像源文件失败")
	}

	// 4. 转为数字（字节）
	sizeBytes, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("获取镜像源文件失败")
	}
	return sizeBytes / 1024 / 1024 / 1024, nil
}
