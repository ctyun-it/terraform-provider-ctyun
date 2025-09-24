package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateEipApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateEipApi

	// 构造请求
	var projectID string = "0"
	var bandwidthID string = ""
	var demandBillingType string = "bandwidth"
	var payVoucherPrice string = "1"
	var lineType string = "163"
	var segmentID string = ""
	var exclusiveName string = ""
	request := &CtvpcCreateEipRequest{
		ClientToken:       "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:          "",
		ProjectID:         &projectID,
		CycleType:         "month",
		CycleCount:        1,
		Name:              "acl11",
		Bandwidth:         1,
		BandwidthID:       &bandwidthID,
		DemandBillingType: &demandBillingType,
		PayVoucherPrice:   &payVoucherPrice,
		LineType:          &lineType,
		SegmentID:         &segmentID,
		ExclusiveName:     &exclusiveName,
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
