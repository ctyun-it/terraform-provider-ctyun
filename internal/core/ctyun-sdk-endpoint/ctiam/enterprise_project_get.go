package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EnterpriseProjectGetApi 查询企业项目
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9425&data=114
type EnterpriseProjectGetApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEnterpriseProjectGetApi(client *ctyunsdk.CtyunClient) *EnterpriseProjectGetApi {
	return &EnterpriseProjectGetApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v1/project/getEnterpriseProjectById",
		},
	}
}

func (this *EnterpriseProjectGetApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EnterpriseProjectGetRequest) (*EnterpriseProjectGetResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	builder.AddParam("id", req.Id)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp enterpriseProjectGetRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	return &EnterpriseProjectGetResponse{
		Id:          resp.Id,
		ProjectName: resp.ProjectName,
		Status:      resp.Status,
		HwProjectId: resp.HwProjectId,
		Description: resp.Description,
		CreateTime:  resp.CreateTime,
	}, nil
}

type enterpriseProjectGetRealResponse struct {
	Id          string `json:"id"`
	ProjectName string `json:"projectName"`
	Status      int    `json:"status"`
	HwProjectId string `json:"hwProjectId"`
	Description string `json:"description"`
	CreateTime  int64  `json:"createTime"`
}

type EnterpriseProjectGetResponse struct {
	Id          string
	ProjectName string
	Status      int
	HwProjectId string
	Description string
	CreateTime  int64
}

type EnterpriseProjectGetRequest struct {
	Id string
}
