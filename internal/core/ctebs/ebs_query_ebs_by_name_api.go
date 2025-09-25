package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsQueryEbsByNameApi
/* 基于资源池ID和云硬盘名称查询云硬盘详情。
 */type EbsQueryEbsByNameApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsQueryEbsByNameApi(client *core.CtyunClient) *EbsQueryEbsByNameApi {
	return &EbsQueryEbsByNameApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs/info-by-name-ebs",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsQueryEbsByNameApi) Do(ctx context.Context, credential core.Credential, req *EbsQueryEbsByNameRequest) (*EbsQueryEbsByNameResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("diskName", req.DiskName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsQueryEbsByNameResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsQueryEbsByNameRequest struct {
	RegionID string /*  资源池ID。  */
	DiskName string /*  云硬盘名称。  */
}

type EbsQueryEbsByNameResponse struct {
	StatusCode  int32                               `json:"statusCode"`            /*  返回状态码(800为成功，900为失败)。  */
	Message     *string                             `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                             `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsQueryEbsByNameReturnObjResponse `json:"returnObj"`             /*  云硬盘信息查询结果。  */
	ErrorCode   *string                             `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，请参考错误码。  */
	Error       *string                             `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码，请参考错误码。  */
}

type EbsQueryEbsByNameReturnObjResponse struct {
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
	DiskStatus       *string                                          `json:"diskStatus,omitempty"`       /*  参考 <a href='https://www.ctyun.cn/document/10027696/10168629'>云硬盘使用状态</a>  */
	CreateTime       int64                                            `json:"createTime"`                 /*  创建时刻，epoch时戳，精度毫秒。  */
	UpdateTime       int64                                            `json:"updateTime"`                 /*  更新时刻，epoch时戳，精度毫秒。  */
	ExpireTime       int64                                            `json:"expireTime"`                 /*  过期时刻，epoch时戳，精度毫秒。  */
	IsSystemVolume   *bool                                            `json:"isSystemVolume"`             /*  只有为系统盘时才返回该字段。  */
	IsPackaged       *bool                                            `json:"isPackaged"`                 /*  是否随云主机一起订购。  */
	InstanceName     *string                                          `json:"instanceName,omitempty"`     /*  绑定的云主机名称，有挂载时才返回  */
	InstanceID       *string                                          `json:"instanceID,omitempty"`       /*  绑定的云主机ID，有挂载时才返回。  */
	InstanceStatus   *string                                          `json:"instanceStatus,omitempty"`   /*  云主机状态，参考<a href='https://www.ctyun.cn/document/10027696/10168629'>云主机使用状态</a>  */
	MultiAttach      *bool                                            `json:"multiAttach"`                /*  是否是共享云硬盘。  */
	Attachments      []*EbsQueryEbsByNameReturnObjAttachmentsResponse `json:"attachments"`                /*  挂载信息。如果是共享挂载云硬盘，则返回多项；无挂载时不返回该字段。  */
	ProjectID        *string                                          `json:"projectID,omitempty"`        /*  云硬盘所属的企业项目ID。  */
	IsEncrypt        *bool                                            `json:"isEncrypt"`                  /*  是否是加密盘。  */
	KmsUUID          *string                                          `json:"kmsUUID,omitempty"`          /*  加密盘密钥UUID，是加密盘时才返回。  */
	OnDemand         *bool                                            `json:"onDemand"`                   /*  是否按需订购，按需时才返回该字段。  */
	CycleType        *string                                          `json:"cycleType,omitempty"`        /*  包周期类型，year：年，month：月。非按需时才返回。  */
	CycleCount       int32                                            `json:"cycleCount"`                 /*  包周期数，非按需时才返回。  */
	RegionID         *string                                          `json:"regionID,omitempty"`         /*  资源池ID。  */
	AzName           *string                                          `json:"azName,omitempty"`           /*  多可用区下的可用区名称。  */
	DiskFreeze       *bool                                            `json:"diskFreeze,omitempty"`       /*  云硬盘是否已冻结。  */
	ProvisionedIops  int32                                            `json:"provisionedIops"`            /*  XSSD类型盘的预配置iops，未配置返回0，其他类型盘不返回。  */
	VolumeSource     *string                                          `json:"volumeSource,omitempty"`     /*  云硬盘源快照ID，若不是从快照创建的则返回null。  */
	SnapshotPolicyID *string                                          `json:"snapshotPolicyID,omitempty"` /*  云硬盘绑定的快照策略ID，若没有绑定则返回null。  */
}

type EbsQueryEbsByNameReturnObjAttachmentsResponse struct {
	InstanceID   *string `json:"instanceID,omitempty"`   /*  绑定的云主机ID。  */
	AttachmentID *string `json:"attachmentID,omitempty"` /*  挂载ID。  */
	Device       *string `json:"device,omitempty"`       /*  挂载的设备名，比如/dev/sda。  */
}
