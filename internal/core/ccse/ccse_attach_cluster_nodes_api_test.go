package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseAttachClusterNodesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseAttachClusterNodesApi

	// 构造请求
	var isSyncClusterResourceLabels bool = false
	var cpuCFSQuota bool = true
	request := &CcseAttachClusterNodesRequest{
		ClusterId: "47281b02f87757478f20b1827c97cadf",
		RegionId:  "bb9fdb42056f11eda1610242ac110002",
		Instances: []*CcseAttachClusterNodesInstancesRequest{
			{
				InstanceId: "",
				AzName:     "",
			},
		},
		VmType:      "ecs",
		Runtime:     "",
		ImageType:   1,
		ImageUuid:   "",
		LoginType:   "",
		Password:    "",
		KeyName:     "KeyPair-a589",
		KeyPairId:   "ba425a97-9ad9-2d45-e21a-770d2ebeb477",
		Labels:      &CcseAttachClusterNodesLabelsRequest{},
		Annotations: &CcseAttachClusterNodesAnnotationsRequest{},
		Taints: []*CcseAttachClusterNodesTaintsRequest{
			{
				Key:    "",
				Value:  "",
				Effect: "",
			},
		},
		VisibilityPostHostScript: "",
		VisibilityHostScript:     "",
		KubeletArgs: &CcseAttachClusterNodesKubeletArgsRequest{
			KubeAPIQPS:           50,
			KubeAPIBurst:         100,
			MaxPods:              110,
			RegistryPullQPS:      5,
			RegistryBurst:        10,
			PodPidsLimit:         -1,
			EventRecordQPS:       50,
			EventBurst:           100,
			TopologyManagerScope: "container",
			CpuCFSQuota:          &cpuCFSQuota,
		},
		IsSyncClusterResourceLabels: &isSyncClusterResourceLabels,
		ResourceLabels:              &CcseAttachClusterNodesResourceLabelsRequest{},
		KubeletDirectory:            "",
		ContainerDataDirectory:      "",
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
