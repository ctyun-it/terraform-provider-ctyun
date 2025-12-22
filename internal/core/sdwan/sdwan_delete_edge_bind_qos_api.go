package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanDeleteEdgeBindQosApi
/* 智能网关绑定qos */
type SdwanDeleteEdgeBindQosApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanDeleteEdgeBindQosApi(client *core.CtyunClient) *SdwanDeleteEdgeBindQosApi {
	return &SdwanDeleteEdgeBindQosApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-bind-qos/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanDeleteEdgeBindQosApi) Do(ctx context.Context, credential core.Credential, req *SdwanDeleteEdgeBindQosRequest) (*SdwanDeleteEdgeBindQosResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanDeleteEdgeBindQosRequest
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
	var resp SdwanDeleteEdgeBindQosResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanDeleteEdgeBindQosRequest struct {
	QosID      string   `json:"qosID"`      /*  qos策略ID  */
	EdgeIDList []string `json:"edgeIDList"` /*  盒子列表  ，值类型为string  */
}

type SdwanDeleteEdgeBindQosResponse struct {
	StatusCode  int32                                      `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                    `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                    `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                    `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanDeleteEdgeBindQosReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                    `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanDeleteEdgeBindQosReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作日志Id  */
}
