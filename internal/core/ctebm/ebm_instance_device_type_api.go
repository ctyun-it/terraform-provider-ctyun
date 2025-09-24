package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmInstanceDeviceTypeApi
/* 根据实例ID查询对应套餐信息
 */type EbmInstanceDeviceTypeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmInstanceDeviceTypeApi(client *core.CtyunClient) *EbmInstanceDeviceTypeApi {
	return &EbmInstanceDeviceTypeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/instance-device-type",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmInstanceDeviceTypeApi) Do(ctx context.Context, credential core.Credential, req *EbmInstanceDeviceTypeRequest) (*EbmInstanceDeviceTypeResponse, error) {
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
	var resp EbmInstanceDeviceTypeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmInstanceDeviceTypeRequest struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string /*  实例UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */
}

type EbmInstanceDeviceTypeResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmInstanceDeviceTypeReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmInstanceDeviceTypeReturnObjResponse struct {
	Id int32 `json:"id"` /*  套餐ID
	 */Region *string `json:"region"` /*  资源池
	 */AzName *string `json:"azName"` /*  可用区
	 */DeviceType *string `json:"deviceType"` /*  套餐类型
	 */NameZh *string `json:"nameZh"` /*  物理机中文名
	 */NameEn *string `json:"nameEn"` /*  英文名
	 */CpuSockets int32 `json:"cpuSockets"` /*  物理cpu数量
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
	 */SmartNicExist *bool `json:"smartNicExist"` /*  是否有智能网卡
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
	 */EnableShadowVpc *bool `json:"enableShadowVpc"` /*  是否支持存储高速网络；如支持存储高速网络则会占用对应可用网卡数量
	 */Project *string `json:"project"` /*  项目信息
	 */CreateTime *string `json:"createTime"` /*  创建时间(epoch 时间格式)
	 */UpdateTime *string `json:"updateTime"` /*  最后更新时间(epoch 时间格式)
	 */
}
