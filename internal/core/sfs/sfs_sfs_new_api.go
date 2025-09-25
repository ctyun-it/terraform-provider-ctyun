package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsNewApi
/* 创建文件系统
 */type SfsSfsNewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsNewApi(client *core.CtyunClient) *SfsSfsNewApi {
	return &SfsSfsNewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/new",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsNewApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsNewRequest) (*SfsSfsNewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsNewRequest
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
	var resp SfsSfsNewResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsNewRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性；要求单个云平台账户内唯一；参考订单幂等性说明  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID  */
	IsEncrypt   *bool  `json:"isEncrypt"`             /*  是否加密盘，默认false  */
	KmsUUID     string `json:"kmsUUID,omitempty"`     /*  如果是加密盘，需要提供kms的uuid  */
	ProjectID   string `json:"projectID,omitempty"`   /*  企业项目ID，默认为“0”，即default（默认企业项目）  */
	SfsType     string `json:"sfsType,omitempty"`     /*  存储类型，capacity(标准型)，performance（性能型）  */
	SfsProtocol string `json:"sfsProtocol,omitempty"` /*  协议类型，nfs/cifs。nfs适用于Linux，cifs适用于Windows  */
	Name        string `json:"name,omitempty"`        /*  文件系统名称；单账户单资源池下，命名需唯一  */
	SfsSize     int32  `json:"sfsSize,omitempty"`     /*  大小，单位GB，最小500GB  */
	OnDemand    *bool  `json:"onDemand"`              /*  是否按需下单；默认为true  */
	CycleType   string `json:"cycleType,omitempty"`   /*  包周期类型，year/month；onDemand为false时，必须指定  */
	CycleCount  int32  `json:"cycleCount,omitempty"`  /*  包周期数。onDemand为false时必须指定；周期最大长度不能超过3年  */
	AzName      string `json:"azName,omitempty"`      /*  多可用区资源池下，必须指定可用区  */
	Vpc         string `json:"vpc,omitempty"`         /*  虚拟私有云ID，获取：<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8658&data=94&isNormal=1&vid=88" target="_blank">查询用户VPC列表</a>
	注：在多可用区类型资源池下，vpcID通常为“vpc-”开头，非多可用区类型资源池vpcID为uuid格式  */
	Subnet string `json:"subnet,omitempty"` /*  子网ID，获取：<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8658&data=94&isNormal=1&vid=88" target="_blank">查询用户VPC列表</a>   */
}

type SfsSfsNewResponse struct {
	StateCode int32                       `json:"stateCode"` /*  返回状态码(800为成功，900为失败)  */
	Message   string                      `json:"message"`   /*  失败时的错误描述，一般为英文描述  */
	MsgDesc   string                      `json:"msgDesc"`   /*  失败时的错误描述，一般为中文描述  */
	ReturnObj *SfsSfsNewReturnObjResponse `json:"returnObj"` /*  returnObj  */
	MsgCode   string                      `json:"msgCode"`   /*  业务细分码，为product.module.code三段式码；参考结果码  */
}

type SfsSfsNewReturnObjResponse struct {
	MasterOrderID        string                                 `json:"masterOrderID"`        /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态；参考结果码  */
	MasterOrderNO        string                                 `json:"masterOrderNO"`        /*  订单号  */
	MasterResourceID     string                                 `json:"masterResourceID"`     /*  主资源ID  */
	MasterResourceStatus string                                 `json:"masterResourceStatus"` /*  主资源状态。参考主资源状态  */
	RegionID             string                                 `json:"regionID"`             /*  资源所属资源池ID  */
	Resources            []*SfsSfsNewReturnObjResourcesResponse `json:"resources"`            /*  资源明细列表，参考资源明细对象  */
}

type SfsSfsNewReturnObjResourcesResponse struct {
	ResourceID       string                                             `json:"resourceID"`       /*  单项资源的变配、续订、退订等需要该资源项的ID  */
	ResourceType     string                                             `json:"resourceType"`     /*  资源类型，SFS_TURBO（按需），SFS_TURBOC（包周期）  */
	OrderID          string                                             `json:"orderID"`          /*  无需关心  */
	StartTime        int64                                              `json:"startTime"`        /*  启动时刻，epoch 时戳，毫秒精度  */
	CreateTime       int64                                              `json:"createTime"`       /*  创建时刻，epoch 时戳，毫秒精度  */
	UpdateTime       int64                                              `json:"updateTime"`       /*  更新时刻，epoch 时戳，毫秒精度  */
	Status           int32                                              `json:"status"`           /*  资源状态  */
	IsMaster         *bool                                              `json:"isMaster"`         /*  是否是主资源项  */
	ItemValue        int32                                              `json:"itemValue"`        /*  无需关心  */
	UID              string                                             `json:"UID"`              /*  弹性文件系统内部唯一 ID  */
	SfsStatus        int32                                              `json:"sfsStatus"`        /*  弹性文件状态序号  */
	MasterOrderID    string                                             `json:"masterOrderID"`    /*  订单 ID  */
	Name             string                                             `json:"name"`             /*  文件系统名称  */
	MasterResourceID string                                             `json:"masterResourceID"` /*  主资源 ID  */
	ResourceConfig   *SfsSfsNewReturnObjResourcesResourceConfigResponse `json:"resourceConfig"`   /*  资源开通相关信息，无需关心  */
}

type SfsSfsNewReturnObjResourcesResourceConfigResponse struct{}
