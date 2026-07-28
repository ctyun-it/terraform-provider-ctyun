package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtiamSasEnterpriseProjectApi
/* 启用停用企业项目 */
type CtiamSasEnterpriseProjectApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamSasEnterpriseProjectApi(client *core.CtyunClient) *CtiamSasEnterpriseProjectApi {
	return &CtiamSasEnterpriseProjectApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/sasEnterpriseProject",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamSasEnterpriseProjectApi) Do(ctx context.Context, credential core.Credential, req *CtiamSasEnterpriseProjectRequest) (*CtiamSasEnterpriseProjectResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("projectId", req.ProjectId)
	ctReq.AddParam("status", strconv.FormatInt(int64(req.Status), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamSasEnterpriseProjectResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamSasEnterpriseProjectRequest struct {
	ProjectId string `json:"projectId"` /*  企业项目id  */
	Status    int32  `json:"status"`    /*  状态（1 启用；2 停用）  */
}

type CtiamSasEnterpriseProjectResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	ReturnObj  *string `json:"returnObj"`  /*  返回参数  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
