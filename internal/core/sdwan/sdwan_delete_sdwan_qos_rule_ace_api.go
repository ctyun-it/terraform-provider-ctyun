package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanDeleteSdwanQosRuleAceApi
/* 删除qos规则五元组 */
type SdwanDeleteSdwanQosRuleAceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanDeleteSdwanQosRuleAceApi(client *core.CtyunClient) *SdwanDeleteSdwanQosRuleAceApi {
	return &SdwanDeleteSdwanQosRuleAceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/qos-rule-ace/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanDeleteSdwanQosRuleAceApi) Do(ctx context.Context, credential core.Credential, req *SdwanDeleteSdwanQosRuleAceRequest) (*SdwanDeleteSdwanQosRuleAceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanDeleteSdwanQosRuleAceRequest
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
	var resp SdwanDeleteSdwanQosRuleAceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanDeleteSdwanQosRuleAceRequest struct {
	QosID     string   `json:"qosID"`     /*  qos策略ID  */
	RuleID    string   `json:"ruleID"`    /*  规则id  */
	AceIDList []string `json:"aceIDList"` /*  ace id列表，值类型为string  */
}

type SdwanDeleteSdwanQosRuleAceResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
}
