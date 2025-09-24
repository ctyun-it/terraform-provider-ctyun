package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcQueryDelegateRoleApi
/* 查询委托角色
 */type CtvpcQueryDelegateRoleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcQueryDelegateRoleApi(client *core.CtyunClient) *CtvpcQueryDelegateRoleApi {
	return &CtvpcQueryDelegateRoleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/delegate/query-role",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcQueryDelegateRoleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcQueryDelegateRoleRequest) (*CtvpcQueryDelegateRoleResponse, error) {
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
	var resp CtvpcQueryDelegateRoleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcQueryDelegateRoleRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  区域ID  */
}

type CtvpcQueryDelegateRoleResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcQueryDelegateRoleReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcQueryDelegateRoleReturnObjResponse struct {
	Check    *bool                                              `json:"check"`    /*  角色校验是否通过  */
	RoleList []*CtvpcQueryDelegateRoleReturnObjRoleListResponse `json:"roleList"` /*  角色列表  */
}

type CtvpcQueryDelegateRoleReturnObjRoleListResponse struct {
	Name    *string `json:"name,omitempty"` /*  角色名称  */
	RawType int32   `json:"type"`           /*  角色类型  */
}
