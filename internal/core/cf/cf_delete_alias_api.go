package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfDeleteAliasApi
/* 删除别名 */
type CfDeleteAliasApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfDeleteAliasApi(client *core.CtyunClient) *CfDeleteAliasApi {
	return &CfDeleteAliasApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/openapi/v1/functions/{functionName}/aliases/{aliasName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfDeleteAliasApi) Do(ctx context.Context, credential core.Credential, req *CfDeleteAliasRequest) (*CfDeleteAliasResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder = builder.ReplaceUrl("aliasName", req.AliasName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfDeleteAliasResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfDeleteAliasRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称，函数必须存在  */
	AliasName    string `json:"aliasName"`    /*  别名名称，别名必须存在  */
	RegionId     string `json:"regionId"`     /*  资源池id，标识不同的地区，如：华东1、西南1  */
}

type CfDeleteAliasResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码。0表示成功，非0表示不成功  */
	Error      *string `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    *string `json:"message"`    /*  错误描述信息  */
	ReturnObj  *bool   `json:"returnObj"`  /*  删除结果，删除成功时为true  */
}
