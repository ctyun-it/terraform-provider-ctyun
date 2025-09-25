package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsShowVolumeV41Api
/* 基于磁盘ID查询云硬盘详情<br /><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsShowVolumeV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsShowVolumeV41Api(client *core.CtyunClient) *CtecsShowVolumeV41Api {
	return &CtecsShowVolumeV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/volume/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsShowVolumeV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsShowVolumeV41Request) (*CtecsShowVolumeV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("diskID", req.DiskID)
	if req.RegionID != "" {
		ctReq.AddParam("regionID", req.RegionID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsShowVolumeV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsShowVolumeV41Request struct {
	DiskID   string /*  磁盘ID，您可以查看<a href="https://www.ctyun.cn/document/10027696/10027930">产品定义-云硬盘</a>来了解云硬盘 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7338&data=48">云硬盘列表查询</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7332&data=48&isNormal=1&vid=45">创建云硬盘</a>  */
	RegionID string /*  如本地语境支持保存regionID，那么建议传递。资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
}

type CtecsShowVolumeV41Response struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中或失败)  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码  */
	Error       string                               `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                               `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                               `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsShowVolumeV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsShowVolumeV41ReturnObjResponse struct {
	DiskName         string                                            `json:"diskName,omitempty"`         /*  磁盘名  */
	DiskID           string                                            `json:"diskID,omitempty"`           /*  磁盘ID  */
	DiskSize         int32                                             `json:"diskSize,omitempty"`         /*  磁盘大小，单位为GB  */
	DiskType         string                                            `json:"diskType,omitempty"`         /*  云硬盘规格类型，取值为：<br />●SATA：普通IO<br />●SAS：高IO<br />●SSD：超高IO<br />●FAST-SSD：极速型SSD<br />●XSSD-0、XSSD-1、XSSD-2：X系列云硬盘  */
	DiskMode         string                                            `json:"diskMode,omitempty"`         /*  云硬盘磁盘模式，取值为：<br />●VBD（Virtual Block Device）：虚拟块存储设备<br />●ISCSI （Internet Small Computer System Interface）：小型计算机系统接口<br />●FCSAN（Fibre Channel SAN）：光纤通道协议的SAN网络  */
	DiskStatus       string                                            `json:"diskStatus,omitempty"`       /*  云硬盘使用状态 deleting/creating/detaching，具体请参考<a href="https://www.ctyun.cn/document/10027696/10168629">云硬盘使用状态</a>  */
	CreateTime       int32                                             `json:"createTime,omitempty"`       /*  创建时刻，epoch时戳，精度毫秒  */
	UpdateTime       int32                                             `json:"updateTime,omitempty"`       /*  更新时刻，epoch时戳，精度毫秒  */
	ExpireTime       int32                                             `json:"expireTime,omitempty"`       /*  过期时刻，epoch时戳，精度毫秒  */
	IsSystemVolume   *bool                                             `json:"isSystemVolume"`             /*  是否系统盘，只有为系统盘时才返回该字段  */
	IsPackaged       *bool                                             `json:"isPackaged"`                 /*  是否是云主机成套资源  */
	InstanceName     string                                            `json:"instanceName,omitempty"`     /*  绑定的云主机名，有挂载时才返回  */
	InstanceID       string                                            `json:"instanceID,omitempty"`       /*  绑定云主机resourceUUID，有挂载时才返回  */
	InstanceStatus   string                                            `json:"instanceStatus,omitempty"`   /*  云主机状态  */
	MultiAttach      *bool                                             `json:"multiAttach"`                /*  是否共享云硬盘  */
	Attachments      []*CtecsShowVolumeV41ReturnObjAttachmentsResponse `json:"attachments"`                /*  挂载信息。如果是共享挂载云硬盘，有多项  */
	ProjectID        string                                            `json:"projectID,omitempty"`        /*  资源所属企业项目ID  */
	IsEncrypt        *bool                                             `json:"isEncrypt"`                  /*  是否加密盘  */
	KmsUUID          string                                            `json:"kmsUUID,omitempty"`          /*  加密盘密钥UUID，是加密盘时才返回  */
	OnDemand         *bool                                             `json:"onDemand"`                   /*  是否按需订购，按需时才返回该字段  */
	CycleType        string                                            `json:"cycleType,omitempty"`        /*  month/year，非按需时返回  */
	CycleCount       int32                                             `json:"cycleCount,omitempty"`       /*  包周期数，非按需时返回  */
	RegionID         string                                            `json:"regionID,omitempty"`         /*  资源池ID  */
	AzName           string                                            `json:"azName,omitempty"`           /*  多可用区下的可用区名字  */
	DiskFreeze       string                                            `json:"diskFreeze,omitempty"`       /*  云硬盘是否已冻结   */
	ProvisionedIops  int32                                             `json:"provisionedIops,omitempty"`  /*  XSSD类型盘的预配置iops，未配置返回0，其他类型盘不返回  */
	VolumeSource     string                                            `json:"volumeSource,omitempty"`     /*  云硬盘源快照ID，若不是从快照创建的则返回null  */
	SnapshotPolicyID string                                            `json:"snapshotPolicyID,omitempty"` /*  云硬盘绑定的快照策略ID，若没有绑定则返回null   */
}

type CtecsShowVolumeV41ReturnObjAttachmentsResponse struct {
	InstanceID   string `json:"instanceID,omitempty"`   /*  绑定云主机实例UUID  */
	AttachmentID string `json:"attachmentID,omitempty"` /*  挂载ID  */
	Device       string `json:"device,omitempty"`       /*  挂载设备名，比如/dev/sda  */
}
