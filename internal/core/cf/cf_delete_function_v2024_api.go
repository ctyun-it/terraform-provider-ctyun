package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfDeleteFunctionV2024Api
/* 删除函数 */
type CfDeleteFunctionV2024Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfDeleteFunctionV2024Api(client *core.CtyunClient) *CfDeleteFunctionV2024Api {
	return &CfDeleteFunctionV2024Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/openapi/v1/functions/{functionName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfDeleteFunctionV2024Api) Do(ctx context.Context, credential core.Credential, req *CfDeleteFunctionV2024Request) (*CfDeleteFunctionV2024Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfDeleteFunctionV2024Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfDeleteFunctionV2024Request struct {
	FunctionName string `json:"functionName"` /*  函数名  */
	RegionId     string `json:"regionId"`     /*  资源池id  */
}

type CfDeleteFunctionV2024Response struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string `json:"message"`    /*  错误描述信息  */
}
