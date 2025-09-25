package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// PolicyDeleteApi 删除自定义策略
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9169&data=114
type PolicyDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPolicyDeleteApi(client *ctyunsdk.CtyunClient) *PolicyDeleteApi {
	return &PolicyDeleteApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/policy/deletePolicy",
		},
	}
}

func (this *PolicyDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PolicyDeleteRequest) (*PolicyDeleteResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	realRequest := policyDeleteRealRequest{
		PolicyId: req.PolicyId,
	}
	_, err := builder.WriteJson(realRequest)
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
	return &PolicyDeleteResponse{}, nil
}

type policyDeleteRealRequest struct {
	PolicyId string `json:"policyId"`
}

type PolicyDeleteRequest struct {
	PolicyId string
}

type PolicyDeleteResponse struct {
}
