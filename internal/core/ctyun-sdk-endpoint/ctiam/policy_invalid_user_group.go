package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// PolicyInvalidUserGroupApi 用户组取消权限
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9157&data=114
type PolicyInvalidUserGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPolicyInvalidUserGroupApi(client *ctyunsdk.CtyunClient) *PolicyInvalidUserGroupApi {
	return &PolicyInvalidUserGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/perm/invalidGroupPolicy",
		},
	}
}

func (this *PolicyInvalidUserGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PolicyInvalidUserGroupRequest) (*PolicyInvalidUserGroupResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&policyInvalidUserGroupRealRequest{
		PrivilegeId: req.PrivilegeId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &PolicyInvalidUserGroupResponse{}, nil
}

type policyInvalidUserGroupRealRequest struct {
	PrivilegeId string `json:"privilegeId"`
}

type PolicyInvalidUserGroupRequest struct {
	PrivilegeId string
}

type PolicyInvalidUserGroupResponse struct {
}
