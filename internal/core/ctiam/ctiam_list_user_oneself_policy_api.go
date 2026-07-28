package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamListUserOneselfPolicyApi
type CtiamListUserOneselfPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamListUserOneselfPolicyApi(client *core.CtyunClient) *CtiamListUserOneselfPolicyApi {
	return &CtiamListUserOneselfPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/perm/listUserOneselfPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamListUserOneselfPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamListUserOneselfPolicyRequest) (*CtiamListUserOneselfPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("userId", req.UserId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamListUserOneselfPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamListUserOneselfPolicyRequest struct {
	UserId string `json:"userId"` /*  用户id  */
}

type CtiamListUserOneselfPolicyResponse struct {
	StatusCode *string                                      `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                      `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                      `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamListUserOneselfPolicyReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamListUserOneselfPolicyReturnObjResponse struct {
	Permissions []*CtiamListUserOneselfPolicyReturnObjPermissionsResponse `json:"permissions"` /*  权限集合  */
}

type CtiamListUserOneselfPolicyReturnObjPermissionsResponse struct {
	Id         *string `json:"id"`         /*  权限id  */
	RegionId   *string `json:"regionId"`   /*  资源池id  */
	AccountId  *string `json:"accountId"`  /*  账号id  */
	PolicyId   *string `json:"policyId"`   /*  策略id  */
	RangeType  *string `json:"rangeType"`  /*  授权范围：GLOBAL_SERVICE：全局，PROJECT_SERVICE：资源池  */
	PolicyName *string `json:"policyName"` /*  策略名称  */
	RegionName *string `json:"regionName"` /*  资源池名称  */
}
