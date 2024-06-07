package ctiam

import (
	"context"
	"encoding/json"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// PolicyGetApi 查询策略详情
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9173&data=114
type PolicyGetApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPolicyGetApi(client *ctyunsdk.CtyunClient) *PolicyGetApi {
	return &PolicyGetApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v1/policy/getPolicyById",
		},
	}
}

func (this *PolicyGetApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PolicyGetRequest) (*PolicyGetResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("policyId", req.PolicyId)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp policyGetRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}

	content := resp.PolicyContent
	var policyGetPolicyContentRequest PolicyGetPolicyContentRequest
	if content != "" {
		var contentResp policyGetPolicyContentRealResponse
		err2 := json.Unmarshal([]byte(content), &contentResp)
		if err2 != nil {
			return nil, ctyunsdk.ErrorAfterResponse(err2, response)
		}

		var requests []PolicyGetPolicyContentStatementRequest
		for _, st := range contentResp.Statement {
			requests = append(requests, PolicyGetPolicyContentStatementRequest{
				Resource: st.Resource,
				Action:   st.Action,
				Effect:   st.Effect,
			})
		}
		policyGetPolicyContentRequest = PolicyGetPolicyContentRequest{
			Version:   contentResp.Version,
			Statement: requests,
		}
	}

	return &PolicyGetResponse{
		Id:                resp.Id,
		PolicyName:        resp.PolicyName,
		PolicyType:        resp.PolicyType,
		PolicyRange:       resp.PolicyRange,
		PolicyDescription: resp.PolicyDescription,
		PolicyContent:     policyGetPolicyContentRequest,
	}, nil
}

type policyGetPolicyContentStatementRealResponse struct {
	Resource []string `json:"Resource"`
	Action   []string `json:"Action"`
	Effect   string   `json:"Effect"`
}

type policyGetPolicyContentRealResponse struct {
	Version   string                                        `json:"Version"`
	Statement []policyGetPolicyContentStatementRealResponse `json:"Statement"`
}

type policyGetRealResponse struct {
	Id                string `json:"id"`
	PolicyName        string `json:"policyName"`
	PolicyType        int    `json:"policyType"`
	PolicyRange       int    `json:"policyRange"`
	PolicyDescription string `json:"policyDescription"`
	PolicyContent     string `json:"policyContent"`
}

type PolicyGetRequest struct {
	PolicyId string
}

type PolicyGetPolicyContentStatementRequest struct {
	Resource []string `json:"Resource"`
	Action   []string `json:"Action"`
	Effect   string   `json:"Effect"`
}

type PolicyGetPolicyContentRequest struct {
	Version   string                                   `json:"Version"`
	Statement []PolicyGetPolicyContentStatementRequest `json:"Statement"`
}

type PolicyGetResponse struct {
	Id                string
	PolicyName        string
	PolicyType        int
	PolicyRange       int
	PolicyDescription string
	PolicyContent     PolicyGetPolicyContentRequest
}
