package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsPermissionRuleDeleteApi
/* 删除已有权限组规则
 */type SfsPermissionRuleDeleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsPermissionRuleDeleteApi(client *core.CtyunClient) *SfsPermissionRuleDeleteApi {
	return &SfsPermissionRuleDeleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/permission-rule/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsPermissionRuleDeleteApi) Do(ctx context.Context, credential core.Credential, req *SfsPermissionRuleDeleteRequest) (*SfsPermissionRuleDeleteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsPermissionRuleDeleteRequest
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
	var resp SfsPermissionRuleDeleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsPermissionRuleDeleteRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池 ID  */
	PermissionRuleID string `json:"permissionRuleID,omitempty"` /*  权限组规则的fuid  */
}

type SfsPermissionRuleDeleteResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                    `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                    `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsPermissionRuleDeleteReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                    `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                                    `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsPermissionRuleDeleteReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  删除权限组规则的操作号  */
}
