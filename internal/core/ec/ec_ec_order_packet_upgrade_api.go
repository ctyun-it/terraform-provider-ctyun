package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcOrderPacketUpgradeApi
/* 支持云间高速带宽包变配，目前只支持升配，即增加带宽 */
type EcEcOrderPacketUpgradeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcOrderPacketUpgradeApi(client *core.CtyunClient) *EcEcOrderPacketUpgradeApi {
	return &EcEcOrderPacketUpgradeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/packet/upgrade",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcOrderPacketUpgradeApi) Do(ctx context.Context, credential core.Credential, req *EcEcOrderPacketUpgradeRequest) (*EcEcOrderPacketUpgradeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcOrderPacketUpgradeRequest
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
	var resp EcEcOrderPacketUpgradeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcOrderPacketUpgradeRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	EcID        string  `json:"ecID"`                  /*  云间高速ID  */
	RegionID    string  `json:"regionID"`              /*  资源池ID, 例:100054c0416811e9a6690242ac110002  */
	Bandwidth   int32   `json:"bandwidth"`             /*  带宽，单位MB  */
	ResourceID  string  `json:"resourceID"`            /*  云间高速带宽包资源ID  */
}

type EcEcOrderPacketUpgradeResponse struct {
	StatusCode  *int32                                   `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *EcEcOrderPacketUpgradeReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
}

type EcEcOrderPacketUpgradeReturnObjResponse struct {
	MasterOrderID *string `json:"masterOrderID"` /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO *string `json:"masterOrderNO"` /*  订单号  */
	RegionID      *string `json:"regionID"`      /*  资源池ID  */
}
