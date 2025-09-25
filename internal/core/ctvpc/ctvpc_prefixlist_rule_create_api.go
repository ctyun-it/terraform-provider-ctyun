package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcPrefixlistRuleCreateApi
/* 创建 prefixlist_rule
 */type CtvpcPrefixlistRuleCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcPrefixlistRuleCreateApi(client *core.CtyunClient) *CtvpcPrefixlistRuleCreateApi {
	return &CtvpcPrefixlistRuleCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/prefixlist_rule/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcPrefixlistRuleCreateApi) Do(ctx context.Context, credential core.Credential, req *CtvpcPrefixlistRuleCreateRequest) (*CtvpcPrefixlistRuleCreateResponse, error) {
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
	var resp CtvpcPrefixlistRuleCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcPrefixlistRuleCreateRequest struct {
	RegionID        string                                             `json:"regionID,omitempty"`     /*  资源池ID  */
	PrefixListID    string                                             `json:"prefixListID,omitempty"` /*  prefixlistID  */
	PrefixListRules []*CtvpcPrefixlistRuleCreatePrefixListRulesRequest `json:"prefixListRules"`        /*  接口业务数据  */
}

type CtvpcPrefixlistRuleCreatePrefixListRulesRequest struct {
	Cidr string `json:"cidr,omitempty"` /*  前缀列表条目,cidr  */
}

type CtvpcPrefixlistRuleCreateResponse struct {
	StatusCode  int32                                         `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcPrefixlistRuleCreateReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcPrefixlistRuleCreateReturnObjResponse struct {
	PrefixListID *string `json:"prefixListID,omitempty"` /*  prefixlist id  */
}
