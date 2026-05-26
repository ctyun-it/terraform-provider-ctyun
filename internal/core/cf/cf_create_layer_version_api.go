package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfCreateLayerVersionApi
/* 创建层版本 */
type CfCreateLayerVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfCreateLayerVersionApi(client *core.CtyunClient) *CfCreateLayerVersionApi {
	return &CfCreateLayerVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/layers/{layerName}/versions",
			ContentType:  "application/json",
		},
	}
}

func (a *CfCreateLayerVersionApi) Do(ctx context.Context, credential core.Credential, req *CfCreateLayerVersionRequest) (*CfCreateLayerVersionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("layerName", req.LayerName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfCreateLayerVersionRequest
		RegionId  interface{} `json:"regionId,omitempty"`
		LayerName interface{} `json:"layerName,omitempty"`
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
	var resp CfCreateLayerVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfCreateLayerVersionRequest struct {
	LayerName         string                           `json:"layerName"`             /*  只能包含字母、数字、下划线和中划线。只能字母开头。长度在 1-64字符之间  */
	RegionId          string                           `json:"regionId"`              /*  资源池id  */
	Description       *string                          `json:"description,omitempty"` /*  描述信息  */
	CompatibleRuntime []string                         `json:"compatibleRuntime"`     /*  运行时环境列表  */
	Code              *CfCreateLayerVersionCodeRequest `json:"code,omitempty"`        /*  层代码  */
}

type CfCreateLayerVersionCodeRequest struct {
	OssBucketName *string `json:"ossBucketName,omitempty"` /*  oss的bucket  */
	OssObjectName string  `json:"ossObjectName"`           /*  oss的name  */
}

type CfCreateLayerVersionResponse struct {
	StatusCode *int32                                 `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string                                `json:"error"`      /*  错误码  */
	Message    *string                                `json:"message"`    /*  信息  */
	ReturnObj  *CfCreateLayerVersionReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfCreateLayerVersionReturnObjResponse struct {
	LayerName         *string                                    `json:"layerName"`         /*  层名  */
	Version           *int32                                     `json:"version"`           /*  版本  */
	Description       *string                                    `json:"description"`       /*  描述信息  */
	CompatibleRuntime []*string                                  `json:"compatibleRuntime"` /*  运行时环境列表  */
	Ctrn              *string                                    `json:"ctrn"`              /*  唯一标识  */
	Code              *CfCreateLayerVersionReturnObjCodeResponse `json:"code"`              /*  代码配置  */
	Codesize          *int32                                     `json:"codesize"`          /*  代码大小  */
	CodeChecksum      *string                                    `json:"codeChecksum"`      /*  代码校验码  */
	CreateTime        *string                                    `json:"createTime"`        /*  创建时间  */
	Acl               *int32                                     `json:"acl"`               /*  公共层/私有层  */
}

type CfCreateLayerVersionReturnObjCodeResponse struct {
	OssBucketName *string `json:"ossBucketName"` /*  oss的bucket  */
	OssObjectName *string `json:"ossObjectName"` /*  oss的name  */
}
