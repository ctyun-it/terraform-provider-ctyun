package ctiam

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// IdpDeleteApi 删除身份供应商
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9184&data=114
type IdpDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewIdpDeleteApi(client *ctyunsdk.CtyunClient) *IdpDeleteApi {
	return &IdpDeleteApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/identityProvider/deleteIdP",
		},
		client: client,
	}

}

func (this IdpDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *IdpDeleteRequest) (*IdpDeleteResponse, ctyunsdk.CtyunRequestError) {
	builder := this.CtyunRequestBuilder.WithCredential(&credential)
	builder, err := builder.WriteJson(&idpDeleteRealRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	send, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	err = send.ParseByStandardModelWithCheck(nil)
	return &IdpDeleteResponse{}, err
}

type idpDeleteRealRequest struct {
	Id int64 `json:"id"`
}

type IdpDeleteRequest struct {
	Id int64
}

type IdpDeleteResponse struct {
}
