package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosCreateZMSMigrationTaskApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosCreateZMSMigrationTaskApi

	// 构造请求
	request := &ZosCreateZMSMigrationTaskRequest{
		RegionID:         "332232xxxxx5123866f0",
		MigrationName:    "cxx-hkpblz-bbt",
		StorageType:      "STANDARD",
		AclConf:          "match-resource",
		ConflictMode:     "IGNORE",
		MigrateStartTime: "1970-01-04 08:00:00",
		MigrateEndTime:   "2025-01-01 08:00:00",
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
