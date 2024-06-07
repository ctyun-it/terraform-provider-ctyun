package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// UserUpdateApi 修改用户
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9148&data=114
type UserUpdateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewUserUpdateApi(client *ctyunsdk.CtyunClient) *UserUpdateApi {
	return &UserUpdateApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/user/updateUser",
		},
		client: client,
	}
}

func (this UserUpdateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *UserUpdateRequest) (*UserUpdateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.CtyunRequestBuilder.WithCredential(&credential)
	builder, err := builder.WriteJson(&userUpdateRealRequest{
		UserId:      req.UserId,
		Remark:      req.Remark,
		LoginEmail:  req.LoginEmail,
		MobilePhone: req.MobilePhone,
		UserName:    req.UserName,
		Prohibit:    req.Prohibit,
	})
	if err != nil {
		return nil, err
	}
	send, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	resp := &userUpdateRealResponse{}
	err = send.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return &UserUpdateResponse{
		LoginEmail:  resp.LoginEmail,
		AccountId:   resp.AccountId,
		MobilePhone: resp.MobilePhone,
		Remark:      resp.Remark,
		UserName:    resp.UserName,
	}, nil
}

type userUpdateRealRequest struct {
	UserId      string `json:"userId"`
	Remark      string `json:"remark,"`
	LoginEmail  string `json:"loginEmail"`
	MobilePhone string `json:"mobilePhone"`
	UserName    string `json:"userName"`
	Prohibit    int    `json:"prohibit"`
}

type userUpdateRealResponse struct {
	LoginEmail  string `json:"loginEmail"`
	AccountId   string `json:"accountId"`
	MobilePhone string `json:"mobilePhone"`
	Remark      string `json:"remark"`
	UserName    string `json:"userName"`
}

type UserUpdateRequest struct {
	UserId      string
	Remark      string
	LoginEmail  string
	MobilePhone string
	UserName    string
	Prohibit    int
}

type UserUpdateResponse struct {
	LoginEmail  string
	AccountId   string
	MobilePhone string
	Remark      string
	UserName    string
}
