package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUpdateSecurityGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUpdateSecurityGroupApi(client *ctyunsdk.CtyunClient) *TeledbUpdateSecurityGroupApi {
	return &TeledbUpdateSecurityGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/securityGroup/change",
		},
	}
}

func (this *TeledbUpdateSecurityGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateSecurityGroupRequest, header *TeledbUpdateSecurityGroupRequestHeader) (updateResp *TeledbUpdateSecurityGroupResponse, err error) {
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
	if req.NewSecurityGroupId == "" {
		err = errors.New("missing required field: NewSecurityGroupId")
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	updateResp = &TeledbUpdateSecurityGroupResponse{}
	err = resp.Parse(updateResp)
	if err != nil {
		return
	}
	return updateResp, nil
}

type TeledbUpdateSecurityGroupRequest struct {
	SecurityGroupId    string `json:"securityGroupId"`    // 原安全组ID，不能为空
	InstanceId         string `json:"instanceId" `        // 实例ID，不能为空
	NewSecurityGroupId string `json:"newSecurityGroupId"` // 新安全组ID，不能为空
}

type TeledbUpdateSecurityGroupRequestHeader struct {
	ProjectID *string `json:"projectId,omitempty"`
}

type TeledbUpdateSecurityGroupResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}
