package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// PolicyAttachUserApi 为用户授权
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=13946&data=114
type PolicyAttachUserApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPolicyAttachUserApi(client *ctyunsdk.CtyunClient) *PolicyAttachUserApi {
	return &PolicyAttachUserApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/perm/attachPolicyToUser",
		},
	}
}

func (this *PolicyAttachUserApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PolicyAttachUserRequest) (*PolicyAttachUserResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&policyAttachUserRealRequest{
		UserId:    req.UserId,
		RangeType: req.RangeType,
		PolicyIds: req.PolicyIds,
		RegionIds: req.RegionIds,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp policyAttachUserRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	var policyList []PolicyAttachUserPolicyListResponse
	for _, policy := range resp.PolicyList {
		policyList = append(policyList, PolicyAttachUserPolicyListResponse{
			PolicyId:   policy.PolicyId,
			PolicyName: policy.PolicyName,
			PolicyType: policy.PolicyType,
		})
	}
	var privilegeList []PolicyAttachUserPolicyPrivilegeMessageResponse
	for _, privilege := range resp.PrivilegeMessage {
		privilegeList = append(privilegeList, PolicyAttachUserPolicyPrivilegeMessageResponse{
			PrivilegeId: privilege.PrivilegeId,
			PloyId:      privilege.PloyId,
		})
	}
	return &PolicyAttachUserResponse{
		AccountId:        resp.AccountId,
		UserId:           resp.UserId,
		RangeType:        resp.RangeType,
		PolicyIds:        resp.PolicyIds,
		PolicyList:       policyList,
		PrivilegeMessage: privilegeList,
	}, nil
}

type policyAttachUserRealRequest struct {
	UserId    string   `json:"userId"`
	RangeType string   `json:"rangeType"`
	PolicyIds []string `json:"policyIds"`
	RegionIds []string `json:"regionIds"`
}

type policyAttachUserRealResponse struct {
	AccountId  string   `json:"accountId"`
	UserId     string   `json:"userId"`
	RangeType  string   `json:"rangeType"`
	PolicyIds  []string `json:"policyIds"`
	PolicyList []struct {
		PolicyId   string `json:"policyId"`
		PolicyName string `json:"policyName"`
		PolicyType string `json:"policyType"`
	} `json:"policyList"`
	PrivilegeMessage []struct {
		PrivilegeId string `json:"privilegeId"`
		PloyId      string `json:"ployId"`
	} `json:"privilegeMessage"`
}

type PolicyAttachUserRequest struct {
	UserId    string
	RangeType string
	PolicyIds []string
	RegionIds []string
}

type PolicyAttachUserPolicyListResponse struct {
	PolicyId   string
	PolicyName string
	PolicyType string
}

type PolicyAttachUserPolicyPrivilegeMessageResponse struct {
	PrivilegeId string
	PloyId      string
}

type PolicyAttachUserResponse struct {
	AccountId        string
	UserId           string
	RangeType        string
	PolicyIds        []string
	PolicyList       []PolicyAttachUserPolicyListResponse
	PrivilegeMessage []PolicyAttachUserPolicyPrivilegeMessageResponse
}
