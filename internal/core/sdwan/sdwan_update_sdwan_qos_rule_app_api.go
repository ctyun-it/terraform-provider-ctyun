package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanUpdateSdwanQosRuleAppApi
/* 修改qos规则应用组 */
type SdwanUpdateSdwanQosRuleAppApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanUpdateSdwanQosRuleAppApi(client *core.CtyunClient) *SdwanUpdateSdwanQosRuleAppApi {
	return &SdwanUpdateSdwanQosRuleAppApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/qos-rule-app/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanUpdateSdwanQosRuleAppApi) Do(ctx context.Context, credential core.Credential, req *SdwanUpdateSdwanQosRuleAppRequest) (*SdwanUpdateSdwanQosRuleAppResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanUpdateSdwanQosRuleAppRequest
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
	var resp SdwanUpdateSdwanQosRuleAppResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanUpdateSdwanQosRuleAppRequest struct {
	QosID   string                                      `json:"qosID"`             /*  qos策略ID  */
	RuleID  string                                      `json:"ruleID"`            /*  规则id  */
	ItemID  string                                      `json:"itemID"`            /*  qos策略应用实例ID  */
	AppList []*SdwanUpdateSdwanQosRuleAppAppListRequest `json:"appList,omitempty"` /*  app列表  */
}

type SdwanUpdateSdwanQosRuleAppAppListRequest struct {
	AppName string `json:"appName"` /*  app名称  */
	AppID   string `json:"appID"`   /*  app id  */
}

type SdwanUpdateSdwanQosRuleAppResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码   */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
	Error       *string   `json:"error"`       /*  业务细分码，为product.module.code三段式码   */
}
