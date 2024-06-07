package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// UserRemoveUserGroupApi 用户批量移出用户组
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=13945&data=114
type UserRemoveUserGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewUserRemoveUserGroupApi(client *ctyunsdk.CtyunClient) *UserRemoveUserGroupApi {
	return &UserRemoveUserGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/user/removeGroupFromUser",
		},
	}
}

func (this *UserRemoveUserGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserRemoveUserGroupRequest) (*UserRemoveUserGroupResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&userRemoveUserGroupRealRequest{
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
	return &UserRemoveUserGroupResponse{}, nil
}

type userRemoveUserGroupRealRequest struct {
	UserId   string   `json:"userId"`
	GroupIds []string `json:"groupIds"`
}

type UserRemoveUserGroupRequest struct {
	UserId   string
	GroupIds []string
}

type UserRemoveUserGroupResponse struct {
}
