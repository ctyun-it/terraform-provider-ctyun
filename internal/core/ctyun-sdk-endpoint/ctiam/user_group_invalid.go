package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// UserGroupInvalidApi 注销用户组
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9150&data=114
type UserGroupInvalidApi struct {
	builder ctyunsdk.CtyunRequestBuilder
	client  *ctyunsdk.CtyunClient
}

func NewUserGroupInvalidApi(client *ctyunsdk.CtyunClient) *UserGroupInvalidApi {
	return &UserGroupInvalidApi{
		builder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/userGroup/invalidGroup",
		},
		client: client,
	}
}

func (this *UserGroupInvalidApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserGroupInvalidRequest) (*UserGroupInvalidResponse, ctyunsdk.CtyunRequestError) {
	builder := this.builder.WithCredential(&credential)
	builder, err := builder.WriteJson(&userGroupInvalidRealRequest{
		GroupId: req.GroupId,
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
	return &UserGroupInvalidResponse{}, nil
}

type userGroupInvalidRealRequest struct {
	GroupId string `json:"groupId"`
}

type UserGroupInvalidRequest struct {
	GroupId string
}

type UserGroupInvalidResponse struct {
}
