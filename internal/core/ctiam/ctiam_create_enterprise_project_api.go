package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamCreateEnterpriseProjectApi
/* 创建企业项目 */
type CtiamCreateEnterpriseProjectApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamCreateEnterpriseProjectApi(client *core.CtyunClient) *CtiamCreateEnterpriseProjectApi {
	return &CtiamCreateEnterpriseProjectApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/createEnterpriseProject",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamCreateEnterpriseProjectApi) Do(ctx context.Context, credential core.Credential, req *CtiamCreateEnterpriseProjectRequest) (*CtiamCreateEnterpriseProjectResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamCreateEnterpriseProjectRequest
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
	var resp CtiamCreateEnterpriseProjectResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamCreateEnterpriseProjectRequest struct {
	ProjectName string `json:"projectName"` /*  企业项目名称  */
	Description string `json:"description"` /*  描述  */
}

type CtiamCreateEnterpriseProjectResponse struct {
	StatusCode *string                                        `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                        `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                        `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamCreateEnterpriseProjectReturnObjResponse `json:"returnObj"`  /*  返回体  */
}

type CtiamCreateEnterpriseProjectReturnObjResponse struct {
	ProjectId *string `json:"projectId"` /*  企业项目id  */
}
