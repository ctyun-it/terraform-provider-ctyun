package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsResetPasswordApi 更新云主机密码
// https://www.ctyun.cn/document/10026730/10106390
type EcsResetPasswordApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsResetPasswordApi(client *ctyunsdk.CtyunClient) *EcsResetPasswordApi {
	return &EcsResetPasswordApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/reset-password",
		},
	}
}

func (this *EcsResetPasswordApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsResetPasswordRequest) (*EcsResetPasswordResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsResetPasswordRealRequest{
		RegionID:    req.RegionId,
		InstanceID:  req.InstanceId,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse EcsResetPasswordRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &EcsResetPasswordResponse{
		InstanceId: realResponse.InstanceID,
	}, nil
}

type ecsResetPasswordRealRequest struct {
	RegionID    string `json:"regionID"`
	InstanceID  string `json:"instanceID"`
	NewPassword string `json:"newPassword"`
}

type EcsResetPasswordRealResponse struct {
	InstanceID string `json:"instanceID"`
}

type EcsResetPasswordRequest struct {
	RegionId    string
	InstanceId  string
	NewPassword string
}

type EcsResetPasswordResponse struct {
	InstanceId string
}
