package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageCreateEcsDataDiskImageApi
/* 使用指定的云主机数据盘来创建 1 份私有数据盘镜像<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分
 */type CtimageCreateEcsDataDiskImageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageCreateEcsDataDiskImageApi(client *core.CtyunClient) *CtimageCreateEcsDataDiskImageApi {
	return &CtimageCreateEcsDataDiskImageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/image/create-from-data-disk",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageCreateEcsDataDiskImageApi) Do(ctx context.Context, credential core.Credential, req *CtimageCreateEcsDataDiskImageRequest) (*CtimageCreateEcsDataDiskImageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtimageCreateEcsDataDiskImageRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtimageCreateEcsDataDiskImageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageCreateEcsDataDiskImageRequest struct {
	DataDiskID                string                                        `json:"dataDiskID,omitempty"`      /*  数据盘 ID。注意：所指定的数据盘应是挂载在传入的 instanceID 参数所指定的云主机、磁盘模式不为 FCSAN 或 ISCSI、磁盘大小不大于您所能创建的 1 份私有数据盘镜像的最大大小、云硬盘状态为 in-use、未加密的云硬盘。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8290&data=87&isNormal=1&vid=81" target="_blank">查询云主机的云硬盘列表</a>接口来查询 1 台云主机挂载的云硬盘，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7336&data=48&isNormal=1&vid=45" target="_blank">查询云硬盘详情（基于 diskID）</a>接口来查询 1 块云硬盘的详细信息  */
	ImageName                 string                                        `json:"imageName,omitempty"`       /*  镜像名称。注意：<br />1. 长度为 2~32 个字符，只能由数字、字母、- 组成，不能以数字、- 开头，且不能以 - 结尾<br />2. 不能与您已有的私有镜像的名称重复。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您已有的私有镜像  */
	InstanceID                string                                        `json:"instanceID,omitempty"`      /*  云主机 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87&isNormal=1&vid=81" target="_blank">查询云主机列表</a>接口来查询您的云主机列表，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8310&data=87&isNormal=1&vid=81" target="_blank">查询一台云主机详细信息</a>接口来查询 1 台云主机的详细信息。注意：所指定的云主机应是云主机状态为 stopped 的云主机。仅在传入的 regionID 参数所指定的资源池具备运行中的云主机创建私有镜像的功能时可指定云主机状态为 running 的云主机  */
	RegionID                  string                                        `json:"regionID,omitempty"`        /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表。注意：此接口仅支持具备从云主机数据盘创建私有镜像的功能的资源池  */
	Description               string                                        `json:"description,omitempty"`     /*  描述信息。注意：长度为 1~128 个字符，不能以空格开头或结尾  */
	EnableImageIntegrityCheck *bool                                         `json:"enableImageIntegrityCheck"` /*  用于表示是否启用镜像完整性校验的标识。默认 false。注意：仅在传入的 regionID 参数所指定的资源池具备镜像完整性校验功能时生效  */
	Labels                    []*CtimageCreateEcsDataDiskImageLabelsRequest `json:"labels"`                    /*  标签列表。注意：<br />1. 列表中最多 10 个标签<br />2. 标签键不可重复<br />3. 单个标签键或值应满足长度为 1~32 个字符，不能换行，且不能以空格开头或结尾  */
	ProjectID                 string                                        `json:"projectID,omitempty"`       /*  企业项目 ID。默认 0（即 default 企业项目）。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=7246&data=114&isNormal=1&vid=107" target="_blank">查询企业项目列表</a>接口来查询您可以使用的企业项目 ID  */
}

type CtimageCreateEcsDataDiskImageLabelsRequest struct {
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签键  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签值  */
}

type CtimageCreateEcsDataDiskImageResponse struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功<br />900：失败  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                                          `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                          `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageCreateEcsDataDiskImageReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageCreateEcsDataDiskImageReturnObjResponse struct {
	Images []*CtimageCreateEcsDataDiskImageReturnObjImagesResponse `json:"images"` /*  镜像列表。应包含 1 个镜像  */
}

