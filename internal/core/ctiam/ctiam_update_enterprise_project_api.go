package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamUpdateEnterpriseProjectApi
/* 修改企业项目 */
type CtiamUpdateEnterpriseProjectApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamUpdateEnterpriseProjectApi(client *core.CtyunClient) *CtiamUpdateEnterpriseProjectApi {
	return &CtiamUpdateEnterpriseProjectApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/updateEnterpriseProject",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamUpdateEnterpriseProjectApi) Do(ctx context.Context, credential core.Credential, req *CtiamUpdateEnterpriseProjectRequest) (*CtiamUpdateEnterpriseProjectResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamUpdateEnterpriseProjectRequest
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
	var resp CtiamUpdateEnterpriseProjectResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamUpdateEnterpriseProjectRequest struct {
	Id          string  `json:"id"`                    /*  企业项目id  */
	ProjectName *string `json:"projectName,omitempty"` /*  企业项目名称  */
	Description *string `json:"description,omitempty"` /*  企描述  */
}

type CtiamUpdateEnterpriseProjectResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
