package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfPublishFunctionVersionApi
/* 发布/创建版本 */
type CfPublishFunctionVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfPublishFunctionVersionApi(client *core.CtyunClient) *CfPublishFunctionVersionApi {
	return &CfPublishFunctionVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/functions/{functionName}/versions",
			ContentType:  "application/json",
		},
	}
}

func (a *CfPublishFunctionVersionApi) Do(ctx context.Context, credential core.Credential, req *CfPublishFunctionVersionRequest) (*CfPublishFunctionVersionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfPublishFunctionVersionRequest
		RegionId     interface{} `json:"regionId,omitempty"`
		FunctionName interface{} `json:"functionName,omitempty"`
	}{
		req, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfPublishFunctionVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfPublishFunctionVersionRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称  */
	RegionId     string `json:"regionId"`     /*  资源池id  */
	Description  string `json:"description"`  /*  描述  */
}

type CfPublishFunctionVersionResponse struct {
	StatusCode *int32                                     `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string                                    `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                    `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfPublishFunctionVersionReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfPublishFunctionVersionReturnObjResponse struct {
	VersionId   *string `json:"versionId"`   /*  每次版本后自动增1  */
	Description *string `json:"description"` /*  版本描述  */
	CreateTime  *int32  `json:"createTime"`  /*  版本创建unix时间戳(秒)  */
}
