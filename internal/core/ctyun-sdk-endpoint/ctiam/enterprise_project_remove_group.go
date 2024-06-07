package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EnterpriseProjectRemoveGroupApi 移除企业项目关联用户组
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9430&data=114
type EnterpriseProjectRemoveGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEnterpriseProjectRemoveGroupApi(client *ctyunsdk.CtyunClient) *EnterpriseProjectRemoveGroupApi {
	return &EnterpriseProjectRemoveGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/project/removeEpFromGroup",
		},
	}
}

func (this *EnterpriseProjectRemoveGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EnterpriseProjectRemoveGroupRequest) (*EnterpriseProjectRemoveGroupResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	_, err := builder.WriteJson(&enterpriseProjectRemoveGroupRealRequest{
		ProjectId: req.ProjectId,
		GroupIds:  req.GroupIds,
	})
	if err != nil {
		return nil, err
	}

	_, err = this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	return &EnterpriseProjectRemoveGroupResponse{}, nil
}

type EnterpriseProjectRemoveGroupResponse struct {
}

type enterpriseProjectRemoveGroupRealRequest struct {
	ProjectId string   `json:"projectId"`
	GroupIds  []string `json:"groupIds"`
}

type EnterpriseProjectRemoveGroupRequest struct {
	ProjectId string
	GroupIds  []string
}
