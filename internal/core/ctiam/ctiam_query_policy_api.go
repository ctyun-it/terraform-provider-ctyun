package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryPolicyApi
/* 根据账户ID查询所有策略 */
type CtiamQueryPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryPolicyApi(client *core.CtyunClient) *CtiamQueryPolicyApi {
	return &CtiamQueryPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/policy/queryPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryPolicyRequest) (*CtiamQueryPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryPolicyRequest
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
	var resp CtiamQueryPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryPolicyRequest struct {
	PageNum           int32   `json:"pageNum"`                     /*  页码  */
	PageSize          int32   `json:"pageSize"`                    /*  每页条数  */
	PolicyType        *int32  `json:"policyType,omitempty"`        /*  策略类型（1：系统策略、2：自定义策略）  */
	PolicyRange       *int32  `json:"policyRange,omitempty"`       /*  策略范围（1：资源池、2：全局）  */
	PolicyName        *string `json:"policyName,omitempty"`        /*  策略名称  */
	PolicyDescription *string `json:"policyDescription,omitempty"` /*  策略描述  */
}

type CtiamQueryPolicyResponse struct {
	StatusCode *string                            `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamQueryPolicyReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                            `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                            `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamQueryPolicyReturnObjResponse struct {
	PageNum  *int32                                   `json:"pageNum"`  /*  页码  */
	PageSize *int32                                   `json:"pageSize"` /*  每页条数  */
	Total    *int32                                   `json:"total"`    /*  总条数  */
	Pages    *int32                                   `json:"pages"`    /*  总页数  */
	List     []*CtiamQueryPolicyReturnObjListResponse `json:"list"`     /*  策略列表  */
}

type CtiamQueryPolicyReturnObjListResponse struct {
	Id                *string `json:"id"`                /*  策略id  */
	AccountId         *string `json:"accountId"`         /*  账号id  */
	PolicyName        *string `json:"policyName"`        /*  策略名称  */
	PolicyRange       *string `json:"policyRange"`       /*  策略范围  */
	PolicyType        *string `json:"policyType"`        /*  策略类型  */
	PolicyDescription *string `json:"policyDescription"` /*  策略描述  */
	PolicyContent     *string `json:"policyContent"`     /*  策略内容  */
	ProductName       *string `json:"productName"`       /*  产品名称  */
}
