package mysql

import (
	"context"
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

type TeledbUpdateSecurityGroupRequest struct{}

type TeledbUpdateSecurityGroupRequestHeader struct{}

type TeledbUpdateSecurityGroupResponse struct{}

func (this *TeledbUpdateSecurityGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateSecurityGroupRequest, header *TeledbUpdateSecurityGroupRequestHeader) (updateResp *TeledbUpdateSecurityGroupResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
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
