package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
	"strconv"
)

// EnterpriseProjectStatusUpdateApi 启用停用企业项目
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9426&data=114
type EnterpriseProjectStatusUpdateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEnterpriseProjectStatusUpdateApi(client *ctyunsdk.CtyunClient) *EnterpriseProjectStatusUpdateApi {
	return &EnterpriseProjectStatusUpdateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/project/sasEnterpriseProject",
		},
	}
}

func (this *EnterpriseProjectStatusUpdateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EnterpriseProjectStatusUpdateRequest) (*EnterpriseProjectStatusUpdateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	// ctiam的bug相关
	_, _ = builder.WriteJson(struct {
	}{})

	builder.AddParam("projectId", req.ProjectId).
		AddParam("status", strconv.Itoa(req.Status))

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &EnterpriseProjectStatusUpdateResponse{}, nil
}

type EnterpriseProjectStatusUpdateResponse struct {
}

type EnterpriseProjectStatusUpdateRequest struct {
	ProjectId string
	Status    int
}
