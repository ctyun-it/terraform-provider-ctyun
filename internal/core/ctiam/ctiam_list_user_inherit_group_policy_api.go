package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamListUserInheritGroupPolicyApi
type CtiamListUserInheritGroupPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamListUserInheritGroupPolicyApi(client *core.CtyunClient) *CtiamListUserInheritGroupPolicyApi {
	return &CtiamListUserInheritGroupPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/perm/listUserInheritGroupPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamListUserInheritGroupPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamListUserInheritGroupPolicyRequest) (*CtiamListUserInheritGroupPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("userId", req.UserId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamListUserInheritGroupPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamListUserInheritGroupPolicyRequest struct {
	UserId string `json:"userId"` /*  用户ID  */
}

type CtiamListUserInheritGroupPolicyResponse struct {
	StatusCode *string                                           `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                           `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                           `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamListUserInheritGroupPolicyReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamListUserInheritGroupPolicyReturnObjResponse struct {
	Permissions []*CtiamListUserInheritGroupPolicyReturnObjPermissionsResponse `json:"permissions"` /*  权限集合  */
}

type CtiamListUserInheritGroupPolicyReturnObjPermissionsResponse struct {
	Id                *string `json:"id"`                /*  授权id  */
	RegionId          *string `json:"regionId"`          /*  资源池id  */
	GroupId           *string `json:"groupId"`           /*  用户组id  */
	AccountId         *string `json:"accountId"`         /*  账号ID  */
	PolicyId          *string `json:"policyId"`          /*  策略id  */
	RangeType         *string `json:"rangeType"`         /*  授权范围（GLOBAL_SERVICE：全局，PROJECT_SERVICE :   资源池）  */
	PolicyName        *string `json:"policyName"`        /*  策略名称  */
	PolicyDescription *string `json:"policyDescription"` /*  策略描述  */
	GroupName         *string `json:"groupName"`         /*  用户组名称  */
	GroupIntro        *string `json:"groupIntro"`        /*  用户组描述  */
	RegionName        *string `json:"regionName"`        /*  资源名称  */
}
