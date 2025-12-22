package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetQosRuleCanConfigCountApi
/* 查询qos可配置规则数量 */
type SdwanGetQosRuleCanConfigCountApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetQosRuleCanConfigCountApi(client *core.CtyunClient) *SdwanGetQosRuleCanConfigCountApi {
	return &SdwanGetQosRuleCanConfigCountApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/qos-rule/count",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetQosRuleCanConfigCountApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetQosRuleCanConfigCountRequest) (*SdwanGetQosRuleCanConfigCountResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("qosID", req.QosID)
	if req.QosName != nil && *req.QosName != "" {
		ctReq.AddParam("qosName", *req.QosName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetQosRuleCanConfigCountResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetQosRuleCanConfigCountRequest struct {
	QosID   string  `json:"qosID"`             /*  qos策略ID  */
	QosName *string `json:"qosName,omitempty"` /*  qos策略名称  */
}

type SdwanGetQosRuleCanConfigCountResponse struct {
	StatusCode  int32                                           `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                         `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                         `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                         `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetQosRuleCanConfigCountReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                         `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetQosRuleCanConfigCountReturnObjResponse struct {
	CanConfigCount int32 `json:"canConfigCount"` /*  可配置数量  */
}
