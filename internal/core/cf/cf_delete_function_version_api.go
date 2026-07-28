package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfDeleteFunctionVersionApi
/* 删除版本 */
type CfDeleteFunctionVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfDeleteFunctionVersionApi(client *core.CtyunClient) *CfDeleteFunctionVersionApi {
	return &CfDeleteFunctionVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/openapi/v1/functions/{functionName}/versions/{versionId}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfDeleteFunctionVersionApi) Do(ctx context.Context, credential core.Credential, req *CfDeleteFunctionVersionRequest) (*CfDeleteFunctionVersionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder = builder.ReplaceUrl("versionId", req.VersionId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfDeleteFunctionVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfDeleteFunctionVersionRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称  */
	VersionId    string `json:"versionId"`    /*  版本ID  */
	RegionId     string `json:"regionId"`     /*  资源池id  */
}

type CfDeleteFunctionVersionResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string `json:"error"`      /*  错误码  */
	Message    *string `json:"message"`    /*  信息  */
	ReturnObj  *bool   `json:"returnObj"`  /*  删除结果，删除成功时为true  */
}
