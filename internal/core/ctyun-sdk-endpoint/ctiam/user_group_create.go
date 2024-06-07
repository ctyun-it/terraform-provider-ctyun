package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// UserGroupCreateApi 创建用户组
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=8136&data=114
type UserGroupCreateApi struct {
	builder ctyunsdk.CtyunRequestBuilder
	client  *ctyunsdk.CtyunClient
}

func NewUserGroupCreateApi(client *ctyunsdk.CtyunClient) *UserGroupCreateApi {
	return &UserGroupCreateApi{
		builder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/userGroup/createGroup",
		},
		client: client,
	}
}

func (this *UserGroupCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserGroupCreateRequest) (*UserGroupCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.builder.WithCredential(&credential)
	builder, err := builder.WriteJson(&userGroupCreateRealRequest{
		GroupName:  req.GroupName,
		GroupIntro: req.GroupIntro,
	})
	if err != nil {
		return nil, err
	}
	send, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	resp := &userGroupCreateRealResponse{}
	err = send.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return &UserGroupCreateResponse{
		Id: resp.Id,
	}, nil
}

type userGroupCreateRealRequest struct {
	GroupName  string `json:"groupName"`
	GroupIntro string `json:"groupIntro"`
}

type userGroupCreateRealResponse struct {
	Id string `json:"id"`
}

type UserGroupCreateRequest struct {
	GroupName  string
	GroupIntro string
}

type UserGroupCreateResponse struct {
	Id string
}
