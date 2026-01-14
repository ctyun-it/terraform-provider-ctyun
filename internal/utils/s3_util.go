package utils

import "github.com/aws/aws-sdk-go/service/s3"

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
