package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbmListInstanceV4plusApi
/* 通过参数批量查询物理机信息
 */type EbmListInstanceV4plusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmListInstanceV4plusApi(client *core.CtyunClient) *EbmListInstanceV4plusApi {
	return &EbmListInstanceV4plusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/list-instance",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmListInstanceV4plusApi) Do(ctx context.Context, credential core.Credential, req *EbmListInstanceV4plusRequest) (*EbmListInstanceV4plusResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("azName", req.AzName)
	if req.ResourceID != nil {
		ctReq.AddParam("resourceID", *req.ResourceID)
	}
	if req.Ip != nil {
		ctReq.AddParam("ip", *req.Ip)
	}
	if req.InstanceName != nil {
		ctReq.AddParam("instanceName", *req.InstanceName)
	}
	if req.VpcID != nil {
		ctReq.AddParam("vpcID", *req.VpcID)
	}
	if req.SubnetID != nil {
		ctReq.AddParam("subnetID", *req.SubnetID)
	}
	if req.DeviceType != nil {
		ctReq.AddParam("deviceType", *req.DeviceType)
	}
	if req.DeviceUUIDList != nil {
		ctReq.AddParam("deviceUUIDList", *req.DeviceUUIDList)
	}
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if req.InstanceUUIDList != nil {
		ctReq.AddParam("instanceUUIDList", *req.InstanceUUIDList)
	}
	if req.InstanceUUID != nil {
		ctReq.AddParam("instanceUUID", *req.InstanceUUID)
	}
	if req.Status != nil {
		ctReq.AddParam("status", *req.Status)
	}
	if req.Sort != nil {
		ctReq.AddParam("sort", *req.Sort)
	}
	if req.Asc != nil {
		ctReq.AddParam("asc", strconv.FormatBool(*req.Asc))
	}
	if req.VipID != nil {
		ctReq.AddParam("vipID", *req.VipID)
	}
	if req.VolumeUUID != nil {
		ctReq.AddParam("volumeUUID", *req.VolumeUUID)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmListInstanceV4plusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmListInstanceV4plusRequest struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */ResourceID *string /*  物理机资源ID
	 */Ip *string /*  弹性ip，公网IP地址
	 */InstanceName *string /*  实例名称
	 */VpcID *string /*  虚拟私有云ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-虚拟私有云</a>来了解虚拟私有云<br /> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94">查询VPC列表</a><br /><a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4811&data=94">创建VPC</a><br />注：在多可用区类型资源池下，vpcID通常为“vpc-”开头，非多可用区类型资源池vpcID为uuid格式
	 */SubnetID *string /*  子网UUID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10098380">基本概念</a>来查找子网的相关定义<br /> <a href="https://www.ctyun.cn/document/10026755/10040797">查询子网列表</a><br /><a href="https://www.ctyun.cn/document/10026755/10040804">创建子网</a><br/>注：在多可用区类型资源池下，subnetID通常以“subnet-”开头；非多可用区类型资源池subnetID为uuid格式
	 */DeviceType *string /*  物理机套餐类型<br /><a href="https://www.ctyun.cn/document/10027724/10040107">查询资源池内物理机套餐</a><br /><a href="https://www.ctyun.cn/document/10027724/10040124">查询指定物理机的套餐信息</a>
	 */DeviceUUIDList *string /*  设备uuid 用,分隔，您可以调用<a href="https://www.ctyun.cn/document/10027724/10040100">查询单台物理机</a>和<a href="https://www.ctyun.cn/document/10027724/10040106">批量查询物理机</a>获取设备uuid
	 */QueryContent *string /*  对instanceName，内网IP，displayName这些字段模糊查询
	 */InstanceUUIDList *string /*  实例UUID 用,分隔，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */InstanceUUID *string /*  实例UUID，您可以调用<a href="https://www.ctyun.cn/document/10026730/10069118">根据订单号查询uuid</a>获取实例UUID
	 */Status *string /*  实例状态，取值范围：CREATING(创建中)，STARTING(开机中)，RUNNING(运行中)，STOPPING(关机中)，RESTARTING(重启中)，ERROR(故障中)，REINSTALLING(重装系统中)，RESETTING_PASSWORD(重置密码中)，ADDING_NETWORK(添加网卡中)，DELETING_NETWORK(删除网卡中)<br>注：该参数大小写不敏感（如CREATING可填写为creating）
	 */Sort *string /*  排序类型。取值范围：[expire_time]。expire_time表示按到期时间排序，默认为降序排序
	 */Asc *bool /*  排序参数。true表示升序，false表示降序。当未指定排序类型 sort 时，此参数无效
	 */VipID *string /*  vip_id
	 */VolumeUUID *string /*  云硬盘UUID
	 */PageNo int32 /*  页码，默认值:1
	 */PageSize int32 /*  每页记录数目，取值范围:[1~10000]，默认值:10，单页最大记录不超过10000
	 */
}

type EbmListInstanceV4plusResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmListInstanceV4plusReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmListInstanceV4plusReturnObjResponse struct {
	CurrentCount int32 `json:"currentCount"` /*  当前页数量
	 */TotalCount int32 `json:"totalCount"` /*  总记录数
	 */TotalPage int32 `json:"totalPage"` /*  总页数
	 */Results []*EbmListInstanceV4plusReturnObjResultsResponse `json:"results"` /*  分页明细,元素类型是results,定义请参考表results
	 */
}

type EbmListInstanceV4plusReturnObjResultsResponse struct {
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
	 */Flavor *EbmListInstanceV4plusReturnObjResultsFlavorResponse `json:"flavor"` /*  规格信息
	 */Interfaces []*EbmListInstanceV4plusReturnObjResultsInterfacesResponse `json:"interfaces"` /*  网卡信息
	 */NetworkInfo []*EbmListInstanceV4plusReturnObjResultsNetworkInfoResponse `json:"networkInfo"` /*  网络信息
	 */RaidDetail *EbmListInstanceV4plusReturnObjResultsRaidDetailResponse `json:"raidDetail"` /*  磁盘阵列信息
	 */AttachedVolumes []*string `json:"attachedVolumes"` /*  挂载的硬盘ID
	 */DeviceDetail *EbmListInstanceV4plusReturnObjResultsDeviceDetailResponse `json:"deviceDetail"` /*  设备信息
	 */Freezing *bool `json:"freezing"` /*  是否冻结
	 */Expired *bool `json:"expired"` /*  是否到期
	 */ReleaseDate *string `json:"releaseDate"` /*  预期释放时间
	 */CreateTime *string `json:"createTime"` /*  创建时间
	 */UpdatedTime *string `json:"updatedTime"` /*  最后更新时间
	 */ExpiredTime *string `json:"expiredTime"` /*  到期时间
	 */OnDemand *bool `json:"onDemand"` /*  付费方式，true表示按量付费; false为包周期
	 */
}

type EbmListInstanceV4plusReturnObjResultsFlavorResponse struct {
	NumaNodeAmount int32 `json:"numaNodeAmount"` /*  cpu numa
	 */NicAmount int32 `json:"nicAmount"` /*  网卡
	 */MemSize int32 `json:"memSize"` /*  内存大小
	 */Ram int32 `json:"ram"` /*  内存数量
	 */Vcpus int32 `json:"vcpus"` /*  vCPU数量
	 */CpuThreadAmount int32 `json:"cpuThreadAmount"` /*  cpu进程数目
	 */DeviceType *string `json:"deviceType"` /*  设备类型
	 */
}

type EbmListInstanceV4plusReturnObjResultsInterfacesResponse struct {
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
	 */SecurityGroups []*EbmListInstanceV4plusReturnObjResultsInterfacesSecurityGroupsResponse `json:"securityGroups"` /*  安全组信息列表
	 */
}

type EbmListInstanceV4plusReturnObjResultsNetworkInfoResponse struct {
	SubnetUUID *string `json:"subnetUUID"` /*  子网ID
	 */VpcName *string `json:"vpcName"` /*  私有云名称
	 */VpcID *string `json:"vpcID"` /*  私有云ID
	 */
}

type EbmListInstanceV4plusReturnObjResultsRaidDetailResponse struct {
	SystemVolume *EbmListInstanceV4plusReturnObjResultsRaidDetailSystemVolumeResponse `json:"systemVolume"` /*  系统盘
	 */DataVolume *EbmListInstanceV4plusReturnObjResultsRaidDetailDataVolumeResponse `json:"dataVolume"` /*  数据盘
	 */
}

type EbmListInstanceV4plusReturnObjResultsDeviceDetailResponse struct {
	CpuSockets int32 `json:"cpuSockets"` /*  物理cpu数量
	 */NumaNodeAmount int32 `json:"numaNodeAmount"` /*  单个cpu numa node数量
	 */CpuAmount int32 `json:"cpuAmount"` /*  单个cpu核数
	 */CpuThreadAmount int32 `json:"cpuThreadAmount"` /*  单个cpu核超线程数量
	 */CpuManufacturer *string `json:"cpuManufacturer"` /*  cpu厂商；Intel ，AMD，Hygon，HiSilicon，loongson等
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
	 */DataVolumeType *string `json:"dataVolumeType"` /*  数据盘介质参数类型； 包含SSD、HDD
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
	 */ComputeIBAmount int32 `json:"computeIBAmount"` /*  计算ib网卡大小
	 */ComputeIBRate int32 `json:"computeIBRate"` /*  计算ib网卡速率(GE)
	 */StorageIBAmount int32 `json:"storageIBAmount"` /*  存储ib网卡大小
	 */StorageIBRate int32 `json:"storageIBRate"` /*  存储ib网卡速率(GE)
	 */ComputeRoCEAmount int32 `json:"computeRoCEAmount"` /*  计算RoCE网卡大小
	 */ComputeRoCERate int32 `json:"computeRoCERate"` /*  计算RoCE网卡速率(GE)
	 */StorageRoCEAmount int32 `json:"storageRoCEAmount"` /*  存储RoCE网卡大小
	 */StorageRoCERate int32 `json:"storageRoCERate"` /*  存储RoCE网卡速率(GE)
	 */SupportCloud *bool `json:"supportCloud"` /*  是否支持云盘
	 */CloudBoot *bool `json:"cloudBoot"` /*  是否支持云盘启动
	 */
}

type EbmListInstanceV4plusReturnObjResultsInterfacesSecurityGroupsResponse struct {
	SecurityGroupID *string `json:"securityGroupID"` /*  安全组UUID
	 */SecurityGroupName *string `json:"securityGroupName"` /*  安全组名称
	 */
}

type EbmListInstanceV4plusReturnObjResultsRaidDetailSystemVolumeResponse struct {
	Uuid *string `json:"uuid"` /*  UUID
	 */VolumeType *string `json:"volumeType"` /*  卷类型
	 */Name *string `json:"name"` /*  raid名称
	 */VolumeDetail *string `json:"volumeDetail"` /*  卷详情
	 */
}

type EbmListInstanceV4plusReturnObjResultsRaidDetailDataVolumeResponse struct {
	Uuid *string `json:"uuid"` /*  UUID
	 */VolumeType *string `json:"volumeType"` /*  卷类型
	 */Name *string `json:"name"` /*  英文名称
	 */VolumeDetail *string `json:"volumeDetail"` /*  卷详情
	 */
}
