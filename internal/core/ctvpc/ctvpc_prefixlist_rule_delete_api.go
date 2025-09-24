package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcPrefixlistRuleDeleteApi
/* 删除 prefixlist_rule
 */type CtvpcPrefixlistRuleDeleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcPrefixlistRuleDeleteApi(client *core.CtyunClient) *CtvpcPrefixlistRuleDeleteApi {
	return &CtvpcPrefixlistRuleDeleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/prefixlist_rule/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcPrefixlistRuleDeleteApi) Do(ctx context.Context, credential core.Credential, req *CtvpcPrefixlistRuleDeleteRequest) (*CtvpcPrefixlistRuleDeleteResponse, error) {
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
	var resp CtvpcPrefixlistRuleDeleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcPrefixlistRuleDeleteRequest struct {
	RegionID          string `json:"regionID,omitempty"`          /*  资源池ID  */
	PrefixListID      string `json:"prefixListID,omitempty"`      /*  prefixlistID  */
	PrefixListRuleIDs string `json:"prefixListRuleIDs,omitempty"` /*  prefixlistRuleIDs  */
}

type CtvpcPrefixlistRuleDeleteResponse struct {
	StatusCode  int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcPrefixlistRuleDeleteReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcPrefixlistRuleDeleteReturnObjResponse struct {
	PrefixListID     *string `json:"prefixListID,omitempty"`     /*  prefixlist id  */
	PrefixListRuleID *string `json:"prefixListRuleID,omitempty"` /*  prefixlist_rule id  */
}
