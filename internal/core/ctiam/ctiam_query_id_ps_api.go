package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryIdPsApi
/* 分页查询身份供应商 */
type CtiamQueryIdPsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryIdPsApi(client *core.CtyunClient) *CtiamQueryIdPsApi {
	return &CtiamQueryIdPsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/identityProvider/queryIdPs",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryIdPsApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryIdPsRequest) (*CtiamQueryIdPsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryIdPsRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamQueryIdPsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryIdPsRequest struct {
	PageNum  int32   `json:"pageNum"`        /*  页码  */
	PageSize int32   `json:"pageSize"`       /*  每页条数  */
	Name     *string `json:"name,omitempty"` /*  身份提供商名称（模糊搜索）  */
}

type CtiamQueryIdPsResponse struct {
	StatusCode *string                          `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                          `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *CtiamQueryIdPsReturnObjResponse `json:"returnObj"`  /*  返回参数  */
	Error      *string                          `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

type CtiamQueryIdPsReturnObjResponse struct {
	List     []*CtiamQueryIdPsReturnObjListResponse `json:"list"`     /*  身份提供商列表  */
	PageNum  *int32                                 `json:"pageNum"`  /*  页码  */
	PageSize *int32                                 `json:"pageSize"` /*  每页条数  */
	Total    *int32                                 `json:"total"`    /*  总条数  */
	Pages    *int32                                 `json:"pages"`    /*  总页数  */
}

type CtiamQueryIdPsReturnObjListResponse struct {
	Id         *int32  `json:"id"`         /*  id  */
	Name       *string `json:"name"`       /*  名称  */
	RawType    *int32  `json:"type"`       /*  类型，0 虚拟用户SSO，1 IAM用户SSO  */
	Protocol   *int32  `json:"protocol"`   /*  协议，0 SAML协议，1 OIDC协议  */
	AccountId  *string `json:"accountId"`  /*  账号id  */
	Remark     *string `json:"remark"`     /*  描述  */
	Status     *int32  `json:"status"`     /*  状态  */
	CreateTime *string `json:"createTime"` /*  创建时间  */
	UpdateTime *string `json:"updateTime"` /*  更新时间  */
	FileName   *string `json:"fileName"`   /*  文件名  */
}
