package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosPutObjectAclApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosPutObjectAclApi

	// 构造请求
	request := &ZosPutObjectAclRequest{
		Bucket:    "bucket1",
		Key:       "obj1",
		VersionID: "USzD.sN0vODAsQ84ncdT20oRiY2lFCD",
		RegionID:  "332232eb-63aa-465e-9028-52e5123866f0",
		ACL:       "private",
		AccessControlPolicy: &ZosPutObjectAclAccessControlPolicyRequest{
			Owner: &ZosPutObjectAclAccessControlPolicyOwnerRequest{
				DisplayName: "dname1",
				ID:          "user1",
			},
			Grants: []*ZosPutObjectAclAccessControlPolicyGrantsRequest{
				{
					Grantee: &ZosPutObjectAclAccessControlPolicyGrantsGranteeRequest{
						EmailAddress: "example@ctyun.cn",
						RawType:      "CanonicalUser",
						DisplayName:  "dname2",
						ID:           "user2",
						URI:          "http://acs.amazonaws.com/groups/global/AuthenticatedUsers",
					},
					Permission: "FULL_CONTROL",
				},
			},
		},
	}

	// 发起调用
	response, err := api.Do(context.Background(), *credential, request)
	if err != nil {
		t.Log("request error:", err)
		t.Fail()
		return
	}
	t.Logf("%+v\n", *response)
}
