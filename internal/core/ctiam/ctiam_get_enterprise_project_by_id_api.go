package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetEnterpriseProjectByIdApi
/* 查询企业项目 */
type CtiamGetEnterpriseProjectByIdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetEnterpriseProjectByIdApi(client *core.CtyunClient) *CtiamGetEnterpriseProjectByIdApi {
	return &CtiamGetEnterpriseProjectByIdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/project/getEnterpriseProjectById",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetEnterpriseProjectByIdApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetEnterpriseProjectByIdRequest) (*CtiamGetEnterpriseProjectByIdResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("id", req.Id)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamGetEnterpriseProjectByIdResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetEnterpriseProjectByIdRequest struct {
	Id string `json:"id"` /*  企业项目id  */
}

type CtiamGetEnterpriseProjectByIdResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *string `json:"returnObj"`  /*  返回参数  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
