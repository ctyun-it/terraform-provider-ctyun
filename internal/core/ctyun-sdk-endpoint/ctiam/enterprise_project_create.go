package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EnterpriseProjectCreateApi 创建企业项目
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9423&data=114
type EnterpriseProjectCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEnterpriseProjectCreateApi(client *ctyunsdk.CtyunClient) *EnterpriseProjectCreateApi {
	return &EnterpriseProjectCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/project/createEnterpriseProject",
		},
	}
}

func (this *EnterpriseProjectCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EnterpriseProjectCreateRequest) (*EnterpriseProjectCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	_, err := builder.WriteJson(&enterpriseProjectCreateRealRequest{
		ProjectName: req.ProjectName,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp enterpriseProjectCreateRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	return &EnterpriseProjectCreateResponse{
		ProjectId: resp.ProjectId,
	}, nil
}

type enterpriseProjectCreateRealResponse struct {
	ProjectId string `json:"projectId"`
}

type EnterpriseProjectCreateResponse struct {
	ProjectId string
}

type enterpriseProjectCreateRealRequest struct {
	ProjectName string `json:"projectName"`
	Description string `json:"description"`
}

type EnterpriseProjectCreateRequest struct {
	ProjectName string
	Description string
}
