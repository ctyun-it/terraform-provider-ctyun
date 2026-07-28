package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryQuotaByTypeApi
/* 根据配额类型查询配额列表 */
type CtiamQueryQuotaByTypeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryQuotaByTypeApi(client *core.CtyunClient) *CtiamQueryQuotaByTypeApi {
	return &CtiamQueryQuotaByTypeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/quota/queryQuotaByType",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryQuotaByTypeApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryQuotaByTypeRequest) (*CtiamQueryQuotaByTypeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryQuotaByTypeRequest
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
	var resp CtiamQueryQuotaByTypeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryQuotaByTypeRequest struct {
	RawType int8 `json:"type"` /*  配额类型（1代表用户配额, 2代表用户组配额, 3代表策略配额）  */
}

type CtiamQueryQuotaByTypeResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *string `json:"returnObj"`  /*  返回参数  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
