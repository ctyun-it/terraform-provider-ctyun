package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EnterpriseProjectUpdateApi 修改企业项目
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9424&data=114
type EnterpriseProjectUpdateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEnterpriseProjectUpdateApi(client *ctyunsdk.CtyunClient) *EnterpriseProjectUpdateApi {
	return &EnterpriseProjectUpdateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/project/updateEnterpriseProject",
		},
	}
}

func (this *EnterpriseProjectUpdateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EnterpriseProjectUpdateRequest) (*EnterpriseProjectUpdateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	_, err := builder.WriteJson(&enterpriseProjectUpdateRealRequest{
		Id:          req.Id,
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

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &EnterpriseProjectUpdateResponse{}, nil
}

type EnterpriseProjectUpdateResponse struct {
}

type enterpriseProjectUpdateRealRequest struct {
	Id          string `json:"id"`
	ProjectName string `json:"projectName"`
	Description string `json:"description"`
}

type EnterpriseProjectUpdateRequest struct {
	Id          string
	ProjectName string
	Description string
}
