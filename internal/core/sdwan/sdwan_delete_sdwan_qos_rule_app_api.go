package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanDeleteSdwanQosRuleAppApi
/* 删除qos规则应用组 */
type SdwanDeleteSdwanQosRuleAppApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanDeleteSdwanQosRuleAppApi(client *core.CtyunClient) *SdwanDeleteSdwanQosRuleAppApi {
	return &SdwanDeleteSdwanQosRuleAppApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/qos-rule-app/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanDeleteSdwanQosRuleAppApi) Do(ctx context.Context, credential core.Credential, req *SdwanDeleteSdwanQosRuleAppRequest) (*SdwanDeleteSdwanQosRuleAppResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanDeleteSdwanQosRuleAppRequest
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
	var resp SdwanDeleteSdwanQosRuleAppResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanDeleteSdwanQosRuleAppRequest struct {
	QosID     string   `json:"qosID"`     /*  qos策略ID  */
	RuleID    string   `json:"ruleID"`    /*  规则id  */
	AppIDList []string `json:"appIDList"` /*  app id列表，值类型为string  */
}

type SdwanDeleteSdwanQosRuleAppResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
}
