package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// UserGroupGetApi 根据用户组ID查询用户组信息
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=8142&data=114
type UserGroupGetApi struct {
	builder ctyunsdk.CtyunRequestBuilder
	client  *ctyunsdk.CtyunClient
}

func NewUserGroupGetApi(client *ctyunsdk.CtyunClient) *UserGroupGetApi {
	return &UserGroupGetApi{
		builder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v1/userGroup/getGroupByGroupId",
		},
		client: client,
	}
}

func (this *UserGroupGetApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserGroupGetRequest) (*UserGroupGetResponse, ctyunsdk.CtyunRequestError) {
	builder := this.builder.WithCredential(&credential)
	builder.AddParam("groupId", req.GroupId)
	send, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	resp := &userGroupGetRealResponse{}
	err = send.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return &UserGroupGetResponse{
		Id:         resp.Id,
		GroupName:  resp.GroupName,
		GroupIntro: resp.GroupIntro,
		AccountId:  resp.AccountId,
	}, nil
}

type userGroupGetRealResponse struct {
	Id         string `json:"id"`
	GroupName  string `json:"groupName"`
	GroupIntro string `json:"groupIntro"`
	AccountId  string `json:"accountId"`
}

type UserGroupGetRequest struct {
	GroupId string
}

type UserGroupGetResponse struct {
	Id         string
	GroupName  string
	GroupIntro string
	AccountId  string
}
