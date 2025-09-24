package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcEipQueryCreatePriceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcEipQueryCreatePriceApi

	// 构造请求
	var projectID string = "0"
	var bandwidthID string = ""
	var demandBillingType string = "bandwidth"
	request := &CtvpcEipQueryCreatePriceRequest{
		ClientToken:       "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:          "",
		ProjectID:         &projectID,
		CycleType:         "month",
		CycleCount:        1,
		Name:              "",
		Bandwidth:         1,
		BandwidthID:       &bandwidthID,
		DemandBillingType: &demandBillingType,
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
