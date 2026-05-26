package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfGetFunctionCodeApi
/* 获取函数的代码下载链接 */
type CfGetFunctionCodeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetFunctionCodeApi(client *core.CtyunClient) *CfGetFunctionCodeApi {
	return &CfGetFunctionCodeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions/{functionName}/code",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetFunctionCodeApi) Do(ctx context.Context, credential core.Credential, req *CfGetFunctionCodeRequest) (*CfGetFunctionCodeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.Qualifier != nil && *req.Qualifier != "" {
		ctReq.AddParam("qualifier", *req.Qualifier)
	}
	if req.Expire != nil && *req.Expire != 0 {
		ctReq.AddParam("expire", strconv.FormatInt(int64(*req.Expire), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetFunctionCodeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetFunctionCodeRequest struct {
	FunctionName string  `json:"functionName"`        /*  函数名称  */
	RegionId     string  `json:"regionId"`            /*  资源池id  */
	Qualifier    *string `json:"qualifier,omitempty"` /*  服务的版本或者别名，默认是LATEST  */
	Expire       *int32  `json:"expire,omitempty"`    /*  下载链接的失效时间，单位分钟，默认120  */
}

type CfGetFunctionCodeResponse struct {
	StatusCode *int32                              `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string                             `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                             `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfGetFunctionCodeReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfGetFunctionCodeReturnObjResponse struct {
	Url *string `json:"url"` /*  公网下载链接  */
}
