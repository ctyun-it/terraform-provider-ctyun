package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcLIstInstanceBandwidthPacketApi
/* 查询实例带宽包 */
type EcLIstInstanceBandwidthPacketApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcLIstInstanceBandwidthPacketApi(client *core.CtyunClient) *EcLIstInstanceBandwidthPacketApi {
	return &EcLIstInstanceBandwidthPacketApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/instance-bandwidth-packet/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcLIstInstanceBandwidthPacketApi) Do(ctx context.Context, credential core.Credential, req *EcLIstInstanceBandwidthPacketRequest) (*EcLIstInstanceBandwidthPacketResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.EcID != nil && *req.EcID != "" {
		ctReq.AddParam("ecID", *req.EcID)
	}
	if req.CgwID != nil && *req.CgwID != "" {
		ctReq.AddParam("cgwID", *req.CgwID)
	}
	if req.IbpID != nil && *req.IbpID != "" {
		ctReq.AddParam("ibpID", *req.IbpID)
	}
	if req.InstanceID != nil && *req.InstanceID != "" {
		ctReq.AddParam("instanceID", *req.InstanceID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcLIstInstanceBandwidthPacketResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcLIstInstanceBandwidthPacketRequest struct {
	EcID       *string `json:"ecID,omitempty"`       /*  云间高速ID  */
	CgwID      *string `json:"cgwID,omitempty"`      /*  云网关ID。组合查询：ecID+cgwID查询云间高速云网关下的带宽包详细信息  */
	IbpID      *string `json:"ibpID,omitempty"`      /*  实例带宽包ID。组合查询：ibpID查询某个带宽包详细信息  */
	InstanceID *string `json:"instanceID,omitempty"` /*  实例ID。组合查询：ecID+cgwID+instanceID查询云间高速云网关下指定实例的带宽包详细信息  */
}

type EcLIstInstanceBandwidthPacketResponse struct {
	StatusCode  *int32                                          `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                         `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                         `json:"message"`     /*  失败时的错误描述，一般为英文描述   */
	Description *string                                         `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                         `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcLIstInstanceBandwidthPacketReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                         `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcLIstInstanceBandwidthPacketReturnObjResponse struct {
	CurrentCount *int32                                                   `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                                   `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                                   `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcLIstInstanceBandwidthPacketReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组   */
}

type EcLIstInstanceBandwidthPacketReturnObjResultsResponse struct {
	IbpID      *string `json:"ibpID"`      /*  实例带宽包ID  */
	IbpName    *string `json:"ibpName"`    /*  实例带宽包名字  */
	EcID       *string `json:"ecID"`       /*  云间高速实例ID  */
	CgwID      *string `json:"cgwID"`      /*  云网关ID  */
	InstanceID *string `json:"instanceID"` /*  实例唯一ID信息  */
	Bandwidth  *int32  `json:"bandwidth"`  /*  带宽，单位MB  */
}
