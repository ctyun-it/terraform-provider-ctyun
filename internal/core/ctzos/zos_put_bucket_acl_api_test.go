package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosPutBucketAclApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosPutBucketAclApi

	// 构造请求
	request := &ZosPutBucketAclRequest{
		Bucket:   "bucket1",
		RegionID: "332232eb-63aa-465e-9028-52e5123866f0",
		ACL:      "private",
		AccessControlPolicy: &ZosPutBucketAclAccessControlPolicyRequest{
			Owner: &ZosPutBucketAclAccessControlPolicyOwnerRequest{
				DisplayName: "xxx",
				ID:          "00000000-0000-0000-0000-000000000000",
			},
			Grants: []*ZosPutBucketAclAccessControlPolicyGrantsRequest{
				{
					Grantee: &ZosPutBucketAclAccessControlPolicyGrantsGranteeRequest{
						EmailAddress: "example@ctyun.cn",
						RawType:      "CanonicalUser",
						DisplayName:  "xxx",
						ID:           "00000000-0000-0000-0000-000000000000",
						URI:          "xxx",
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
