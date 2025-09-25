package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcChangeFilteRulePriorityApi
/* 调整过滤规则优先级
 */type CtvpcChangeFilteRulePriorityApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcChangeFilteRulePriorityApi(client *core.CtyunClient) *CtvpcChangeFilteRulePriorityApi {
	return &CtvpcChangeFilteRulePriorityApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/mirrorflow/change-filter-rule-priority",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcChangeFilteRulePriorityApi) Do(ctx context.Context, credential core.Credential, req *CtvpcChangeFilteRulePriorityRequest) (*CtvpcChangeFilteRulePriorityResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcChangeFilteRulePriorityResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcChangeFilteRulePriorityRequest struct {
	RegionID            string   `json:"regionID,omitempty"`       /*  区域ID  */
	MirrorFilterID      string   `json:"mirrorFilterID,omitempty"` /*  过滤条件 ID  */
	MirrorFilterRuleIDs []string `json:"mirrorFilterRuleIDs"`      /*  调整后的规则列表，优先级按照数据的顺序  */
}

type CtvpcChangeFilteRulePriorityResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
