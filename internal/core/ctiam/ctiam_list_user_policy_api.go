package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamListUserPolicyApi
/* 通过用户ID查询权限 */
type CtiamListUserPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamListUserPolicyApi(client *core.CtyunClient) *CtiamListUserPolicyApi {
	return &CtiamListUserPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/perm/listUserPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamListUserPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamListUserPolicyRequest) (*CtiamListUserPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("userId", req.UserId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamListUserPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamListUserPolicyRequest struct {
	UserId string `json:"userId"` /*  用户id  */
}

type CtiamListUserPolicyResponse struct {
	StatusCode *string                               `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamListUserPolicyReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                               `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                               `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamListUserPolicyReturnObjResponse struct {
	Permissions []*CtiamListUserPolicyReturnObjPermissionsResponse `json:"permissions"` /*  权限列表  */
}

type CtiamListUserPolicyReturnObjPermissionsResponse struct {
	Id                *string `json:"id"`                /*  权限id  */
	RegionId          *string `json:"regionId"`          /*  资源池id  */
	GroupId           *string `json:"groupId"`           /*  用户组id  */
	AccountId         *string `json:"accountId"`         /*  账号id  */
	PolicyId          *string `json:"policyId"`          /*  策略id  */
	RangeType         *string `json:"rangeType"`         /*  授权范围  */
	PolicyName        *string `json:"policyName"`        /*  策略名称  */
	PolicyDescription *string `json:"policyDescription"` /*  策略描述  */
	GroupName         *string `json:"groupName"`         /*  用户组名称  */
	GroupIntro        *string `json:"groupIntro"`        /*  用户组描述  */
}
