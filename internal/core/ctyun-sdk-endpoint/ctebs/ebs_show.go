package ctebs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EbsShowApi 云硬盘信息查询（基于diskID）
// https://www.ctyun.cn/document/10027696/10110677
type EbsShowApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEbsShowApi(client *ctyunsdk.CtyunClient) *EbsShowApi {
	return &EbsShowApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/ebs/info-ebs",
		},
	}
}

func (this *EbsShowApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EbsShowRequest) (*EbsShowResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("regionID", req.RegionId)
	builder.AddParam("diskID", req.DiskId)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtebs, builder)
	if err != nil {
		return nil, err
	}

	result := &ebsShowRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}

	var attach []EbsAttachment
	for _, realAttachment := range result.Attachments {
		attach = append(attach, EbsAttachment{
			InstanceId:   realAttachment.InstanceID,
			AttachmentId: realAttachment.AttachmentID,
			Device:       realAttachment.Device,
		})
	}

	return &EbsShowResponse{
		DiskName:       result.DiskName,
		DiskID:         result.DiskID,
		DiskSize:       result.DiskSize,
		DiskType:       result.DiskType,
		DiskMode:       result.DiskMode,
		DiskStatus:     result.DiskStatus,
		CreateTime:     result.CreateTime,
		UpdateTime:     result.UpdateTime,
		ExpireTime:     result.ExpireTime,
		IsSystemVolume: result.IsSystemVolume,
		IsPackaged:     result.IsPackaged,
		InstanceName:   result.InstanceName,
		InstanceID:     result.InstanceID,
		InstanceStatus: result.InstanceStatus,
		MultiAttach:    result.MultiAttach,
		Attachments:    attach,
		ProjectID:      result.ProjectID,
		IsEncrypt:      result.IsEncrypt,
		KmsUUID:        result.KmsUUID,
		OnDemand:       result.OnDemand,
		CycleType:      &result.CycleType,
		CycleCount:     &result.CycleCount,
		RegionID:       result.RegionID,
		AzName:         result.AzName,
	}, nil
}

type ebsShowRealResponse struct {
	DiskName       string              `json:"diskName"`            // 磁盘名
	DiskID         string              `json:"diskID"`              // 磁盘ID
	DiskSize       int64               `json:"diskSize"`            // 磁盘大小（GB）
	DiskType       string              `json:"diskType"`            // 磁盘规格类型 SATA/SAS/SSD-genric/SSD/FAST-SSD
	DiskMode       string              `json:"diskMode"`            // 磁盘模式。VBD/ISCSI/FCSAN
	DiskStatus     string              `json:"diskStatus"`          // 云硬盘使用状态 deleting/creating/detaching，具体请参考云硬盘使用状态
	CreateTime     int64               `json:"createTime"`          // 创建时刻，epoch时戳，精度毫秒
	UpdateTime     int64               `json:"updateTime"`          // 更新时刻，epoch时戳，精度毫秒
	ExpireTime     int64               `json:"expireTime"`          // 过期时刻，epoch时戳，精度毫秒
	IsSystemVolume bool                `json:"isSystemVolume"`      // 是否系统盘，只有为系统盘时才返回该字段
	IsPackaged     bool                `json:"isPackaged"`          // 是否是云主机成套资源
	InstanceName   string              `json:"instanceName"`        // 绑定的云主机名，有挂载时才返回
	InstanceID     string              `json:"instanceID"`          // 绑定云主机resourceUUID，有挂载时才返回
	InstanceStatus string              `json:"instanceStatus"`      // 云主机状态starting/restarting/stopping，具体参考云主机状态，有挂载时才返回
	MultiAttach    bool                `json:"multiAttach"`         // 是否共享云硬盘
	Attachments    []ebsRealAttachment `json:"attachments"`         // 挂载信息。如果是共享挂载云硬盘，有多项		参考表attachment
	ProjectID      string              `json:"projectID,omitempty"` // 资源所属企业项目id
	IsEncrypt      bool                `json:"isEncrypt"`           // 是否加密盘
	KmsUUID        string              `json:"kmsUUID"`             // 加密盘密钥UUID，是加密盘时才返回
	OnDemand       bool                `json:"onDemand"`            // 是否按需订购，按需时才返回该字段
	CycleType      string              `json:"cycleType"`           // month/year，非按需时返回
	CycleCount     int64               `json:"cycleCount"`          // 包周期数，非按需时返回
	RegionID       string              `json:"regionID"`            // 资源池ID
	AzName         string              `json:"azName"`              // 多可用区下的可用区名字
}

type ebsRealAttachment struct {
	InstanceID   string `json:"instanceID"`   // 绑定云主机实例UUID
	AttachmentID string `json:"attachmentID"` // 挂载ID
	Device       string `json:"device"`       // 挂载设备名，比如/dev/sda
}

type EbsShowRequest struct {
	RegionId string
	DiskId   string
}

type EbsShowResponse struct {
	DiskName       string          // 磁盘名
	DiskID         string          // 磁盘ID
	DiskSize       int64           // 磁盘大小（GB）
	DiskType       string          // 磁盘规格类型 SATA/SAS/SSD-genric/SSD/FAST-SSD
	DiskMode       string          // 磁盘模式。VBD/ISCSI/FCSAN
	DiskStatus     string          // 云硬盘使用状态 deleting/creating/detaching，具体请参考云硬盘使用状态
	CreateTime     int64           // 创建时刻，epoch时戳，精度毫秒
	UpdateTime     int64           // 更新时刻，epoch时戳，精度毫秒
	ExpireTime     int64           // 过期时刻，epoch时戳，精度毫秒
	IsSystemVolume bool            // 是否系统盘，只有为系统盘时才返回该字段
	IsPackaged     bool            // 是否是云主机成套资源
	InstanceName   string          // 绑定的云主机名，有挂载时才返回
	InstanceID     string          // 绑定云主机resourceUUID，有挂载时才返回
	InstanceStatus string          // 云主机状态starting/restarting/stopping，具体参考云主机状态，有挂载时才返回
	MultiAttach    bool            // 是否共享云硬盘
	Attachments    []EbsAttachment // 挂载信息。如果是共享挂载云硬盘，有多项		参考表attachment
	ProjectID      string          // 资源所属企业项目id
	IsEncrypt      bool            // 是否加密盘
	KmsUUID        string          // 加密盘密钥UUID，是加密盘时才返回
	OnDemand       bool            // 是否按需订购，按需时才返回该字段
	CycleType      *string         // month/year，非按需时返回
	CycleCount     *int64          // 包周期数，非按需时返回
	RegionID       string          // 资源池ID
	AzName         string          // 多可用区下的可用区名字
}

type EbsAttachment struct {
	InstanceId   string // 绑定云主机实例UUID
	AttachmentId string // 挂载ID
	Device       string // 挂载设备名，比如/dev/sda
}
