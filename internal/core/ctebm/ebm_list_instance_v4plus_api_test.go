package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbmListInstanceV4plusApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbmListInstanceV4plusApi

	// 构造请求
	var resourceID string = "f13e86bd53b4426594796987faaea01c"
	var ip string = "100.124.28.100"
	var instanceName string = "pm-10901"
	var vpcID string = "a64ddb73-3021-5a8d-9abd-b0a12e429690"
	var subnetID string = "c45512bf-7919-55a9-a106-5f4aa9194c7c"
	var deviceType string = "physical.t3.large"
	var deviceUUIDList string = "d-gutvjbzdfquu9ufiaheuerrpjhxz,d-gutvjbzdfquu9ufiaheuerrpjhxa"
	var queryContent string = "192.168.0.210"
	var instanceUUIDList string = "ss-rtkjng87ovlyda013spdb0shvfaj,ss-rtkjng87ovlyda013spdb0shvfat"
	var instanceUUID string = "ss-wcwkjt374iq5ce8orekenziqdfs5"
	var status string = "RUNNING"
	var sort string = "expire_time"
	var asc bool = true
	var vipID string = "havip-b3lbmyhj27"
	var volumeUUID string = "614aadc6-81d5-43de-8bbb-0397371d2ed8"
	request := &EbmListInstanceV4plusRequest{
		RegionID:         "73f321ea-62ff-11ec-a8bc-005056898fe0",
		AzName:           "",
		ResourceID:       &resourceID,
		Ip:               &ip,
		InstanceName:     &instanceName,
		VpcID:            &vpcID,
		SubnetID:         &subnetID,
		DeviceType:       &deviceType,
		DeviceUUIDList:   &deviceUUIDList,
		QueryContent:     &queryContent,
		InstanceUUIDList: &instanceUUIDList,
		InstanceUUID:     &instanceUUID,
		Status:           &status,
		Sort:             &sort,
		Asc:              &asc,
		VipID:            &vipID,
		VolumeUUID:       &volumeUUID,
		PageNo:           1,
		PageSize:         10,
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
