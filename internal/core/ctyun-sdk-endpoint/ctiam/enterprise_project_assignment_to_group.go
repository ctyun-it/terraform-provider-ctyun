package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EnterpriseProjectAssignmentToGroupApi 用户组与企业项目关联
type EnterpriseProjectAssignmentToGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEpAssignmentToGroupApi(client *ctyunsdk.CtyunClient) *EnterpriseProjectAssignmentToGroupApi {
	return &EnterpriseProjectAssignmentToGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/project/assignmentEpToGroup",
		},
	}
}

func (this *EnterpriseProjectAssignmentToGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EnterpriseProjectAssignmentToGroupRequest) (*EnterpriseProjectAssignmentToGroupResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)

	_, err := builder.WriteJson(&enterpriseProjectAssignmentToGroupRealRequest{
		ProjectId: req.ProjectId,
		GroupIds:  req.GroupIds,
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
	return &EnterpriseProjectAssignmentToGroupResponse{}, nil
}

type EnterpriseProjectAssignmentToGroupResponse struct {
}

type enterpriseProjectAssignmentToGroupRealRequest struct {
	ProjectId string   `json:"projectId"`
	GroupIds  []string `json:"groupIds"`
}

type EnterpriseProjectAssignmentToGroupRequest struct {
	ProjectId string
	GroupIds  []string
}
