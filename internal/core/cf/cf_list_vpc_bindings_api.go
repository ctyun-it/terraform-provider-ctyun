package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfListVpcBindingsApi
/* 列出VPC绑定配置 */
type CfListVpcBindingsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListVpcBindingsApi(client *core.CtyunClient) *CfListVpcBindingsApi {
	return &CfListVpcBindingsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions/{functionName}/vpc-bindings",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListVpcBindingsApi) Do(ctx context.Context, credential core.Credential, req *CfListVpcBindingsRequest) (*CfListVpcBindingsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListVpcBindingsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListVpcBindingsRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称  */
	RegionId     string `json:"regionId"`     /*  资源池id  */
}

type CfListVpcBindingsResponse struct {
	StatusCode *int32    `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string   `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string   `json:"message"`    /*  错误描述信息  */
	ReturnObj  []*string `json:"returnObj"`  /*  vpcIds  */
}
