package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosRoleDetailApi
/* 获取角色详情。
 */type ZosRoleDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosRoleDetailApi(client *core.CtyunClient) *ZosRoleDetailApi {
	return &ZosRoleDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/role/detail",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosRoleDetailApi) Do(ctx context.Context, credential core.Credential, req *ZosRoleDetailRequest) (*ZosRoleDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("roleName", req.RoleName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosRoleDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosRoleDetailRequest struct {
	RegionID string /*  区域 ID  */
	RoleName string /*  角色名  */
}

type ZosRoleDetailResponse struct {
	StatusCode  int64                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                          `json:"message,omitempty"`     /*  状态描述  */
	ReturnObj   *ZosRoleDetailReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	Description string                          `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                          `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                          `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosRoleDetailReturnObjResponse struct {
	Status             *bool                                     `json:"status"`                       /*  角色状态  */
	Fuser_last_updated string                                    `json:"fuser_last_updated,omitempty"` /*  最近更新时间  */
	Role_arn           string                                    `json:"role_arn,omitempty"`           /*  角色arn  */
	MaxSessionDuration int64                                     `json:"maxSessionDuration,omitempty"` /*  最大会话时长  */
	Role_name          string                                    `json:"role_name,omitempty"`          /*  角色名  */
	Note               string                                    `json:"note,omitempty"`               /*  角色备注  */
	Policies           []*ZosRoleDetailReturnObjPoliciesResponse `json:"policies"`                     /*  该角色下绑定的策略列表  */
	Created_time       string                                    `json:"created_time,omitempty"`       /*  角色创建时间  */
}

type ZosRoleDetailReturnObjPoliciesResponse struct {
	Note        string `json:"note,omitempty"`        /*  策略备注  */
	Bind_date   string `json:"bind_date,omitempty"`   /*  策略绑定时间  */
	Policy_name string `json:"policy_name,omitempty"` /*  策略名  */
}
