package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbsQueryEbsListApi
/* 查询某可用区全部云硬盘详情。
 */type EbsQueryEbsListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsQueryEbsListApi(client *core.CtyunClient) *EbsQueryEbsListApi {
	return &EbsQueryEbsListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs/list-ebs",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsQueryEbsListApi) Do(ctx context.Context, credential core.Credential, req *EbsQueryEbsListRequest) (*EbsQueryEbsListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.DiskStatus != nil {
		ctReq.AddParam("diskStatus", *req.DiskStatus)
	}
	if req.AzName != nil {
		ctReq.AddParam("azName", *req.AzName)
	}
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	if req.DiskType != nil {
		ctReq.AddParam("diskType", *req.DiskType)
	}
	if req.DiskMode != nil {
		ctReq.AddParam("diskMode", *req.DiskMode)
	}
	if req.MultiAttach != nil {
		ctReq.AddParam("multiAttach", *req.MultiAttach)
	}
	if req.IsSystemVolume != nil {
		ctReq.AddParam("isSystemVolume", *req.IsSystemVolume)
	}
	if req.IsEncrypt != nil {
		ctReq.AddParam("isEncrypt", *req.IsEncrypt)
	}
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsQueryEbsListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsQueryEbsListRequest struct {
	RegionID       string  /*  资源池ID。  */
	PageNo         int32   /*  默认为1。  */
	PageSize       int32   /*  默认为10，最大值为300。  */
	DiskStatus     *string /*  云硬盘状态。取值为：<br>●attached：已挂载。<br>●unattached：未挂载。<br>●detaching：卸载中。<br>●creating：创建中。<br>●expired：已过期。<br>●freezing：已冻结。  */
	AzName         *string /*  可用区。  */
	ProjectID      *string /*  企业项目。  */
	DiskType       *string /*  云硬盘类型。取值为：<br>●SATA：普通IO。<br>●SAS：高IO。<br>●SSD：超高IO。<br>●FAST-SSD：极速型SSD、<br>●XSSD-0、XSSD-1、XSSD-2：X系列云硬盘。  */
	DiskMode       *string /*  云硬盘模式。取值为：<br>●VBD：虚拟块存储设备。<br>●ISCSI：小型计算机系统接口。<br>●FCSAN：光纤通道协议的SAN网络。  */
	MultiAttach    *string /*  是否共享盘。取值为：<br>●true：共享盘。<br>●false：非共享盘。  */
	IsSystemVolume *string /*  是否为系统盘。取值为：<br>●true：系统盘。<br>●false：数据盘。  */
	IsEncrypt      *string /*  是否加密盘。取值为：<br>●true：加密盘。<br>●false：非加密盘。  */
	QueryContent   *string /*  模糊匹配盘的信息，包括盘ID、盘名称、挂载主机ID、挂载主机名称等信息。  */
}

type EbsQueryEbsListResponse struct {
	StatusCode  int32                             `json:"statusCode"`            /*  返回状态码(800为成功，900为失败)。  */
	Message     *string                           `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                           `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsQueryEbsListReturnObjResponse `json:"returnObj"`             /*  返回数据体。  */
	ErrorCode   *string                           `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，请参考错误码。  */
	Error       *string                           `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码，请参考错误码。  */
}

type EbsQueryEbsListReturnObjResponse struct {
	DiskList     []*EbsQueryEbsListReturnObjDiskListResponse `json:"diskList"`     /*  返回数据集合。  */
	DiskTotal    int32                                       `json:"diskTotal"`    /*  总数。  */
	CurrentCount int32                                       `json:"currentCount"` /*  当前页记录数目。  */
	TotalCount   int32                                       `json:"totalCount"`   /*  总记录数。  */
	TotalPage    int32                                       `json:"totalPage"`    /*  总页数。  */
}

