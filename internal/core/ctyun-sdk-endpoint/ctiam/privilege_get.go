package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// PrivilegeGetApi 根据授权id查询授权信息
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=13979&data=114
type PrivilegeGetApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPrivilegeGetApi(client *ctyunsdk.CtyunClient) *PrivilegeGetApi {
	return &PrivilegeGetApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v1/perm/queryPrivilegeById",
		},
	}
}

func (this *PrivilegeGetApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PrivilegeGetRequest) (*PrivilegeGetResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("privilegeId", req.PrivilegeId)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp privilegeGetRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	return &PrivilegeGetResponse{
		PrivilegeId:   resp.PrivilegeId,
		RegionId:      resp.RegionId,
		Id:            resp.Id,
		AccountId:     resp.AccountId,
		PolicyId:      resp.PolicyId,
		RangeType:     resp.RangeType,
		PrincipalType: resp.PrincipalType,
	}, nil
}

type PrivilegeGetRequest struct {
	PrivilegeId string
}

type privilegeGetRealResponse struct {
	PrivilegeId   string `json:"privilegeId"`
	RegionId      string `json:"regionId"`
	Id            string `json:"id"`
	AccountId     string `json:"accountId"`
	PolicyId      string `json:"policyId"`
	RangeType     string `json:"rangeType"`
	PrincipalType string `json:"principalType"`
}

type PrivilegeGetResponse struct {
	PrivilegeId   string
	RegionId      string
	Id            string
	AccountId     string
	PolicyId      string
	RangeType     string
	PrincipalType string
}
