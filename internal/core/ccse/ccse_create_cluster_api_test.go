package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseCreateClusterApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseCreateClusterApi

	// 构造请求
	var autoGenerateSecurityGroup bool = false
	var enableApiServerEip bool = false
	var enableSnat bool = false
	var installNginxIngress bool = false
	var autoRenewStatus bool = false
	var pluginCstorcsiEnabled bool = false
	var pluginCcseMonitorEnabled bool = false
	var enableAls bool = false
	var enableAlsCubeEventer bool = false
	var pluginNginxIngressEnabled bool = false
	var customScriptBase64 bool = false
	var enablePostUserScript bool = false
	var enableHostName bool = false
	var nodeUnschedulable bool = false
	var deleteProtection bool = false
	var ipvlan bool = false
	var networkPolicy bool = false
	var syncNodeLabels bool = false
	var syncNodeTaints bool = false
	var enableAffinityGroup bool = false
	var cpuManagerPolicyEnable bool = false
	var customCAEnable bool = false
	request := &CcseCreateClusterRequest{
		RegionId:  "b342b77ef26b11ecb0ac0242ac110002",
		ResPoolId: "bb9fdb42056f11eda1610242ac110002",
		ClusterBaseInfo: &CcseCreateClusterClusterBaseInfoRequest{
			SubnetUuid:                "subnet-7xipneei60",
			NetworkPlugin:             "calico",
			ClusterDomain:             "cluster.local",
			PodSubnetUuidList:         []string{},
			AutoGenerateSecurityGroup: &autoGenerateSecurityGroup,
			SecurityGroupUuid:         "sg-8judj901ou",
			StartPort:                 30000,
			EndPort:                   32767,
			EnableApiServerEip:        &enableApiServerEip,
			EnableSnat:                &enableSnat,
			NatGatewaySpec:            "small",
			ElbProdCode:               "standardI",
			NodeLabels:                &CcseCreateClusterClusterBaseInfoNodeLabelsRequest{},
			PodCidr:                   "192.168.0.0/16",
			InstallNginxIngress:       &installNginxIngress,
			NginxIngressLBSpec:        "",
			NginxIngressLBNetWork:     "internal",
			BillMode:                  "1",
			CycleType:                 "3",
			CycleCnt:                  1,
			AutoRenewStatus:           &autoRenewStatus,
			AutoRenewCycleType:        "3",
			AutoRenewCycleCount:       "1",
			ContainerRuntime:          "containerd",
			Timezone:                  "Asia/Shanghai (UTC+08:00)",
			ClusterVersion:            "1.23.3",
			DeployType:                "single",
			AzInfos: []*CcseCreateClusterClusterBaseInfoAzInfosRequest{
				{
					AzName: "cn-huabei2-tj1A-public-ctcloud",
					Size:   1,
				},
			},
			ServiceCidr:               "172.26.0.0/16",
			VpcUuid:                   "vpc-vu2t769u81",
			ClusterName:               "test-123",
			KubeProxy:                 "iptables",
			PluginCstorcsiAk:          "8d7b21fregtegeb2dda072c81b80133",
			PluginCstorcsiSk:          "54786hfvjkdsbbfrtgbtgr2dda0234454",
			PluginCstorcsiEnabled:     &pluginCstorcsiEnabled,
			PluginCcseMonitorEnabled:  &pluginCcseMonitorEnabled,
			ClusterSeries:             "cce.standard",
			ProjectId:                 "0",
			EnableAls:                 &enableAls,
			EnableAlsCubeEventer:      &enableAlsCubeEventer,
			PluginNginxIngressEnabled: &pluginNginxIngressEnabled,
			CustomScriptBase64:        &customScriptBase64,
			HostScript:                "",
			EnablePostUserScript:      &enablePostUserScript,
			PostUserScript:            "",
			EnableHostName:            &enableHostName,
			HostNamePrefix:            "",
			HostNamePostfix:           "",
			NodeTaints:                "",
			NodeUnschedulable:         &nodeUnschedulable,
			ClusterDesc:               "",
			DeleteProtection:          &deleteProtection,
			SeriesType:                "",
			Ipvlan:                    &ipvlan,
			NetworkPolicy:             &networkPolicy,
			IpStackType:               "ipv4",
			ServiceCidrV6:             "",
			ClusterAlias:              "cce-test",
			NodeScale:                 "50",
			SyncNodeLabels:            &syncNodeLabels,
			SyncNodeTaints:            &syncNodeTaints,
			EnableAffinityGroup:       &enableAffinityGroup,
			AffinityGroupUuid:         "8c66ee54-f922-75a8-1e13-af1f682f15dc",
			DelegateName:              "ecsadmintrust",
			ResourceLabels:            &CcseCreateClusterClusterBaseInfoResourceLabelsRequest{},
			CpuManagerPolicyEnable:    &cpuManagerPolicyEnable,
			San:                       "192.168.1.1",
			CustomCAEnable:            &customCAEnable,
			CustomCA:                  "",
		},
		MasterHost: &CcseCreateClusterMasterHostRequest{
			Cpu:         4,
			Mem:         8,
			ItemDefName: "s6.xlarge.2",
			ItemDefType: "S6",
			Size:        1,
			SysDisk: &CcseCreateClusterMasterHostSysDiskRequest{
				ItemDefName: "SATA",
				Size:        120,
			},
			DataDisks: []*CcseCreateClusterMasterHostDataDisksRequest{
				{
					ItemDefName: "SATA",
					Size:        120,
					DecTypeId:   "b0d48eaf-e164-4873-9b89-e3dc85c8bed0",
				},
			},
		},
		SlaveHost: &CcseCreateClusterSlaveHostRequest{
			Cpu:         4,
			Mem:         8,
			ItemDefName: "s6.xlarge.2",
			ItemDefType: "S6",
			Size:        1,
			SysDisk: &CcseCreateClusterSlaveHostSysDiskRequest{
				ItemDefName: "SATA",
				Size:        120,
			},
			DataDisks: []*CcseCreateClusterSlaveHostDataDisksRequest{
				{
					ItemDefName: "SATA",
					Size:        120,
					DecTypeId:   "b0d48eaf-e164-4873-9b89-e3dc85c8bed0",
				},
			},
			ForeignMirrorId: "49317128-641c-4a54-b0f1-a47992933413",
			MirrorType:      1,
			MirrorName:      "CTyunOS-23.01-CCND_CCSE_40_08-x86_64",
			AzInfos: []*CcseCreateClusterSlaveHostAzInfosRequest{
				{
					AzName: "cn-huabei2-tj1A-public-ctcloud",
					Size:   1,
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
