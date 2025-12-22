package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamListGroupPolicyApi
/* 通过用户组ID查询权限 */
type CtiamListGroupPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamListGroupPolicyApi(client *core.CtyunClient) *CtiamListGroupPolicyApi {
	return &CtiamListGroupPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/perm/listGroupPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamListGroupPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamListGroupPolicyRequest) (*CtiamListGroupPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("groupId", req.GroupId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamListGroupPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamListGroupPolicyRequest struct {
	GroupId string `json:"groupId"` /*  用户组id  */
}

type CtiamListGroupPolicyResponse struct {
	StatusCode *string                                `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamListGroupPolicyReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                                `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamListGroupPolicyReturnObjResponse struct {
	Permissions []*CtiamListGroupPolicyReturnObjPermissionsResponse `json:"permissions"` /*  权限列表  */
}

type CtiamListGroupPolicyReturnObjPermissionsResponse struct {
	Id                *string `json:"id"`                /*  权限id  */
	RegionId          *string `json:"regionId"`          /*  资源池id  */
	GroupId           *string `json:"groupId"`           /*  用户组id  */
	AccountId         *string `json:"accountId"`         /*  账号id  */
	PolicyId          *string `json:"policyId"`          /*  策略id  */
	RangeType         *string `json:"rangeType"`         /*  授权范围  */
	PolicyName        *string `json:"policyName"`        /*  策略名称  */
	PolicyDescription *string `json:"policyDescription"` /*  策略描述  */
	GroupName         *string `json:"groupName"`         /*  用户组名称  */
	GroupIntro        *string `json:"groupIntro"`        /*  描述  */
	RegionName        *string `json:"regionName"`        /*  资源池名称  */
}
