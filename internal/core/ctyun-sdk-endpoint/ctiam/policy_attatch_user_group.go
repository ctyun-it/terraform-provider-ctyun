package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// PolicyAttachUserGroupApi 创建用户组权限
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9156&data=114
type PolicyAttachUserGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPolicyAttachUserGroupApi(client *ctyunsdk.CtyunClient) *PolicyAttachUserGroupApi {
	return &PolicyAttachUserGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/perm/attachGroupPolicy",
		},
	}
}

func (this *PolicyAttachUserGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PolicyAttachUserGroupRequest) (*PolicyAttachUserGroupResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&policyAttachUserGroupRealRequest{
		UserGroupId: req.UserGroupId,
		RangeType:   req.RangeType,
		PolicyIds:   req.PolicyIds,
		RegionIds:   req.RegionIds,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp policyAttachUserGroupRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	var policyList []PolicyAttachUserGroupPolicyListResponse
	for _, policy := range resp.PolicyList {
		policyList = append(policyList, PolicyAttachUserGroupPolicyListResponse{
			PolicyId:   policy.PolicyId,
			PolicyName: policy.PolicyName,
			PolicyType: policy.PolicyType,
		})
	}
	var privilegeList []PolicyAttachUserGroupPolicyPrivilegeMessageResponse
	for _, privilege := range resp.PrivilegeMessage {
		privilegeList = append(privilegeList, PolicyAttachUserGroupPolicyPrivilegeMessageResponse{
			PrivilegeId: privilege.PrivilegeId,
			PloyId:      privilege.PloyId,
		})
	}
	return &PolicyAttachUserGroupResponse{
		AccountId:        resp.AccountId,
		UserGroupId:      resp.UserGroupId,
		RangeType:        resp.RangeType,
		PolicyIds:        resp.PolicyIds,
		PolicyList:       policyList,
		PrivilegeMessage: privilegeList,
	}, nil
}

type policyAttachUserGroupRealRequest struct {
	UserGroupId string   `json:"userGroupId"`
	RangeType   string   `json:"rangeType"`
	PolicyIds   []string `json:"policyIds"`
	RegionIds   []string `json:"regionIds"`
}

type policyAttachUserGroupRealResponse struct {
	AccountId   string   `json:"accountId"`
	UserGroupId string   `json:"userGroupId"`
	RangeType   string   `json:"rangeType"`
	PolicyIds   []string `json:"policyIds"`
	PolicyList  []struct {
		PolicyId   string `json:"policyId"`
		PolicyName string `json:"policyName"`
		PolicyType string `json:"policyType"`
	} `json:"policyList"`
	PrivilegeMessage []struct {
		PrivilegeId string `json:"privilegeId"`
		PloyId      string `json:"ployId"`
	} `json:"privilegeMessage"`
}

type PolicyAttachUserGroupRequest struct {
	UserGroupId string
	RangeType   string
	PolicyIds   []string
	RegionIds   []string
}

type PolicyAttachUserGroupPolicyListResponse struct {
	PolicyId   string
	PolicyName string
	PolicyType string
}

type PolicyAttachUserGroupPolicyPrivilegeMessageResponse struct {
	PrivilegeId string
	PloyId      string
}

type PolicyAttachUserGroupResponse struct {
	AccountId        string
	UserGroupId      string
	RangeType        string
	PolicyIds        []string
	PolicyList       []PolicyAttachUserGroupPolicyListResponse
	PrivilegeMessage []PolicyAttachUserGroupPolicyPrivilegeMessageResponse
}