type EbsQueryEbsListReturnObjDiskListResponse struct {
	DiskName *string `json:"diskName,omitempty"` /*  云硬盘名称。  */
	DiskID   *string `json:"diskID,omitempty"`   /*  云硬盘ID。  */
	DiskSize int32   `json:"diskSize"`           /*  磁盘大小，单位为GB。  */
	DiskType *string `json:"diskType,omitempty"` /*  磁盘规格类型，取值为：
	●SATA：普通IO。
	●SAS：高IO。
	●SSD：超高IO。
	●FAST-SSD：极速型SSD。
	●XSSD-0、XSSD-1、XSSD-2：X系列云硬盘。  */
	DiskMode *string `json:"diskMode,omitempty"` /*  云硬盘磁盘模式，取值为：
	●VBD（Virtual Block Device）：虚拟块存储设备。
	●ISCSI （Internet Small Computer System Interface）：小型计算机系统接口。
	●FCSAN（Fibre Channel SAN）：光纤通道协议的SAN网络。  */
	DiskStatus      *string                                                `json:"diskStatus,omitempty"`     /*  参考 <a href='https://www.ctyun.cn/document/10027696/10168629'>云硬盘使用状态</a>  */
	CreateTime      int64                                                  `json:"createTime"`               /*  创建时刻，epoch时戳，精度为毫秒。  */
	UpdateTime      int64                                                  `json:"updateTime"`               /*  更新时刻，epoch时戳，精度为毫秒。  */
	ExpireTime      int64                                                  `json:"expireTime"`               /*  过期时刻，epoch时戳，精度为毫秒。  */
	IsSystemVolume  *bool                                                  `json:"isSystemVolume"`           /*  只有为系统盘时才返回该字段。  */
	IsPackaged      *bool                                                  `json:"isPackaged"`               /*  是否随云主机一起订购。  */
	InstanceName    *string                                                `json:"instanceName,omitempty"`   /*  绑定的云主机名称，有挂载时才返回。  */
	InstanceID      *string                                                `json:"instanceID,omitempty"`     /*  绑定的云主机ID，有挂载时才返回。  */
	InstanceStatus  *string                                                `json:"instanceStatus,omitempty"` /*  云主机状态，参考<a href='https://www.ctyun.cn/document/10027696/10168629'>云主机使用状态</a>  */
	MultiAttach     *bool                                                  `json:"multiAttach"`              /*  是否是共享云硬盘。  */
	Attachments     []*EbsQueryEbsListReturnObjDiskListAttachmentsResponse `json:"attachments"`              /*  挂载信息。如果是共享挂载云硬盘，则返回多项；无挂载时不返回该字段。  */
	ProjectID       *string                                                `json:"projectID,omitempty"`      /*  云硬盘所属的企业项目ID。  */
	IsEncrypt       *bool                                                  `json:"isEncrypt"`                /*  是否是加密盘。  */
	KmsUUID         *string                                                `json:"kmsUUID,omitempty"`        /*  加密盘密钥UUID，是加密盘时才返回。  */
	RegionID        *string                                                `json:"regionID,omitempty"`       /*  资源池ID。  */
	AzName          *string                                                `json:"azName,omitempty"`         /*  多可用区下的可用区名字，非多可用区不返回该字段。  */
	DiskFreeze      *bool                                                  `json:"diskFreeze,omitempty"`     /*  云硬盘是否已冻结。  */
	ProvisionedIops int32                                                  `json:"provisionedIops"`          /*  XSSD类型盘的预配置iops，未配置返回0，其他类型云硬盘不返回。  */
}

type EbsQueryEbsListReturnObjDiskListAttachmentsResponse struct {
	InstanceID   *string `json:"instanceID,omitempty"`   /*  绑定的云主机ID。  */
	AttachmentID *string `json:"attachmentID,omitempty"` /*  挂载ID。  */
	Device       *string `json:"device,omitempty"`       /*  挂载的设备名，例如/dev/sda。  */
}
