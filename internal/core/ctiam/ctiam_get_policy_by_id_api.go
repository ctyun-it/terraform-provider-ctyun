package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetPolicyByIdApi
/* 查询策略详情 */
type CtiamGetPolicyByIdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetPolicyByIdApi(client *core.CtyunClient) *CtiamGetPolicyByIdApi {
	return &CtiamGetPolicyByIdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/policy/getPolicyById",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetPolicyByIdApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetPolicyByIdRequest) (*CtiamGetPolicyByIdResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("policyId", req.PolicyId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamGetPolicyByIdResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetPolicyByIdRequest struct {
	PolicyId string `json:"policyId"` /*  策略id  */
}

type CtiamGetPolicyByIdResponse struct {
	StatusCode *string                              `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamGetPolicyByIdReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                              `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                              `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamGetPolicyByIdReturnObjResponse struct {
	Id                *string `json:"id"`                /*  策略id  */
	PolicyName        *string `json:"policyName"`        /*  策略名称  */
	PolicyType        int32   `json:"policyType"`        /*  策略类型  */
	PolicyRange       int32   `json:"policyRange"`       /*  策略范围  */
	PolicyDescription *string `json:"policyDescription"` /*  策略描述  */
	PolicyContent     *string `json:"policyContent"`     /*  策略内容  */
	CreateTime        int64   `json:"createTime"`        /*  创建时间  */
}
