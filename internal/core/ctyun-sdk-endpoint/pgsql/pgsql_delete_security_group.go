package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlDeleteSecurityGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlDeleteSecurityGroupApi(client *ctyunsdk.CtyunClient) *PgsqlDeleteSecurityGroupApi {
	return &PgsqlDeleteSecurityGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/securityGroup/delete",
		},
	}
}

func (this *PgsqlDeleteSecurityGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlDeleteSecurityGroupRequest, header *PgsqlDeleteSecurityGroupRequestHeader) (DeleteResp *PgsqlDeleteSecurityGroupResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}

	if req.SecurityGroupId == "" {
		err = errors.New("missing required field: SecurityGroupId")
		return
	}
	if req.InstanceId == "" {
		err = errors.New("missing required field: InstanceName(实例名称)")
	}
	builder.AddParam("securityGroupId", req.SecurityGroupId)
	builder.AddParam("instanceId", req.InstanceId)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	DeleteResp = &PgsqlDeleteSecurityGroupResponse{}
	err = resp.Parse(DeleteResp)
	if err != nil {
		return
	}
	return DeleteResp, nil
}

type PgsqlDeleteSecurityGroupRequest struct {
	SecurityGroupId string `json:"securityGroupId"` // 原安全组ID，不能为空
	InstanceId      string `json:"instanceId" `     // 实例ID，不能为空
}

type PgsqlDeleteSecurityGroupRequestHeader struct {
	ProjectID *string `json:"projectId,omitempty"`
}

type PgsqlDeleteSecurityGroupResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}
