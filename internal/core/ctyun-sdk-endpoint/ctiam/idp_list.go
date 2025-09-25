package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
	"strconv"
)

// IdpListApi 查看身份供应商详情
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=9186&data=114
type IdpListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewIdpListApi(client *ctyunsdk.CtyunClient) *IdpListApi {
	return &IdpListApi{
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v1/identityProvider/getIdPDetails",
		},
		client: client,
	}
}

func (this *IdpListApi) Do(ctx context.Context, credential ctyunsdk.Credential, r *IdpListRequest) (*IdpListResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential).AddParam("id", strconv.FormatInt(r.Id, 10))
	ctiam, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}
	resp := &idpListRealResponse{}
	err = ctiam.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return &IdpListResponse{
		Id:               resp.Id,
		Protocol:         resp.Protocol,
		AccountId:        resp.AccountId,
		Remark:           resp.Remark,
		Name:             resp.Name,
		Type:             resp.Type,
		Status:           resp.Status,
		CreateTime:       resp.CreateTime,
		UpdateTime:       resp.UpdateTime,
		MetadataDocument: resp.MetadataDocument,
		FileName:         resp.FileName,
	}, nil
}

type idpListRealResponse struct {
	Id               int64  `json:"id"`
	Protocol         int    `json:"protocol"`
	AccountId        string `json:"accountId"`
	Remark           string `json:"remark"`
	Name             string `json:"name"`
	Type             int    `json:"type"`
	Status           int    `json:"status"`
	CreateTime       int64  `json:"createTime"`
	UpdateTime       int64  `json:"updateTime"`
	MetadataDocument string `json:"metadataDocument"`
	FileName         string `json:"fileName"`
}

type IdpListRequest struct {
	Id int64
}

type IdpListResponse struct {
	Id               int64
	Protocol         int
	AccountId        string
	Remark           string
	Name             string
	Type             int
	Status           int
	CreateTime       int64
	UpdateTime       int64
	MetadataDocument string
	FileName         string
}
