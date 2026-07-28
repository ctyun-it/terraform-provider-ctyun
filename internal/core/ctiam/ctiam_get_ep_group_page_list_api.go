package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetEpGroupPageListApi
/* 企业项目关联用户组分页查询 */
type CtiamGetEpGroupPageListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetEpGroupPageListApi(client *core.CtyunClient) *CtiamGetEpGroupPageListApi {
	return &CtiamGetEpGroupPageListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/getEpGroupPageList",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetEpGroupPageListApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetEpGroupPageListRequest) (*CtiamGetEpGroupPageListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamGetEpGroupPageListRequest
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
	var resp CtiamGetEpGroupPageListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetEpGroupPageListRequest struct {
	Id       string `json:"id"`       /*  企业项目id  */
	PageNum  int32  `json:"pageNum"`  /*  页码  */
	PageSize int32  `json:"pageSize"` /*  每页条数  */
}

type CtiamGetEpGroupPageListResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *string `json:"returnObj"`  /*  返回参数  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
