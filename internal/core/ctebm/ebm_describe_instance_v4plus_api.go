package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmDescribeInstanceV4plusApi
/* 通过参数查询单台物理机信息
 */type EbmDescribeInstanceV4plusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmDescribeInstanceV4plusApi(client *core.CtyunClient) *EbmDescribeInstanceV4plusApi {
	return &EbmDescribeInstanceV4plusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/describe-instance",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmDescribeInstanceV4plusApi) Do(ctx context.Context, credential core.Credential, req *EbmDescribeInstanceV4plusRequest) (*EbmDescribeInstanceV4plusResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("azName", req.AzName)
	ctReq.AddParam("instanceUUID", req.InstanceUUID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmDescribeInstanceV4plusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmDescribeInstanceV4plusRequest struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string /*  实例UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */
}

type EbmDescribeInstanceV4plusResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmDescribeInstanceV4plusReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmDescribeInstanceV4plusReturnObjResponse struct {
	RegionID *string `json:"regionID"` /*  资源池ID
	 */Region *string `json:"region"` /*  资源池名
	 */AzName *string `json:"azName"` /*  可用区名
	 */ResourceID *string `json:"resourceID"` /*  资源ID
	 */InstanceUUID *string `json:"instanceUUID"` /*  物理机ID
	 */DeviceUUID *string `json:"deviceUUID"` /*  设备ID
	 */DeviceType *string `json:"deviceType"` /*  设备类型
	 */DisplayName *string `json:"displayName"` /*  物理机展示名
	 */InstanceName *string `json:"instanceName"` /*  物理机名称
	 */Description *string `json:"description"` /*  物理机描述
	 */ZabbixName *string `json:"zabbixName"` /*  zabbixName
	 */SystemVolumeRaidID *string `json:"systemVolumeRaidID"` /*  本地系统盘raid id
	 */DataVolumeRaidID *string `json:"dataVolumeRaidID"` /*  数据盘raid id
	 */ImageID *string `json:"imageID"` /*  镜像ID
	 */OsType int32 `json:"osType"` /*  操作系统类型编号
	 */OsTypeName *string `json:"osTypeName"` /*  操作系统类型
	 */VpcID *string `json:"vpcID"` /*  主网卡网络ID
	 */VpcName *string `json:"vpcName"` /*  私有云名称
	 */SubnetID *string `json:"subnetID"` /*  主网卡子网id
	 */PublicIP *string `json:"publicIP"` /*  公网IPIPv4地址
	 */PrivateIP *string `json:"privateIP"` /*  主网卡私有IPv4地址
	 */PublicIPv6 *string `json:"publicIPv6"` /*  公网IPv6地址
	 */PrivateIPv6 *string `json:"privateIPv6"` /*  私有IPv6地址
	 */EbmState *string `json:"ebmState"` /*  物理机状态
	 */Flavor *EbmDescribeInstanceV4plusReturnObjFlavorResponse `json:"flavor"` /*  规格信息
	 */Interfaces []*EbmDescribeInstanceV4plusReturnObjInterfacesResponse `json:"interfaces"` /*  网卡信息
	//NetworkInfo *EbmDescribeInstanceV4plusReturnObjNetworkInfoResponse `json:"networkInfo"` /*  网络信息
	*/RaidDetail *EbmDescribeInstanceV4plusReturnObjRaidDetailResponse `json:"raidDetail"` /*  磁盘阵列信息
	 */AttachedVolumes []*string `json:"attachedVolumes"` /*  挂载的硬盘ID
	 */DeviceDetail *EbmDescribeInstanceV4plusReturnObjDeviceDetailResponse `json:"deviceDetail"` /*  设备信息
	 */Freezing *bool `json:"freezing"` /*  是否冻结
	 */Expired *bool `json:"expired"` /*  是否到期
	 */ReleaseDate *string `json:"releaseDate"` /*  预期释放时间
	 */CreateTime *string `json:"createTime"` /*  创建时间
	 */UpdatedTime *string `json:"updatedTime"` /*  最后更新时间
	 */ExpiredTime *string `json:"expiredTime"` /*  到期时间
	 */OnDemand *bool `json:"onDemand"` /*  付费方式，true表示按量付费，false为包周期
	 */
}

type EbmDescribeInstanceV4plusReturnObjFlavorResponse struct {
	NumaNodeAmount int32 `json:"numaNodeAmount"` /*  cpu numa
	 */NicAmount int32 `json:"nicAmount"` /*  网卡
	 */MemSize int32 `json:"memSize"` /*  内存大小
	 */Ram int32 `json:"ram"` /*  内存数量
	 */Vcpus int32 `json:"vcpus"` /*  vCPU数量
	 */CpuThreadAmount int32 `json:"cpuThreadAmount"` /*  cpu进程数目
	 */DeviceType *string `json:"deviceType"` /*  设备类型
	 */
}

type EbmDescribeInstanceV4plusReturnObjInterfacesResponse struct {
	InterfaceUUID *string `json:"interfaceUUID"` /*  网卡UUID
	 */Master *bool `json:"master"` /*  是否为主网卡
	 */VpcUUID *string `json:"vpcUUID"` /*  网络UUID
	 */SubnetUUID *string `json:"subnetUUID"` /*  子网UUID
	 */PortUUID *string `json:"portUUID"` /*  Port UUID
	 */Ipv4 *string `json:"ipv4"` /*  IPv4地址
	 */Ipv4Gateway *string `json:"ipv4Gateway"` /*  IPv4网关
	 */Ipv6 *string `json:"ipv6"` /*  IPv6地址
	 */Ipv6Gateway *string `json:"ipv6Gateway"` /*  IPv4网关
	 */VipUUIDList []*string `json:"vipUUIDList"` /*  vip UUID列表
	 */VipList []*string `json:"vipList"` /*  vip列表
	 */SecurityGroups []*EbmDescribeInstanceV4plusReturnObjInterfacesSecurityGroupsResponse `json:"securityGroups"` /*  安全组信息列表
	 */
}

type EbmDescribeInstanceV4plusReturnObjNetworkInfoResponse struct {
	SubnetUUID *string `json:"subnetUUID"` /*  c45512bf-7919-55a9-a106-5f4aa9194c7c
	 */VpcName *string `json:"vpcName"` /*  vpc-f60a
	 */VpcID *string `json:"vpcID"` /*  a64ddb73-3021-5a8d-9abd-b0a12e429690
	 */
}

type EbmDescribeInstanceV4plusReturnObjRaidDetailResponse struct {
	SystemVolume *EbmDescribeInstanceV4plusReturnObjRaidDetailSystemVolumeResponse `json:"systemVolume"` /*  系统盘
	 */DataVolume *EbmDescribeInstanceV4plusReturnObjRaidDetailDataVolumeResponse `json:"dataVolume"` /*  数据盘
	 */
}

type EbmDescribeInstanceV4plusReturnObjDeviceDetailResponse struct {
	CpuSockets int32 `json:"cpuSockets"` /*  物理cpu数量
	 */NumaNodeAmount int32 `json:"numaNodeAmount"` /*  单个cpu numa node数量
	 */CpuAmount int32 `json:"cpuAmount"` /*  单个cpu核数
	 */CpuThreadAmount int32 `json:"cpuThreadAmount"` /*  单个cpu核超线程数量
	 */CpuManufacturer *string `json:"cpuManufacturer"` /*  cpu厂商；Intel，AMD，Hygon，HiSilicon，Loongson等
	 */CpuModel *string `json:"cpuModel"` /*  cpu型号
	 */CpuFrequency *string `json:"cpuFrequency"` /*  cpu频率(G)
	 */MemAmount int32 `json:"memAmount"` /*  内存数
	 */MemSize int32 `json:"memSize"` /*  内存大小(G)
	 */MemFrequency int32 `json:"memFrequency"` /*  内存频率(MHz)
	 */NicAmount int32 `json:"nicAmount"` /*  网卡数
	 */NicRate int32 `json:"nicRate"` /*  网卡传播速率(GE)
	 */SystemVolumeAmount int32 `json:"systemVolumeAmount"` /*  系统盘数量
	 */SystemVolumeSize int32 `json:"systemVolumeSize"` /*  系统盘单盘大小(GB)
	 */SystemVolumeType *string `json:"systemVolumeType"` /*  系统盘介质类型； 包含SSD、HDD
	 */SystemVolumeInterface *string `json:"systemVolumeInterface"` /*  系统盘接口类型；包含SAS、SATA、NVMe
	 */SystemVolumeDescription *string `json:"systemVolumeDescription"` /*  系统盘描述
	 */DataVolumeAmount int32 `json:"dataVolumeAmount"` /*  数据盘数量
	 */DataVolumeSize int32 `json:"dataVolumeSize"` /*  数据盘单盘大小(GB)
	 */DataVolumeInterface *string `json:"dataVolumeInterface"` /*  数据盘接口；包含SAS、SATA、NVMe
	 */DataVolumeType *string `json:"dataVolumeType"` /*  数据盘介质类型； 包含SSD、HDD
	 */DataVolumeDescription *string `json:"dataVolumeDescription"` /*  数据盘描述
	 */SmartNicExist *bool `json:"smartNicExist"` /*  是否有智能网卡，true为弹性裸金属; false为标准裸金属
	 */NvmeVolumeAmount int32 `json:"nvmeVolumeAmount"` /*  NVME硬盘数量
	 */NvmeVolumeSize int32 `json:"nvmeVolumeSize"` /*  NVME硬盘单盘大小(GB)
	 */NvmeVolumeType *string `json:"nvmeVolumeType"` /*  NVME介质类型； 包含SSD、HDD
	 */NvmeVolumeInterface *string `json:"nvmeVolumeInterface"` /*  NVME接口类型；包含SAS、SATA、NVMe
	 */GpuAmount int32 `json:"gpuAmount"` /*  GPU数目
	 */GpuSize int32 `json:"gpuSize"` /*  GPU显存
	 */GpuManufacturer *string `json:"gpuManufacturer"` /*  GPU厂商；Nvidia，Huawei，Cambricon等
	 */GpuModel *string `json:"gpuModel"` /*  GPU型号
	 */SupportCloud *bool `json:"supportCloud"` /*  是否支持云盘
	 */CloudBoot *bool `json:"cloudBoot"` /*  是否支持云盘启动
	 */
}

type EbmDescribeInstanceV4plusReturnObjInterfacesSecurityGroupsResponse struct {
	SecurityGroupID *string `json:"securityGroupID"` /*  71386230-c9e7-5465-85fc-29e26f79806d
	 */SecurityGroupName *string `json:"securityGroupName"` /*  Default-Security-Group
	 */
}

type EbmDescribeInstanceV4plusReturnObjRaidDetailSystemVolumeResponse struct {
	Uuid *string `json:"uuid"` /*  UUID
	 */VolumeType *string `json:"volumeType"` /*  卷类型
	 */Name *string `json:"name"` /*  raid名称
	 */VolumeDetail *string `json:"volumeDetail"` /*  卷详情
	 */
}

type EbmDescribeInstanceV4plusReturnObjRaidDetailDataVolumeResponse struct {
	Uuid *string `json:"uuid"` /*  UUID
	 */VolumeType *string `json:"volumeType"` /*  卷类型
	 */Name *string `json:"name"` /*  raid名称
	 */VolumeDetail *string `json:"volumeDetail"` /*  卷详情
	 */
}
