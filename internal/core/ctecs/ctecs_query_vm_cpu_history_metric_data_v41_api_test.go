package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsQueryVmCpuHistoryMetricDataV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsQueryVmCpuHistoryMetricDataV41Api

	// 构造请求
	request := &CtecsQueryVmCpuHistoryMetricDataV41Request{
		RegionID:     "100054c0416811e9a6690242ac110002",
		DeviceIDList: []string{},
		Period:       14400,
		StartTime:    "1665305264",
		EndTime:      "1667441264",
		PageNo:       1,
		Page:         1,
		PageSize:     20,
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
