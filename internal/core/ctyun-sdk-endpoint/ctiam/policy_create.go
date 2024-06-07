package ctiam

import (
	"context"
	"encoding/json"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// PolicyCreateApi 创建自定义策略
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9167&data=114
type PolicyCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPolicyCreateApi(client *ctyunsdk.CtyunClient) *PolicyCreateApi {
	return &PolicyCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/policy/createPolicy",
		},
	}
}

func (this *PolicyCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PolicyCreateRequest) (*PolicyCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	var statments []policyCreatePolicyContentStatementRealRequest
	for _, statment := range req.PolicyContent.Statement {
		statments = append(statments, policyCreatePolicyContentStatementRealRequest{
			Resource: statment.Resource,
			Action:   statment.Action,
			Effect:   statment.Effect,
		})
	}
	request := policyCreatePolicyContentRealRequest{
		Version:   req.PolicyContent.Version,
		Statement: statments,
	}
	rawBytes, e := json.Marshal(&request)
	if e != nil {
		return nil, ctyunsdk.ErrorBeforeRequest(e)
	}

	realRequest := policyCreateRealRequest{
		PolicyName:        req.PolicyName,
		PolicyRange:       req.PolicyRange,
		PolicyDescription: req.PolicyDescription,
		PolicyContent:     string(rawBytes),
	}
	_, err := builder.WriteJson(realRequest)
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp policyCreateRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	return &PolicyCreateResponse{
		PolicyId: resp.PolicyId,
	}, nil
}

type policyCreatePolicyContentStatementRealRequest struct {
	Resource []string `json:"Resource"`
	Action   []string `json:"Action"`
	Effect   string   `json:"Effect"`
}

type policyCreatePolicyContentRealRequest struct {
	Version   string                                          `json:"Version"`
	Statement []policyCreatePolicyContentStatementRealRequest `json:"Statement"`
}

type policyCreateRealRequest struct {
	PolicyName        string `json:"policyName"`
	PolicyContent     string `json:"policyContent"`
	PolicyRange       int    `json:"policyRange"`
	PolicyDescription string `json:"policyDescription"`
}

type policyCreateRealResponse struct {
	PolicyId string `json:"policyId"`
}

type PolicyCreatePolicyContentStatementRequest struct {
	Resource []string
	Action   []string
	Effect   string
}

type PolicyCreatePolicyContentRequest struct {
	Version   string
	Statement []PolicyCreatePolicyContentStatementRequest
}

type PolicyCreateRequest struct {
	PolicyName        string
	PolicyRange       int
	PolicyDescription string
	PolicyContent     PolicyCreatePolicyContentRequest
}

type PolicyCreateResponse struct {
	PolicyId string
}
