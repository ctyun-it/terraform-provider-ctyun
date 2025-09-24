package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type MongodbUpdateSecurityGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbUpdateSecurityGroupApi(client *ctyunsdk.CtyunClient) *MongodbUpdateSecurityGroupApi {
	return &MongodbUpdateSecurityGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/securityGroup/change",
		},
	}
}

func (this *MongodbUpdateSecurityGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbUpdateSecurityGroupRequest, header *MongodbUpdateSecurityGroupRequestHeader) (updateResp *MongodbUpdateSecurityGroupResponse, err error) {
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

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return
	}
	updateResp = &MongodbUpdateSecurityGroupResponse{}
	err = resp.Parse(updateResp)
	if err != nil {
		return
	}
	return updateResp, nil
}

type MongodbUpdateSecurityGroupRequest struct {
	SecurityGroupId    string `json:"securityGroupId"`    // 原安全组ID，不能为空
	InstanceId         string `json:"instanceId" `        // 实例ID，不能为空
	NewSecurityGroupId string `json:"newSecurityGroupId"` // 新安全组ID，不能为空
}

type MongodbUpdateSecurityGroupRequestHeader struct {
	ProjectID *string `json:"projectId,omitempty"`
}

type MongodbUpdateSecurityGroupResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}
