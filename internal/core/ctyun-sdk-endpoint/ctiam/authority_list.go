package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
	"strconv"
)

// AuthorityListApi 根据云服务ID查询云服务权限点
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=77&api=13943&data=114
type AuthorityListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAuthorityListApi(client *ctyunsdk.CtyunClient) *AuthorityListApi {
	return &AuthorityListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v1/service/queryAllAuthorityByServiceId",
		},
	}
}

func (this *AuthorityListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AuthorityListRequest) (*AuthorityListResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("serviceId", strconv.Itoa(req.ServiceId))

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtiam, builder)
	if err != nil {
		return nil, err
	}

	var resp authorityListRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	var authorityList []AuthorityListAuthorityListResponse
	for _, authority := range resp.AuthorityList {
		authorityList = append(authorityList, AuthorityListAuthorityListResponse{
			ServiceId:      authority.ServiceId,
			Name:           authority.Name,
			Code:           authority.Code,
			Description:    authority.Description,
			CtrntemplateId: authority.CtrntemplateId,
		})
	}
	return &AuthorityListResponse{
		AuthorityList: authorityList,
	}, nil
}

type authorityListRealResponse struct {
	AuthorityList []struct {
		ServiceId      int    `json:"serviceId"`
		Name           string `json:"name"`
		Code           string `json:"code"`
		Description    string `json:"description"`
		CtrntemplateId string `json:"ctrntemplateId"`
	} `json:"authorityList"`
}

type AuthorityListRequest struct {
	ServiceId int
}

type AuthorityListAuthorityListResponse struct {
	ServiceId      int
	Name           string
	Code           string
	Description    string
	CtrntemplateId string
}

type AuthorityListResponse struct {
	AuthorityList []AuthorityListAuthorityListResponse
}
