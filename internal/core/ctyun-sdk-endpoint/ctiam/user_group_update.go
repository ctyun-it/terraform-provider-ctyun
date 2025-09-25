package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// UserGroupUpdateApi 修改用户组
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9153&data=114
type UserGroupUpdateApi struct {
	builder ctyunsdk.CtyunRequestBuilder
	client  *ctyunsdk.CtyunClient
}

func NewUserGroupUpdateApi(client *ctyunsdk.CtyunClient) *UserGroupUpdateApi {
	return &UserGroupUpdateApi{
		builder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/userGroup/updateGroup",
		},
		client: client,
	}
}

func (this *UserGroupUpdateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserGroupUpdateRequest) (*UserGroupUpdateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.builder.WithCredential(&credential)
	builder, err := builder.WriteJson(&userGroupUpdateRealRequest{
		Id:         req.Id,
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
	resp := &userGroupUpdateRealResponse{}
	err = send.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return &UserGroupUpdateResponse{
		Id:         resp.Id,
		GroupName:  resp.GroupName,
		GroupIntro: resp.GroupIntro,
		AccountId:  resp.AccountId,
	}, nil
}

type userGroupUpdateRealRequest struct {
	Id         string `json:"id"`
	GroupName  string `json:"groupName"`
	GroupIntro string `json:"groupIntro"`
}

type userGroupUpdateRealResponse struct {
	Id         string `json:"id"`
	GroupName  string `json:"groupName"`
	GroupIntro string `json:"groupIntro"`
	AccountId  string `json:"accountId"`
}

type UserGroupUpdateRequest struct {
	Id         string
	GroupName  string
	GroupIntro string
}

type UserGroupUpdateResponse struct {
	Id         string
	GroupName  string
	GroupIntro string
	AccountId  string
}