type CtimageCreateEcsDataDiskImageReturnObjImagesResponse struct {
	AppVersion                string `json:"appVersion,omitempty"`                /*  应用版本  */
	Architecture              string `json:"architecture,omitempty"`              /*  系统架构。取值范围与对应的请求参数相同  */
	AzName                    string `json:"azName,omitempty"`                    /*  在多可用区资源池下物理机镜像的可用区名称  */
	BootMode                  string `json:"bootMode,omitempty"`                  /*  x86_64 架构非数据盘镜像的启动方式。取值范围（值：描述）：<br />bios：BIOS 启动方式，<br />uefi：UEFI 启动方式  */
	ChargeableImage           *bool  `json:"chargeableImage"`                     /*  用于表示是否是收费镜像的标识  */
	ContainerFormat           string `json:"containerFormat,omitempty"`           /*  容器格式  */
	CreatedTime               int32  `json:"createdTime,omitempty"`               /*  创建时间戳  */
	CreatedTimeStr            string `json:"createdTimeStr,omitempty"`            /*  创建时间  */
	CwaiType                  string `json:"cwaiType,omitempty"`                  /*  云骁智算云主机节点类型。取值范围（值：描述）：<br />control：控制面云主机节点，<br />node：GPU 云主机节点<br /><br />注意：镜像可适用于多节点类型，多个云骁智算云主机节点类型之间以英文逗号（,）分隔，如 control,node  */
	Description               string `json:"description,omitempty"`               /*  描述信息  */
	DestinationAccountID      string `json:"destinationAccountID,omitempty"`      /*  共享镜像接受者的账号 ID  */
	DestinationUser           string `json:"destinationUser,omitempty"`           /*  共享镜像接受者  */
	DiskFormat                string `json:"diskFormat,omitempty"`                /*  磁盘格式。取值范围（值：描述）：<br />qcow2：QCOW2 格式，<br />raw：RAW 格式，<br />vhd：VHD 格式，<br />vmdk：VMDK 格式  */
	DiskID                    string `json:"diskID,omitempty"`                    /*  私有镜像来源的系统盘/数据盘 ID  */
	DiskSize                  int32  `json:"diskSize,omitempty"`                  /*  磁盘容量。单位为 GiB  */
	EnableImageIntegrityCheck *bool  `json:"enableImageIntegrityCheck"`           /*  用于表示是否启用镜像完整性校验的标识  */
	FullECSDiskSize           int32  `json:"fullECSDiskSize,omitempty"`           /*  云主机整机磁盘容量。单位为 GiB  */
	GpuImageCategory          string `json:"gpuImageCategory,omitempty"`          /*  GPU 镜像种类。取值范围（值：描述）：<br />pass_through：GPU 直通镜像，<br />vgpu：vGPU 镜像  */
	HasAcceptedSharedImages   *bool  `json:"hasAcceptedSharedImages"`             /*  用于表示私有镜像的共享列表中是否有镜像状态为 accepted 的共享镜像的标识  */
	ImageClass                string `json:"imageClass,omitempty"`                /*  镜像类别。取值范围（值：描述）：<br />BMS：物理机，<br />ECS：云主机  */
	ImageDisplayName          string `json:"imageDisplayName,omitempty"`          /*  镜像展示名称  */
	ImageID                   string `json:"imageID,omitempty"`                   /*  镜像 ID  */
	ImageIntegrityCheckStatus string `json:"imageIntegrityCheckStatus,omitempty"` /*  镜像完整性校验状态，详见枚举值表格  */
	ImageName                 string `json:"imageName,omitempty"`                 /*  镜像名称  */
	ImageScene                string `json:"imageScene,omitempty"`                /*  镜像场景。取值范围（值：描述）：<br />dev：开发工具，<br />ecommerce：电商，<br />gaming：游戏，<br />website：网站<br /><br />注意：镜像可适用于多场景，多个镜像场景之间以英文逗号（,）分隔，如 ecommerce,website  */
	ImageShareCount           int32  `json:"imageShareCount,omitempty"`           /*  私有镜像的共享数量  */
	ImageSize                 int64  `json:"imageSize,omitempty"`                 /*  镜像大小。单位为 byte  */
	ImageSource               string `json:"imageSource,omitempty"`               /*  私有镜像来源。取值范围（值：描述）：<br />cloud_server：云主机，<br />full_ecs：云主机整机，<br />image_file：镜像文件，<br />metal_server：物理机，<br />snapshot：云主机快照  */
	ImageStatus               string `json:"imageStatus,omitempty"`               /*  镜像状态，详见枚举值表格  */
	ImageSubcategory          string `json:"imageSubcategory,omitempty"`          /*  镜像子种类。取值范围（值：描述）：<br />app：云主机应用镜像，<br />thin_app：轻量型云主机应用镜像<br /><br />注意：镜像可适用于多子种类，多个镜像子种类之间以英文逗号（,）分隔，如 app,thin_app  */
	ImageType                 string `json:"imageType,omitempty"`                 /*  镜像类型。取值范围（值：描述）：<br />（空，即 null）：系统盘镜像，<br />data_disk_image：数据盘镜像，<br />full_ecs_image：整机镜像，<br />iso_image：ISO 镜像  */
	ImageVisibility           string `json:"imageVisibility,omitempty"`           /*  镜像可见类型，详见枚举值表格  */
	MaximumRAM                int32  `json:"maximumRAM,omitempty"`                /*  最大内存。单位为 GiB  */
	MinimumRAM                int32  `json:"minimumRAM,omitempty"`                /*  最小内存。单位为 GiB  */
	OsDistro                  string `json:"osDistro,omitempty"`                  /*  操作系统发行版  */
	OsType                    string `json:"osType,omitempty"`                    /*  操作系统类型。取值范围（值：描述）：<br />linux：Linux 系操作系统，<br />windows：Windows 系操作系统  */
	OsVersion                 string `json:"osVersion,omitempty"`                 /*  操作系统版本  */
	ProjectID                 string `json:"projectID,omitempty"`                 /*  企业项目 ID  */
	SourceAccountID           string `json:"sourceAccountID,omitempty"`           /*  共享镜像提供者的账号 ID  */
	SourceServerID            string `json:"sourceServerID,omitempty"`            /*  私有镜像来源的云主机/云主机快照/物理机 ID  */
	SourceUser                string `json:"sourceUser,omitempty"`                /*  共享镜像提供者  */
	SupportOneClickSFSMount   *bool  `json:"supportOneClickSFSMount"`             /*  用于表示是否支持一键挂载文件系统的标识  */
	SupportXSSD               *bool  `json:"supportXSSD"`                         /*  用于表示是否支持 XSSD 类型盘的标识  */
	TaskID                    string `json:"taskID,omitempty"`                    /*  任务 ID  */
	UpdatedTime               int32  `json:"updatedTime,omitempty"`               /*  更新时间戳  */
	UpdatedTimeStr            string `json:"updatedTimeStr,omitempty"`            /*  更新时间  */
}
