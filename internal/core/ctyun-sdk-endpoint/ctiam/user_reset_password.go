package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// UserResetPasswordApi 修改子用户密码
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=13939&data=114
type UserResetPasswordApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewUserResetPasswordApi(client *ctyunsdk.CtyunClient) *UserResetPasswordApi {
	return &UserResetPasswordApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/user/restPassword",
		},
		client: client,
	}
}

func (this UserResetPasswordApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserResetPasswordRequest) (*UserResetPasswordResponse, ctyunsdk.CtyunRequestError) {
	builder := this.CtyunRequestBuilder.WithCredential(&credential)
	builder, err := builder.WriteJson(&userResetPasswordRealRequest{
		UserId:      req.UserId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	send, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	err = send.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &UserResetPasswordResponse{}, nil
}

type userResetPasswordRealRequest struct {
	UserId      string `json:"userId"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UserResetPasswordResponse struct {
}

type UserResetPasswordRequest struct {
	UserId      string
	OldPassword string
	NewPassword string
}
