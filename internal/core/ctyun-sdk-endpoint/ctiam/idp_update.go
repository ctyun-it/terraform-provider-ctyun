package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// IdpUpdateApi 删除身份供应商
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9182&data=114
type IdpUpdateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewIdpUpdateApi(client *ctyunsdk.CtyunClient) *IdpUpdateApi {
	return &IdpUpdateApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v1/identityProvider/updateIdP",
		},
		client: client,
	}
}

func (this IdpUpdateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *IdpUpdateRequest) (*IdpUpdateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.CtyunRequestBuilder.WithCredential(&credential)
	builder, err := builder.WriteJson(&idpUpdateRealRequest{
		Id:       req.Id,
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
	err = send.ParseByStandardModelWithCheck(nil)
	return &IdpUpdateResponse{}, err
}

type idpUpdateRealRequest struct {
	Id       int64  `json:"id"`
	Remark   string `json:"remark"`
	FileName string `json:"fileName"`
	File     []byte `json:"file"`
}

type IdpUpdateRequest struct {
	Id       int64
	Remark   string
	FileName string
	File     []byte
}

type IdpUpdateResponse struct {
}
