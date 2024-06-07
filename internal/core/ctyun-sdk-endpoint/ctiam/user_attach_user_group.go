package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// UserAttachUserGroupApi 用户批量加入用户组
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=13944&data=114
type UserAttachUserGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewUserAttachUserGroupApi(client *ctyunsdk.CtyunClient) *UserAttachUserGroupApi {
	return &UserAttachUserGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/user/attachGroupToUser",
		},
	}
}

func (this *UserAttachUserGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserAttachUserGroupRequest) (*UserAttachUserGroupResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&userAttachUserGroupRealRequest{
		UserId:   req.UserId,
		GroupIds: req.GroupIds,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &UserAttachUserGroupResponse{}, nil
}

type userAttachUserGroupRealRequest struct {
	UserId   string   `json:"userId"`
	GroupIds []string `json:"groupIds"`
}

type UserAttachUserGroupRequest struct {
	UserId   string
	GroupIds []string
}

type UserAttachUserGroupResponse struct {
}
