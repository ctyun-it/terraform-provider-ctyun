package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamAttachPolicyToUserApi
type CtiamAttachPolicyToUserApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamAttachPolicyToUserApi(client *core.CtyunClient) *CtiamAttachPolicyToUserApi {
	return &CtiamAttachPolicyToUserApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/perm/attachPolicyToUser",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamAttachPolicyToUserApi) Do(ctx context.Context, credential core.Credential, req *CtiamAttachPolicyToUserRequest) (*CtiamAttachPolicyToUserResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamAttachPolicyToUserRequest
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
	var resp CtiamAttachPolicyToUserResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamAttachPolicyToUserRequest struct {
	UserId    string    `json:"userId"`              /*  用户ID  */
	RangeType string    `json:"rangeType"`           /*  授权范围（GLOBAL_SERVICE：全局 、PROJECT_SERVICE：资源池级别）  */
	PolicyIds []string  `json:"policyIds"`           /*  策略ID列表  */
	RegionIds []*string `json:"regionIds,omitempty"` /*  资源池ID列表（rangeType为PROJECT_SERVICE时必填）  */
}

type CtiamAttachPolicyToUserResponse struct {
	StatusCode *string                                   `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *CtiamAttachPolicyToUserReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Message    *string                                   `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                   `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamAttachPolicyToUserReturnObjResponse struct {
	AccountId        *string                                                     `json:"accountId"`        /*  账号id  */
	UserId           *string                                                     `json:"userId"`           /*  用户ID  */
	RangeType        *string                                                     `json:"rangeType"`        /*  授权类型  */
	PolicyIds        []*string                                                   `json:"policyIds"`        /*  策略id列表  */
	RegionIds        []*string                                                   `json:"regionIds"`        /*  资源池列表  */
	PolicyList       []*CtiamAttachPolicyToUserReturnObjPolicyListResponse       `json:"policyList"`       /*  策略列表  */
	PrivilegeMessage []*CtiamAttachPolicyToUserReturnObjPrivilegeMessageResponse `json:"privilegeMessage"` /*  授权信息  */
}

type CtiamAttachPolicyToUserReturnObjPolicyListResponse struct {
	PolicyId          *string `json:"policyId"`          /*  策略id  */
	PolicyName        *string `json:"policyName"`        /*  策略名称  */
	PolicyType        *string `json:"policyType"`        /*  策略类型（CUS_154_03_0001：系统策略、CUS_154_03_0002：自定义策略）  */
	PolicyDescription *string `json:"policyDescription"` /*  策略描述  */
}

type CtiamAttachPolicyToUserReturnObjPrivilegeMessageResponse struct {
	PrivilegeId *string `json:"privilegeId"` /*  授权id  */
	PolicyId    *string `json:"policyId"`    /*  策略id  */
	RegionId    *string `json:"regionId"`    /*  资源池id  */
}
