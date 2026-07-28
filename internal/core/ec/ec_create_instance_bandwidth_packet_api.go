package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcCreateInstanceBandwidthPacketApi
/* 创建实例带宽包 */
type EcCreateInstanceBandwidthPacketApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcCreateInstanceBandwidthPacketApi(client *core.CtyunClient) *EcCreateInstanceBandwidthPacketApi {
	return &EcCreateInstanceBandwidthPacketApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/instance-bandwidth-packet/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EcCreateInstanceBandwidthPacketApi) Do(ctx context.Context, credential core.Credential, req *EcCreateInstanceBandwidthPacketRequest) (*EcCreateInstanceBandwidthPacketResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcCreateInstanceBandwidthPacketRequest
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
	var resp EcCreateInstanceBandwidthPacketResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcCreateInstanceBandwidthPacketRequest struct {
	IbpName      string `json:"ibpName"`      /*  实例带宽包名字  */
	EcID         string `json:"ecID"`         /*  云间高速实例ID  */
	CgwID        string `json:"cgwID"`        /*  云网关ID  */
	InstanceID   string `json:"instanceID"`   /*  实例唯一ID信息  */
	InstanceType int32  `json:"instanceType"` /*  2 CDA实例，其他实例暂时不支持  */
	Bandwidth    int32  `json:"bandwidth"`    /*  带宽，单位MB  */
}

type EcCreateInstanceBandwidthPacketResponse struct {
	StatusCode  *int32                                            `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                           `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                           `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                           `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                           `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcCreateInstanceBandwidthPacketReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                           `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcCreateInstanceBandwidthPacketReturnObjResponse struct {
	IbpID   *string `json:"ibpID"`   /*  实例带宽包ID  */
	Message *string `json:"message"` /*  结果  */
}
