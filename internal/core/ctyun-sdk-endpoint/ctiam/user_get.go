package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// UserGetApi 根据id查询用户详情
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9149&data=114
type UserGetApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewUserGetApi(client *ctyunsdk.CtyunClient) *UserGetApi {
	return &UserGetApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v1/user/getUser",
		},
		client: client,
	}
}

func (this *UserGetApi) Do(ctx context.Context, credential ctyunsdk.Credential, r *UserGetRequest) (*UserGetResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential).AddParam("userId", r.UserId)
	ctiam, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, request)
	if err != nil {
		return nil, err
	}
	resp := &UserGetResponse{}
	err = ctiam.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type UserGetRequest struct {
	UserId string
}

type UserGetResponse struct {
	LoginEmail  string      `json:"loginEmail"`
	AccountId   string      `json:"accountId"`
	MobilePhone string      `json:"mobilePhone"`
	Groups      []UserGroup `json:"groups"`
	Remark      string      `json:"remark"`
	UserName    string      `json:"userName"`
}
