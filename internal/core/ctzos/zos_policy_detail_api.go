package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPolicyDetailApi
/* 获取策略详情。
 */type ZosPolicyDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPolicyDetailApi(client *core.CtyunClient) *ZosPolicyDetailApi {
	return &ZosPolicyDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/policy/detail",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPolicyDetailApi) Do(ctx context.Context, credential core.Credential, req *ZosPolicyDetailRequest) (*ZosPolicyDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("policyName", req.PolicyName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosPolicyDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPolicyDetailRequest struct {
	RegionID   string /*  区域 ID  */
	PolicyName string /*  策略名称  */
}

type ZosPolicyDetailResponse struct {
	StatusCode  int64                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                            `json:"message,omitempty"`     /*  状态描述  */
	ReturnObj   *ZosPolicyDetailReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	Description string                            `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                            `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosPolicyDetailReturnObjResponse struct {
	Note            string                                         `json:"note,omitempty"`            /*  策略备注  */
	Policy_document string                                         `json:"policy_document,omitempty"` /*  策略详情  */
	Policy_role     []*ZosPolicyDetailReturnObjPolicy_roleResponse `json:"policy_role"`               /*  该策略绑定的角色列表  */
	Policy_name     string                                         `json:"policy_name,omitempty"`     /*  策略名称  */
}

type ZosPolicyDetailReturnObjPolicy_roleResponse struct {
	Note      string `json:"note,omitempty"`      /*  角色备注  */
	Role_name string `json:"role_name,omitempty"` /*  角色名称  */
	Bind_date string `json:"bind_date,omitempty"` /*  该策略与该角色绑定时间  */
}
