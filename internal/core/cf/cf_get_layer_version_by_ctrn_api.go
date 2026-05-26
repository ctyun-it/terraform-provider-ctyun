package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetLayerVersionByCtrnApi
/* 通过CTRN获取层版本信息 */
type CfGetLayerVersionByCtrnApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetLayerVersionByCtrnApi(client *core.CtyunClient) *CfGetLayerVersionByCtrnApi {
	return &CfGetLayerVersionByCtrnApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/layerctrn/{layerCtrn}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetLayerVersionByCtrnApi) Do(ctx context.Context, credential core.Credential, req *CfGetLayerVersionByCtrnRequest) (*CfGetLayerVersionByCtrnResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("layerCtrn", req.LayerCtrn)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetLayerVersionByCtrnResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetLayerVersionByCtrnRequest struct {
	LayerCtrn string `json:"layerCtrn"` /*  层的ctrn  */
	RegionId  string `json:"regionId"`  /*  资源池id  */
}

type CfGetLayerVersionByCtrnResponse struct {
	StatusCode *int32                                    `json:"statusCode"` /*  状态码  */
	Error      *string                                   `json:"error"`      /*  错误码  */
	Message    *string                                   `json:"message"`    /*  异常信息  */
	ReturnObj  *CfGetLayerVersionByCtrnReturnObjResponse `json:"returnObj"`  /*  返回结果  */
}

type CfGetLayerVersionByCtrnReturnObjResponse struct {
	LayerName         *string                                       `json:"layerName"`         /*  层名  */
	Version           *int32                                        `json:"version"`           /*  版本  */
	Description       *string                                       `json:"description"`       /*  版本描述信息  */
	CompatibleRuntime []*string                                     `json:"compatibleRuntime"` /*  版本运行时环境列表  */
	Ctrn              *string                                       `json:"ctrn"`              /*  版本唯一标识  */
	Code              *CfGetLayerVersionByCtrnReturnObjCodeResponse `json:"code"`              /*  版本代码配置  */
	Codesize          *int32                                        `json:"codesize"`          /*  代码大小  */
	CodeChecksum      *string                                       `json:"codeChecksum"`      /*  代码校验码  */
	CreateTime        *string                                       `json:"createTime"`        /*  版本创建时间  */
	BuildStatus       *bool                                         `json:"buildStatus"`       /*  版本构建状态，false 表示构建中，true 表示构建成功  */
}

type CfGetLayerVersionByCtrnReturnObjCodeResponse struct {
	OssBucketName *string `json:"ossBucketName"` /*  oss的bucket  */
	OssObjectName *string `json:"ossObjectName"` /*  oss的name  */
}
