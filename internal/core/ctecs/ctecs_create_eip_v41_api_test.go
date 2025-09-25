package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateEipV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateEipV41Api

	// 构造请求
	request := &CtecsCreateEipV41Request{
		ClientToken:       "create-eip-test",
		RegionID:          "bb9fdb42056f11eda1610242ac110002",
		ProjectID:         "0",
		CycleType:         "month",
		CycleCount:        1,
		Name:              "eip-name",
		Bandwidth:         5,
		BandwidthID:       "bandwidth-7hzv449r2j",
		DemandBillingType: "bandwidth",
		PayVoucherPrice:   "1",
		LineType:          "chinaunicom",
		SegmentID:         "seg-1kxkdkbt9e",
		ExclusiveName:     "eip-block-ipv4",
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
