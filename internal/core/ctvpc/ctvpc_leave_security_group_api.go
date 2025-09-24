package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcLeaveSecurityGroupApi
/* 解绑安全组。
 */type CtvpcLeaveSecurityGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcLeaveSecurityGroupApi(client *core.CtyunClient) *CtvpcLeaveSecurityGroupApi {
	return &CtvpcLeaveSecurityGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/leave-security-group",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcLeaveSecurityGroupApi) Do(ctx context.Context, credential core.Credential, req *CtvpcLeaveSecurityGroupRequest) (*CtvpcLeaveSecurityGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcLeaveSecurityGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcLeaveSecurityGroupRequest struct {
	SecurityGroupID string `json:"securityGroupID,omitempty"` /*  安全组ID  */
	RegionID        string `json:"regionID,omitempty"`        /*  区域id  */
	InstanceID      string `json:"instanceID,omitempty"`      /*  实例ID  */
}

type CtvpcLeaveSecurityGroupResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
