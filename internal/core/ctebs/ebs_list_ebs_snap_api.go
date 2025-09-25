package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsListEbsSnapApi
/* 无。
 */type EbsListEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsListEbsSnapApi(client *core.CtyunClient) *EbsListEbsSnapApi {
	return &EbsListEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs_snapshot/list-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsListEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsListEbsSnapRequest) (*EbsListEbsSnapResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.DiskID != nil {
		ctReq.AddParam("diskID", *req.DiskID)
	}
	if req.SnapshotID != nil {
		ctReq.AddParam("snapshotID", *req.SnapshotID)
	}
	if req.SnapshotName != nil {
		ctReq.AddParam("snapshotName", *req.SnapshotName)
	}
	if req.SnapshotStatus != nil {
		ctReq.AddParam("snapshotStatus", *req.SnapshotStatus)
	}
	if req.SnapshotType != nil {
		ctReq.AddParam("snapshotType", *req.SnapshotType)
	}
	if req.VolumeAttr != nil {
		ctReq.AddParam("volumeAttr", *req.VolumeAttr)
	}
	if req.RetentionPolicy != nil {
		ctReq.AddParam("retentionPolicy", *req.RetentionPolicy)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsListEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsListEbsSnapRequest struct {
	RegionID        string  /*  资源池ID，请根据查询资源池列表接口返回值进行传参，获取“regionId”参数。  */
	DiskID          *string /*  云硬盘ID。  */
	SnapshotID      *string /*  云硬盘快照ID。  */
	SnapshotName    *string /*  云硬盘快照名称。支持输入名称的一部分做模糊匹配。  */
	SnapshotStatus  *string /*  云硬盘快照状态。取值为：<br>●available：可用。<br>●freezing：冻结。<br>●creating：创建中。<br>●deleting：删除中。<br>●rollbacking：回滚中。<br>●cloning：从快照创建云硬盘中。<br>●error：错误。  */
	SnapshotType    *string /*  云硬盘快照创建类型。取值为：<br>●manu：手动。<br>●timer：自动。  */
	VolumeAttr      *string /*  云硬盘属性。取值为：<br>●data：数据盘。<br>●system：系统盘。  */
	RetentionPolicy *string /*  云硬盘快照保留策略。取值为：<br>●forever：永久保留。<br>●custom：自定义保留天数。  */
}

type EbsListEbsSnapResponse struct {
	StatusCode  int32                            `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）。  */
	Message     string                           `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                           `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsListEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                          `json:"details,omitempty"`     /*  可忽略。  */
	ErrorCode   *string                          `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                          `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsListEbsSnapReturnObjResponse struct {
	SnapshotSize      int32                                          `json:"snapshotSize"`      /*  云硬盘快照使用量，单位为B。  */
	SnapshotTotal     int32                                          `json:"snapshotTotal"`     /*  云硬盘快照总个数。  */
	SnapshotList      []*EbsListEbsSnapReturnObjSnapshotListResponse `json:"snapshotList"`      /*  云硬盘快照列表。  */
	SnapshotManuTotal int32                                          `json:"snapshotManuTotal"` /*  云硬盘手动快照总个数，只有传diskID时返回此字段。  */
}

type EbsListEbsSnapReturnObjSnapshotListResponse struct {
	AvailabilityZone *string `json:"availabilityZone,omitempty"` /*  可用区名称。  */
	SnapshotName     *string `json:"snapshotName,omitempty"`     /*  云硬盘快照名称。  */
	SnapshotStatus   *string `json:"snapshotStatus,omitempty"`   /*  云硬盘快照状态，<br />creating/deleting/rollbacking<br />/cloning/available/error，<br />分别对应创建中/<br />删除中/回滚中<br />/从快照创建云硬盘中/<br />可用/错误  */
	CreateTime       *string `json:"createTime,omitempty"`       /*  创建时间。  */
	AzName           *string `json:"azName,omitempty"`           /*  可用区名称。  */
	DeleteTime       *string `json:"deleteTime,omitempty"`       /*  删除时间。  */
	Description      *string `json:"description,omitempty"`      /*  描述信息。  */
	ExpireTime       *string `json:"expireTime,omitempty"`       /*  过期时间。  */
	Freezing         *bool   `json:"freezing"`                   /*  是否被冻结。  */
	IsEncryted       *bool   `json:"isEncryted"`                 /*  是否是加密盘。  */
	IsMaz            *bool   `json:"isMaz"`                      /*  是否跨AZ。  */
	RegionID         *string `json:"regionID,omitempty"`         /*  资源池ID。  */
	IsTalkOrder      *bool   `json:"isTalkOrder"`                /*  是否是按需计费资源。  */
	RetentionPolicy  *string `json:"retentionPolicy,omitempty"`  /*  快照保留策略，取值为：
	●custom：自定义保留天数。
	●forever：永久保留。  */
	RetentionTime int64   `json:"retentionTime,omitempty"` /*  快照保留时间。  */
	SnapshotID    *string `json:"snapshotID,omitempty"`    /*  快照ID。  */
	SnapshotType  *string `json:"snapshotType,omitempty"`  /*  快照类型，取值为：
	●manu：手动。
	●timer：自动。  */
	UpdateTime *string `json:"updateTime,omitempty"` /*  更新时间。  */
	VolumeAttr *string `json:"volumeAttr,omitempty"` /*  云硬盘属性，取值为：
	●data：数据盘。
	●system：系统盘。  */
	VolumeName   *string `json:"volumeName,omitempty"`   /*  云硬盘名称。  */
	VolumeSize   int64   `json:"volumeSize"`             /*  云硬盘大小。  */
	VolumeSource *string `json:"volumeSource,omitempty"` /*  云硬盘来源，如果为空，<br />则是普通云硬盘，<br />如果不为空，<br />则是由快照创建而来，<br />显示来源快照ID  */
	VolumeStatus *string `json:"volumeStatus,omitempty"` /*  云硬盘的状态，请参考<a href="https://www.ctyun.cn/document/10027696/10168629">云硬盘使用状态</a>  */
	DiskID       *string `json:"diskID,omitempty"`       /*  云硬盘ID。  */
	VolumeType   *string `json:"volumeType,omitempty"`   /*  云硬盘类型，取值为：
	●SATA：普通IO。
	●SAS：高IO。
	●SSD：超高IO。
	●FAST-SSD：极速型SSD。
	●XSSD-0、XSSD-1、XSSD-2：X系列云硬盘。  */
}
