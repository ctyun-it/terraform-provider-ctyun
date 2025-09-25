package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsInstanceAttachSfsV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsInstanceAttachSfsV41Api

	// 构造请求
	var autoMount bool = false
	request := &CtecsInstanceAttachSfsV41Request{
		RegionID:   "bb9fdb42056f11eda1610242ac110002",
		InstanceID: "b67b7f1f-095b-1249-b379-8dd5cc542a05",
		SysInfoList: []*CtecsInstanceAttachSfsV41SysInfoListRequest{
			{
				FileSysID:    "56644622-41d1-5c35-8794-df55bffaff88",
				FileSysRoute: "55.243.4.20:/mnt/sfs_cap/e3aacef1e0be40559b38a5c2158aa62e_mqfb3ej839cdz7ub",
				MountPoint:   "/mnt/docs",
				Option:       "vers=3,async,nolock,noatime,nodiratime,wsize=1048576,rsize=1048576,timeo=600",
				AutoMount:    &autoMount,
				Protocol:     "NFSv3",
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
