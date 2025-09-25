package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcBatchJoinSecurityGroupApi
/* 批量绑定安全组。
 */type CtvpcBatchJoinSecurityGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcBatchJoinSecurityGroupApi(client *core.CtyunClient) *CtvpcBatchJoinSecurityGroupApi {
	return &CtvpcBatchJoinSecurityGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/batch-join-security-group",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcBatchJoinSecurityGroupApi) Do(ctx context.Context, credential core.Credential, req *CtvpcBatchJoinSecurityGroupRequest) (*CtvpcBatchJoinSecurityGroupResponse, error) {
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
	var resp CtvpcBatchJoinSecurityGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcBatchJoinSecurityGroupRequest struct {
	RegionID           string   `json:"regionID,omitempty"`           /*  区域id  */
	SecurityGroupIDs   []string `json:"securityGroupIDs"`             /*  安全组 ID 数组，最多同时支持 10 个  */
	InstanceID         string   `json:"instanceID,omitempty"`         /*  实例ID。  */
	NetworkInterfaceID *string  `json:"networkInterfaceID,omitempty"` /*  弹性网卡ID。  */
	Action             string   `json:"action,omitempty"`             /*  系统规定参数  */
}

type CtvpcBatchJoinSecurityGroupResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
