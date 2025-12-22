package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageImportImageApi
/* 使用存在对象存储的桶上的指定的镜像文件来创建 1 份云主机私有镜像<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：<br />1. 在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分<br />2. 此接口在多可用区资源池中可响应创建中的私有镜像的 imageID。此时，若您欲删除镜像状态为 error 的私有镜像，则可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4766&data=89&isNormal=1&vid=83" target="_blank">删除私有镜像</a>接口来删除 1 份私有镜像<br />3. 此接口在非多可用区资源池中可响应创建中的私有镜像的 taskID。此时，若您欲删除镜像状态为 error 的私有镜像，则可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=7489&data=89&isNormal=1&vid=83" target="_blank">删除创建私有镜像（镜像文件）任务</a>接口来删除 1 个失败的使用镜像文件来创建私有镜像的任务
 */type CtimageImportImageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageImportImageApi(client *core.CtyunClient) *CtimageImportImageApi {
	return &CtimageImportImageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/image/import",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageImportImageApi) Do(ctx context.Context, credential core.Credential, req *CtimageImportImageRequest) (*CtimageImportImageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtimageImportImageRequest
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
	var resp CtimageImportImageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageImportImageRequest struct {
	ImageFileSource           string                                    `json:"imageFileSource,omitempty"` /*  镜像文件地址，即存储桶内对象的 URL。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5305&data=89&vid=83" target="_blank">校验镜像文件地址</a>接口来校验 1 个镜像文件地址。注意：<br />1. 格式应为 {internetEndpoint}/{bucket}/{key}。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5305&data=89&isNormal=1&vid=83" target="_blank">访问控制 endpoint 查询</a>接口来查询外网 endpoint，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6247&data=105&isNormal=1&vid=99" target="_blank">查询所有桶</a>接口来查询您拥有的桶的列表，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6271&data=105&isNormal=1&vid=99" target="_blank">查询桶信息</a>接口来查询 1 个桶的详细信息，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6273&data=105&isNormal=1&vid=99" target="_blank">查看对象列表</a>接口来查询存储桶内所有对象，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6282&data=105&isNormal=1&vid=99" target="_blank">查询对象是否存在</a>接口来查询 1 个对象的详细信息，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6290&data=105&isNormal=1&vid=99" target="_blank">获取对象 ACL </a>接口来查询 1 个对象的访问权限控制列表<br />2. 所指向的桶应是属于您的、存储类型为 STANDARD、未加密的桶<br />3. 所指向的对象应是所指向的桶内的存储类型为 STANDARD、未加密的对象。此对象在非多可用区资源池中还应可公共读。您需自行确保此对象是 QCOW2、RAW、VHD 或 VMDK 格式的镜像文件<br />4. 在 imageProperties.imageType 参数值为空（系统盘镜像）时，所指向的对象的大小应不超过您所能创建的 1 份私有系统盘镜像的最大大小<br />5. 在 imageProperties.imageType 参数值为 data_disk_image 时，所指向的对象的大小应不超过您所能创建的 1 份私有数据盘镜像的最大大小<br />6. 在 imageProperties.imageType 参数值为 iso_image 时，所指向的对象的大小应不超过 1099511627776 bytes（即 1 TiB）  */
	ImageProperties           *CtimageImportImageImagePropertiesRequest `json:"imageProperties"`           /*  镜像属性  */
	RegionID                  string                                    `json:"regionID,omitempty"`        /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表  */
	EnableImageIntegrityCheck *bool                                     `json:"enableImageIntegrityCheck"` /*  用于表示是否启用镜像完整性校验的标识。默认 false。注意：仅在传入的 regionID 参数所指定的资源池具备镜像完整性校验功能时生效  */
	Labels                    []*CtimageImportImageLabelsRequest        `json:"labels"`                    /*  标签列表。注意：<br />1. 列表中最多 10 个标签<br />2. 标签键不可重复<br />3. 单个标签键或值应满足长度为 1~32 个字符，不能换行，且不能以空格开头或结尾<br />4. 若传入的 regionID 参数所指定的资源池是非多可用区资源池，则此参数仅在此资源池具备从镜像文件创建私有镜像时设置标签的功能时生效  */
	ProjectID                 string                                    `json:"projectID,omitempty"`       /*  企业项目 ID。默认 0（即 default 企业项目）。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=7246&data=114&isNormal=1&vid=107" target="_blank">查询企业项目列表</a>接口来查询您可以使用的企业项目 ID  */
}

type CtimageImportImageImagePropertiesRequest struct {
	ImageName    string `json:"imageName,omitempty"`    /*  镜像名称。注意：<br />1. 长度为 2~32 个字符，只能由数字、字母、- 组成，不能以数字、- 开头，且不能以 - 结尾<br />2. 不能与您已有的私有镜像的名称重复。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您已有的私有镜像  */
	Architecture string `json:"architecture,omitempty"` /*  系统架构。取值范围（值：描述）：<br />x86_64：x86_64 架构（默认值）<br />aarch64：AArch64 架构<br />注意：<br />1. 在 imageProperties.imageType 参数值为 data_disk_image 时不会生效<br />2. 在 imageProperties.imageType 参数值为 iso_image 或在传入的 regionID 参数所指定的资源池不具备从镜像文件创建 AArch64 架构的私有镜像的功能时将始终使用默认值  */
	BootMode     string `json:"bootMode,omitempty"`     /*  x86_64 架构非数据盘镜像的启动方式。取值范围（值：描述）：<br />bios：BIOS 启动方式（默认值）<br />uefi：UEFI 启动方式<br />注意：仅在以下条件均满足时生效<br />1. imageProperties.architecture 参数值为 x86_64<br />2. imageProperties.imageType 参数值为空（系统盘镜像）<br />3. 传入的 regionID 参数所指定的资源池具备 x86_64 UEFI 启动方式功能  */
	Description  string `json:"description,omitempty"`  /*  描述信息。注意：长度为 1~128 个字符，不能以空格开头或结尾  */
	DiskSize     int32  `json:"diskSize,omitempty"`     /*  磁盘容量。单位为 GiB。默认 40。注意：<br />1. 在 imageProperties.imageType 参数值为空（系统盘镜像）时，最小 40，最大为您所能创建的 1 份私有系统盘镜像的最大大小<br />2. 在 imageProperties.imageType 参数值为 data_disk_image 时，最小 10，最大为您所能创建的 1 份私有数据盘镜像的最大大小<br />3. 在 imageProperties.imageType 参数值为 iso_image 时，最小 40，最大 1024<br />4. 若传入的值小于最小值，则自动调整为最小值。若传入的值大于最大值，则自动调整为最大值<br />温馨提示：若传入的值小于传入的 imageFileSource 参数所指定的镜像文件的虚拟大小，则所创建的镜像的磁盘容量将是此虚拟大小  */
	ImageType    string `json:"imageType,omitempty"`    /*  镜像类型。取值范围（值：描述）：<br />（空，即 null 或空字符串）：系统盘镜像（默认值）<br />data_disk_image：数据盘镜像<br />iso_image：ISO 镜像<br />注意：仅在传入的 regionID 参数所指定的资源池具备从镜像文件创建私有系统盘镜像的功能时可使用默认值。数据盘镜像和 ISO 镜像同理  */
	MaximumRAM   int32  `json:"maximumRAM,omitempty"`   /*  最大内存。单位为 GiB。取值范围：0（默认值，即不限制）/1/2/4/8/16/32/64/128/256/512。注意：<br />1. 若同时限制最大和最小内存，则最大内存不能小于最小内存<br />2. 仅在传入的 regionID 参数所指定的资源池具备设置最大/最小内存的功能且 imageType 参数值为空（系统盘镜像）时生效  */
	MinimumRAM   int32  `json:"minimumRAM,omitempty"`   /*  最小内存。单位为 GiB。取值范围：0（默认值，即不限制）/1/2/4/8/16/32/64/128/256/512。注意：<br />1. 若同时限制最大和最小内存，则最大内存不能小于最小内存<br />2. 仅在传入的 regionID 参数所指定的资源池具备设置最大/最小内存的功能且 imageType 参数值为空（系统盘镜像）时生效  */
	OsDistro     string `json:"osDistro,omitempty"`     /*  操作系统发行版。详见枚举值表<br />注意：在 imageType 参数值为 data_disk_image 时不会生效，但您仍需传入此参数或使用默认值来使镜像条目按您的预期区分 Linux 或 Windows  */
	OsVersion    string `json:"osVersion,omitempty"`    /*  操作系统版本。默认与 osDistro 参数值一致。注意：<br />1. 长度为 1~32 个字符，不能换行，且不能以空格开头或结尾<br />2. 在 imageType 参数值为 data_disk_image 时不会生效  */
}

type CtimageImportImageLabelsRequest struct {
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签键  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签值  */
}

type CtimageImportImageResponse struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功<br />900：失败  */
	Error       string                               `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                               `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                               `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageImportImageReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageImportImageReturnObjResponse struct {
	Images []*CtimageImportImageReturnObjImagesResponse `json:"images"` /*  镜像列表。应包含 1 个镜像  */
}

type CtimageImportImageReturnObjImagesResponse struct {
	AppVersion                string `json:"appVersion,omitempty"`                /*  应用版本  */
	Architecture              string `json:"architecture,omitempty"`              /*  系统架构。取值范围与对应的请求参数相同  */
	AzName                    string `json:"azName,omitempty"`                    /*  在多可用区资源池下物理机镜像的可用区名称  */
	BootMode                  string `json:"bootMode,omitempty"`                  /*  x86_64 架构非数据盘镜像的启动方式。取值范围（值：描述）：<br />bios：BIOS 启动方式<br />uefi：UEFI 启动方式  */
	ChargeableImage           *bool  `json:"chargeableImage"`                     /*  用于表示是否是收费镜像的标识  */
	ContainerFormat           string `json:"containerFormat,omitempty"`           /*  容器格式  */
	CreatedTime               int32  `json:"createdTime,omitempty"`               /*  创建时间戳  */
	CreatedTimeStr            string `json:"createdTimeStr,omitempty"`            /*  创建时间  */
	CwaiType                  string `json:"cwaiType,omitempty"`                  /*  云骁智算云主机节点类型。取值范围（值：描述）：<br />control：控制面云主机节点<br />node：GPU 云主机节点<br /><br />注意：镜像可适用于多节点类型，多个云骁智算云主机节点类型之间以英文逗号（,）分隔，如 control,node  */
	Description               string `json:"description,omitempty"`               /*  描述信息  */
	DestinationAccountID      string `json:"destinationAccountID,omitempty"`      /*  共享镜像接受者的账号 ID  */
	DestinationUser           string `json:"destinationUser,omitempty"`           /*  共享镜像接受者  */
	DiskFormat                string `json:"diskFormat,omitempty"`                /*  磁盘格式。取值范围（值：描述）：<br />qcow2：QCOW2 格式，<br />raw：RAW 格式，<br />vhd：VHD 格式，<br />vmdk：VMDK 格式  */
	DiskID                    string `json:"diskID,omitempty"`                    /*  私有镜像来源的系统盘/数据盘 ID  */
	DiskSize                  int32  `json:"diskSize,omitempty"`                  /*  磁盘容量。单位为 GiB  */
	EnableImageIntegrityCheck *bool  `json:"enableImageIntegrityCheck"`           /*  用于表示是否启用镜像完整性校验的标识  */
	FullECSDiskSize           int32  `json:"fullECSDiskSize,omitempty"`           /*  云主机整机磁盘容量。单位为 GiB  */
	GpuImageCategory          string `json:"gpuImageCategory,omitempty"`          /*  GPU 镜像种类。取值范围（值：描述）：<br />pass_through：GPU 直通镜像<br />vgpu：vGPU 镜像  */
	HasAcceptedSharedImages   *bool  `json:"hasAcceptedSharedImages"`             /*  用于表示私有镜像的共享列表中是否有镜像状态为 accepted 的共享镜像的标识  */
	ImageClass                string `json:"imageClass,omitempty"`                /*  镜像类别。取值范围（值：描述）：<br />BMS：物理机，<br />ECS：云主机  */
	ImageDisplayName          string `json:"imageDisplayName,omitempty"`          /*  镜像展示名称  */
	ImageID                   string `json:"imageID,omitempty"`                   /*  镜像 ID  */
	ImageIntegrityCheckStatus string `json:"imageIntegrityCheckStatus,omitempty"` /*  镜像完整性校验状态。详见枚举值表
	 */
	ImageName               string `json:"imageName,omitempty"`        /*  镜像名称  */
	ImageScene              string `json:"imageScene,omitempty"`       /*  镜像场景。取值范围（值：描述）：<br />dev：开发工具<br />ecommerce：电商<br />gaming：游戏<br />website：网站<br /><br />注意：镜像可适用于多场景，多个镜像场景之间以英文逗号（,）分隔，如 ecommerce,website  */
	ImageShareCount         int32  `json:"imageShareCount,omitempty"`  /*  私有镜像的共享数量  */
	ImageSize               int64  `json:"imageSize,omitempty"`        /*  镜像大小。单位为 byte  */
	ImageSource             string `json:"imageSource,omitempty"`      /*  私有镜像来源。取值范围（值：描述）：<br />cloud_server：云主机<br />full_ecs：云主机整机<br />image_file：镜像文件<br />metal_server：物理机<br />snapshot：云主机快照  */
	ImageStatus             string `json:"imageStatus,omitempty"`      /*  镜像状态。详见枚举值表  */
	ImageSubcategory        string `json:"imageSubcategory,omitempty"` /*  镜像子种类。取值范围（值：描述）：<br />app：云主机应用镜像<br />thin_app：轻量型云主机应用镜像<br /><br />注意：镜像可适用于多子种类，多个镜像子种类之间以英文逗号（,）分隔，如 app,thin_app  */
	ImageType               string `json:"imageType,omitempty"`        /*  镜像类型。取值范围（值：描述）：<br />（空，即 null）：系统盘镜像<br />data_disk_image：数据盘镜像<br />full_ecs_image：整机镜像<br />iso_image：ISO 镜像  */
	ImageVisibility         string `json:"imageVisibility,omitempty"`  /*  镜像可见类型。详见枚举值表  */
	MaximumRAM              int32  `json:"maximumRAM,omitempty"`       /*  最大内存。单位为 GiB  */
	MinimumRAM              int32  `json:"minimumRAM,omitempty"`       /*  最小内存。单位为 GiB  */
	OsDistro                string `json:"osDistro,omitempty"`         /*  操作系统发行版  */
	OsType                  string `json:"osType,omitempty"`           /*  操作系统类型。取值范围（值：描述）：<br />linux：Linux 系操作系统，<br />windows：Windows 系操作系统  */
	OsVersion               string `json:"osVersion,omitempty"`        /*  操作系统版本  */
	ProjectID               string `json:"projectID,omitempty"`        /*  企业项目 ID  */
	SourceAccountID         string `json:"sourceAccountID,omitempty"`  /*  共享镜像提供者的账号 ID  */
	SourceServerID          string `json:"sourceServerID,omitempty"`   /*  私有镜像来源的云主机/云主机快照/物理机 ID  */
	SourceUser              string `json:"sourceUser,omitempty"`       /*  共享镜像提供者  */
	SupportOneClickSFSMount *bool  `json:"supportOneClickSFSMount"`    /*  用于表示是否支持一键挂载文件系统的标识  */
	SupportXSSD             *bool  `json:"supportXSSD"`                /*  用于表示是否支持 XSSD 类型盘的标识  */
	TaskID                  string `json:"taskID,omitempty"`           /*  任务 ID  */
	UpdatedTime             int32  `json:"updatedTime,omitempty"`      /*  更新时间戳  */
	UpdatedTimeStr          string `json:"updatedTimeStr,omitempty"`   /*  更新时间  */
}
