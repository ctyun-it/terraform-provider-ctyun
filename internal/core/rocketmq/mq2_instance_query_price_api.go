package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2InstanceQueryPriceApi
/* 开通查价 */
type Mq2InstanceQueryPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2InstanceQueryPriceApi(client *core.CtyunClient) *Mq2InstanceQueryPriceApi {
	return &Mq2InstanceQueryPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/queryPrice",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2InstanceQueryPriceApi) Do(ctx context.Context, credential core.Credential, req *Mq2InstanceQueryPriceRequest) (*Mq2InstanceQueryPriceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2InstanceQueryPriceRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2InstanceQueryPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2InstanceQueryPriceRequest struct {
	RegionId  string  `json:"regionId"`           /*  资源池编码  */
	CycleType string  `json:"cycleType"`          /*  计费模式101:按需 3:按月订购  */
	CycleCnt  *string `json:"cycleCnt,omitempty"` /*  cycleType=3时，必填；订购周期，取值范围：值需大于零，不超过384个月  */
	SpecName  string  `json:"specName"`           /*  规格名  */
	NodeNum   int32   `json:"nodeNum"`            /*  节点数  */
	DiskType  string  `json:"diskType"`           /*  存储类型 SAS 、SSD  */
	DiskSize  int32   `json:"diskSize"`           /*  节点存储空间大小,单位G  */
}

type Mq2InstanceQueryPriceResponse struct {
	StatusCode *string                                 `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                 `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2InstanceQueryPriceReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                                 `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2InstanceQueryPriceReturnObjResponse struct {
	ItemId       *string `json:"itemId"`       /*  订单项ID  */
	TotalPrice   *string `json:"totalPrice"`   /*  订单项总价格  */
	FinalPrice   *string `json:"finalPrice"`   /*  订单项最终价格  */
	ResourceType *string `json:"resourceType"` /*  资源类型  */
}
