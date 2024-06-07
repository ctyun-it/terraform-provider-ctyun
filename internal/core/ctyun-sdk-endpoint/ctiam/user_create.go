package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// UserCreateApi 创建用户
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=8135&data=114
type UserCreateApi struct {
	builder ctyunsdk.CtyunRequestBuilder
	client  *ctyunsdk.CtyunClient
}

func NewUserCreateApi(client *ctyunsdk.CtyunClient) *UserCreateApi {
	return &UserCreateApi{
		builder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/user/createUser",
		},
		client: client,
	}
}

func (this *UserCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserCreateRequest) (*UserCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.builder.WithCredential(&credential)
	builder, err := builder.WriteJson(req)
	if err != nil {
		return nil, err
	}
	send, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	resp := &UserCreateResponse{}
	err = send.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type UserCreateRequest struct {
	LoginEmail         string      `json:"loginEmail"`
	MobilePhone        string      `json:"mobilePhone"`
	Password           string      `json:"password,omitempty"`
	UserName           string      `json:"userName"`
	Remark             string      `json:"remark,omitempty"`
	Groups             []UserGroup `json:"groups"`
	GeneratePassword   bool        `json:"generatePassword"`
	LoginResetPassword bool        `json:"loginResetPassword"`
	SourcePassword     string      `json:"sourcePassword"`
}

type UserGroup struct {
	Id string `json:"id"`
}

type UserCreateResponse struct {
	AccountId        string      `json:"accountId"`
	Groups           []UserGroup `json:"groups"`
	IsVirtualAccount string      `json:"isVirtualAccount"`
	LoginEmail       string      `json:"loginEmail"`
	MobilePhone      string      `json:"mobilePhone"`
	Remark           string      `json:"remark"`
	UserId           string      `json:"userId"`
	UserName         string      `json:"userName"`
}
