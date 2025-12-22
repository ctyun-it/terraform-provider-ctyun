package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamAttachGroupPolicyApi
/* 创建用户组权限 */
type CtiamAttachGroupPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamAttachGroupPolicyApi(client *core.CtyunClient) *CtiamAttachGroupPolicyApi {
	return &CtiamAttachGroupPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/perm/attachGroupPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamAttachGroupPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamAttachGroupPolicyRequest) (*CtiamAttachGroupPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamAttachGroupPolicyRequest
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
	var resp CtiamAttachGroupPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamAttachGroupPolicyRequest struct {
	UserGroupId string    `json:"userGroupId"`         /*  用户组id  */
	RangeType   string    `json:"rangeType"`           /*  授权范围（GLOBAL_SERVICE：全局 、PROJECT_SERVICE：资源池级别）  */
	PolicyIds   []string  `json:"policyIds"`           /*  策略ID列表  */
	RegionIds   []*string `json:"regionIds,omitempty"` /*  资源池ID列表（rangeType为PROJECT_SERVICE时必填）  */
}

type CtiamAttachGroupPolicyResponse struct {
	StatusCode *string                                  `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamAttachGroupPolicyReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                                  `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                  `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamAttachGroupPolicyReturnObjResponse struct {
	AccountId        *string                                                    `json:"accountId"`        /*  账号id  */
	UserGroupId      *string                                                    `json:"userGroupId"`      /*  用户组id  */
	RangeType        *string                                                    `json:"rangeType"`        /*  授权类型  */
	PolicyIds        []*string                                                  `json:"policyIds"`        /*  策略id列表  */
	RegionIds        *string                                                    `json:"regionIds"`        /*  资源池列表  */
	PolicyList       *string                                                    `json:"policyList"`       /*  策略列表  */
	PrivilegeMessage []*CtiamAttachGroupPolicyReturnObjPrivilegeMessageResponse `json:"privilegeMessage"` /*  授权信息  */
}

type CtiamAttachGroupPolicyReturnObjPrivilegeMessageResponse struct {
	PrivilegeId *string `json:"privilegeId"` /*  授权id  */
	PolicyId    *string `json:"policyId"`    /*  策略  */
	RegionId    *string `json:"regionId"`    /*  资源池id  */
}
