package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// UserInvalidApi 注销用户
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9147&data=114
type UserInvalidApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewUserInvalidApi(client *ctyunsdk.CtyunClient) *UserInvalidApi {
	return &UserInvalidApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/user/invalidUser",
		},
		client: client,
	}
}

func (this *UserInvalidApi) Do(ctx context.Context, credential ctyunsdk.Credential, t *UserInvalidRequest) (*UserInvalidResponse, ctyunsdk.CtyunRequestError) {
	builder := this.CtyunRequestBuilder.WithCredential(&credential)
	_, err := builder.WriteJson(t)
	if err != nil {
		return nil, err
	}
	ctiam, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	resp := &UserInvalidResponse{}
	err = ctiam.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type UserInvalidRequest struct {
	UserId string `json:"userId"`
}

type UserInvalidResponse struct {
}
