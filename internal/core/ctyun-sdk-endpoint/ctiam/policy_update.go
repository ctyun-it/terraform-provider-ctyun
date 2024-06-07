package ctiam

import (
	"context"
	"encoding/json"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// PolicyUpdateApi 编辑自定义策略
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9168&data=114
type PolicyUpdateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPolicyUpdateApi(client *ctyunsdk.CtyunClient) *PolicyUpdateApi {
	return &PolicyUpdateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/policy/updatePolicy",
		},
	}
}

func (this *PolicyUpdateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PolicyUpdateRequest) (*PolicyUpdateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	var statments []policyUpdatePolicyContentStatementRealRequest
	for _, statment := range req.PolicyContent.Statement {
		statments = append(statments, policyUpdatePolicyContentStatementRealRequest{
			Resource: statment.Resource,
			Action:   statment.Action,
			Effect:   statment.Effect,
		})
	}
	request := policyUpdatePolicyContentRealRequest{
		Version:   req.PolicyContent.Version,
		Statement: statments,
	}
	rawBytes, e := json.Marshal(&request)
	if e != nil {
		return nil, ctyunsdk.ErrorBeforeRequest(e)
	}

	realRequest := policyUpdateRealRequest{
		PolicyId:          req.PolicyId,
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

	var resp policyUpdateRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	return &PolicyUpdateResponse{
		PolicyId: resp.PolicyId,
	}, nil
}

type policyUpdatePolicyContentStatementRealRequest struct {
	Resource []string `json:"Resource"`
	Action   []string `json:"Action"`
	Effect   string   `json:"Effect"`
}

type policyUpdatePolicyContentRealRequest struct {
	Version   string                                          `json:"Version"`
	Statement []policyUpdatePolicyContentStatementRealRequest `json:"Statement"`
}

type policyUpdateRealRequest struct {
	PolicyId          string `json:"policyId"`
	PolicyName        string `json:"policyName"`
	PolicyContent     string `json:"policyContent"`
	PolicyRange       int    `json:"policyRange"`
	PolicyDescription string `json:"policyDescription"`
}

type policyUpdateRealResponse struct {
	PolicyId string `json:"policyId"`
}

type PolicyUpdatePolicyContentStatementRequest struct {
	Resource []string
	Action   []string
	Effect   string
}

type PolicyUpdatePolicyContentRequest struct {
	Version   string
	Statement []PolicyUpdatePolicyContentStatementRequest
}

type PolicyUpdateRequest struct {
	PolicyId          string
	PolicyName        string
	PolicyRange       int
	PolicyDescription string
	PolicyContent     PolicyUpdatePolicyContentRequest
}

type PolicyUpdateResponse struct {
	PolicyId string
}
