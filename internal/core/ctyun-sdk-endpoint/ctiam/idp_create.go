package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// IdpCreateApi 创建身份提供商
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9183&data=114
type IdpCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewIdpCreateApi(client *ctyunsdk.CtyunClient) *IdpCreateApi {
	return &IdpCreateApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/identityProvider/createIdP",
		},
		client: client,
	}
}

func (this IdpCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *IdpCreateRequest) (*IdpCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.CtyunRequestBuilder.WithCredential(&credential)
	builder, err := builder.WriteJson(&idpCreateRealRequest{
		Name:     req.Name,
		Type:     req.Type,
		Protocol: req.Protocol,
		Remark:   req.Remark,
		FileName: req.FileName,
		File:     req.File,
	})
	if err != nil {
		return nil, err
	}
	send, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	resp := &idpCreateRealResponse{}
	err = send.ParseByStandardModelWithCheck(resp)
	return &IdpCreateResponse{
		Id:         resp.Id,
		Name:       resp.Name,
		Type:       resp.Type,
		Protocol:   resp.Protocol,
		AccountId:  resp.AccountId,
		Remark:     resp.Remark,
		CreateTime: resp.CreateTime,
		Uuid:       resp.Uuid,
	}, err
}

type idpCreateRealRequest struct {
	Name     string `json:"name"`
	Type     int    `json:"type"`
	Protocol int    `json:"protocol"`
	Remark   string `json:"remark"`
	FileName string `json:"fileName"`
	File     []byte `json:"file"`
}

type idpCreateRealResponse struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	Protocol   int    `json:"protocol"`
	AccountId  string `json:"accountId"`
	Remark     string `json:"remark"`
	CreateTime int64  `json:"createTime"`
	Uuid       string `json:"uuid"`
}

type IdpCreateRequest struct {
	Name     string
	Type     int
	Protocol int
	Remark   string
	FileName string
	File     []byte
}

type IdpCreateResponse struct {
	Id         int64
	Name       string
	Type       int
	Protocol   int
	AccountId  string
	Remark     string
	CreateTime int64
	Uuid       string
}
