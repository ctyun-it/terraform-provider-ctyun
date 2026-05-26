package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfDeleteTriggerV2024Api
/* 删除触发器 */
type CfDeleteTriggerV2024Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfDeleteTriggerV2024Api(client *core.CtyunClient) *CfDeleteTriggerV2024Api {
	return &CfDeleteTriggerV2024Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/openapi/v1/functions/{functionName}/triggers/{triggerName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfDeleteTriggerV2024Api) Do(ctx context.Context, credential core.Credential, req *CfDeleteTriggerV2024Request) (*CfDeleteTriggerV2024Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder = builder.ReplaceUrl("triggerName", req.TriggerName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfDeleteTriggerV2024Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfDeleteTriggerV2024Request struct {
	FunctionName string `json:"functionName"` /*  函数名称，该函数必须已存在  */
	TriggerName  string `json:"triggerName"`  /*  触发器名称，该触发器必须已存在  */
	RegionId     string `json:"regionId"`     /*  资源池id，标识不同的地区，如：华东1、西南1  */
}

type CfDeleteTriggerV2024Response struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string `json:"message"`    /*  错误提示信息  */
	ReturnObj  *bool   `json:"returnObj"`  /*  是否成功，true：成功；false：失败  */
}
