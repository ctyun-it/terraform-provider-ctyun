package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtimageDetailImageApi
/* 根据镜像 ID，查询 1 份镜像的详细信息<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：<br />1. 推荐使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=4577&data=97&isNormal=1&vid=91" target="_blank">查询物理机镜像</a>接口来查询物理机镜像<br />2. 在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分
 */type CtimageDetailImageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageDetailImageApi(client *core.CtyunClient) *CtimageDetailImageApi {
	return &CtimageDetailImageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/image/detail",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageDetailImageApi) Do(ctx context.Context, credential core.Credential, req *CtimageDetailImageRequest) (*CtimageDetailImageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("imageID", req.ImageID)
	ctReq.AddParam("regionID", req.RegionID)
	if req.ErrorFree != nil {
		ctReq.AddParam("errorFree", strconv.FormatBool(*req.ErrorFree))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtimageDetailImageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageDetailImageRequest struct {
	ImageID   string `json:"imageID,omitempty"`  /*  镜像 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您可使用的镜像资源  */
	RegionID  string `json:"regionID,omitempty"` /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表  */
	ErrorFree *bool  `json:"errorFree"`          /*  用于表示是否期望特定场景“零错误”响应的标识。默认 true。注意：<br />1. 特定场景是指传入的 imageID 参数所指定的镜像不是您在传入的 regionID 参数所指定的资源池中可使用的镜像资源<br />2. 对于特定场景，若此参数设置为 false（推荐），则此接口会响应失败 statusCode、相应的 error/errorCode 等；否则，此接口会响应成功 statusCode 等，但 returnObj 中的 images 是空的镜像列表<br />3. 此参数已弃用，目前仍可使用，但会在合适的时机移除而不再允许设置。移除后，此接口的响应行为等效于将此参数设置为 false（目前默认 true 是为了在移除前的过度阶段仍保持旧有行为），因此请您尽快适配设置此参数为 false 的情形  */
}

type CtimageDetailImageResponse struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功，<br />900：失败  */
	Error       string                               `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                               `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                               `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageDetailImageReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageDetailImageReturnObjResponse struct {
	Images []*CtimageDetailImageReturnObjImagesResponse `json:"images"` /*  镜像列表。应包含 1 个镜像  */
}

type CtimageDetailImageReturnObjImagesResponse struct {
	AppVersion                string `json:"appVersion,omitempty"`                /*  应用版本  */
	Architecture              string `json:"architecture,omitempty"`              /*  系统架构。取值范围（值：描述）：<br />aarch64：AArch64 架构，<br />loongarch64：LoongArch64 架构，<br />sw_64：sw_64 架构，<br />x86_64：x86_64 架构  */
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
