package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsNewSfsApi
/* 弹性文件开通
 */type SfsSfsNewSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsNewSfsApi(client *core.CtyunClient) *SfsSfsNewSfsApi {
	return &SfsSfsNewSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/new-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsNewSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsNewSfsRequest) (*SfsSfsNewSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsNewSfsRequest
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
	var resp SfsSfsNewSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsNewSfsRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID，例：100054c0416811e9a6690242ac110002  */
	IsEncrypt   *bool  `json:"isEncrypt"`             /*  是否加密盘，true/false，默认 false  */
	KmsUUID     string `json:"kmsUUID,omitempty"`     /*  如果是加密盘，需要提供 kms 的 uuid  */
	ProjectID   string `json:"projectID,omitempty"`   /*  资源所属企业项目 ID，默认为'0'  */
	SfsType     string `json:"sfsType,omitempty"`     /*  弹性文件类型，capacity/performance  */
	SfsProtocol string `json:"sfsProtocol,omitempty"` /*  协议类型，nfs/cifs/nfs,cifs三种类型, nfs 适用于 linux，cifs 适用于 windows；cifs 适用于 windows；nfs,cifs可以挂载至linux和windows，仅部分"regionVersion": "v3.0"的资源池适用。  */
	SfsName     string `json:"sfsName,omitempty"`     /*  弹性文件名。单账户单资源池下，命名需唯一  */
	SfsSize     int32  `json:"sfsSize,omitempty"`     /*  大小，单位 GB，最小 500GB  */
	OnDemand    *bool  `json:"onDemand"`              /*  是否按需下单。true/false，默认为 true  */
	CycleType   string `json:"cycleType,omitempty"`   /*  包周期（subscription）类型，year/month。onDemand 为 false 时，必须指定  */
	CycleCount  int32  `json:"cycleCount,omitempty"`  /*  包周期数。onDemand 为 false 时必须指定。周期最大长度不能超过 3 年  */
	AzName      string `json:"azName,omitempty"`      /*  多可用区资源池下，必须指定可用区。4.0资源池必填  */
	Vpc         string `json:"vpc,omitempty"`         /*  虚拟网 ID  */
	Subnet      string `json:"subnet,omitempty"`      /*  子网 ID  */
}

type SfsSfsNewSfsResponse struct {
	StatusCode  int32                            `json:"statusCode"`  /*  返回状态码(800为成功，900为失败/订单处理中)  */
	Message     string                           `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                           `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsNewSfsReturnObjResponse   `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                           `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                           `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
	ErrorDetail *SfsSfsNewSfsErrorDetailResponse `json:"errorDetail"` /*  错误明细。一般情况下，会对订单侧(bss)的弹性文件订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。其他订单侧(bss)的错误，以sfs.order.procFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息  */
}

type SfsSfsNewSfsReturnObjResponse struct {
	MasterOrderID        string                                    `json:"masterOrderID"`        /*  订单 ID。调用方在拿到 masterOrderID 之后，在若干错误情况下，可以使用 materOrderID 进一步确认订单状态及资源状态  */
	MasterOrderNO        string                                    `json:"masterOrderNO"`        /*  订单号  */
	MasterResourceID     string                                    `json:"masterResourceID"`     /*  主资源 ID  */
	MasterResourceStatus string                                    `json:"masterResourceStatus"` /*  主资源状态。  */
	RegionID             string                                    `json:"regionID"`             /*  资源所属资源池 ID  */
	Resources            []*SfsSfsNewSfsReturnObjResourcesResponse `json:"resources"`            /*  资源明细  */
}

type SfsSfsNewSfsErrorDetailResponse struct{}

type SfsSfsNewSfsReturnObjResourcesResponse struct {
	ResourceID       string `json:"resourceID"`       /*  单项资源的变配、续订、退订等需要该资源项的 ID。比如某个云主机资源作为主资源，对其挂载  */
	ResourceType     string `json:"resourceType"`     /*  资源类型。SFS_TURBO/SFS_TURBOC  */
	OrderID          string `json:"orderID"`          /*  无需关心  */
	StartTime        int64  `json:"startTime"`        /*  启动时刻，epoch 时戳，毫秒精度。例：1589869069561  */
	CreateTime       int64  `json:"createTime"`       /*  创建时刻，epoch 时戳，毫秒精度  */
	UpdateTime       int64  `json:"updateTime"`       /*  更新时刻，epoch 时戳，毫秒精度  */
	Status           int32  `json:"status"`           /*  资源状态，无需关心。  */
	IsMaster         *bool  `json:"isMaster"`         /*  是否是主资源项  */
	ItemValue        int32  `json:"itemValue"`        /*  无需关心  */
	SfsUID           string `json:"sfsUID"`           /*  弹性文件系统内部唯一 ID  */
	SfsStatus        string `json:"sfsStatus"`        /*  弹性文件状态。available/unusable  */
	MasterOrderID    string `json:"masterOrderID"`    /*  订单 ID  */
	SfsName          string `json:"sfsName"`          /*  弹性文件名字  */
	MasterResourceID string `json:"masterResourceID"` /*  主资源 ID  */
}
