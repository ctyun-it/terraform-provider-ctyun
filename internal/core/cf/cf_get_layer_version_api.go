package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetLayerVersionApi
/* 获取层版本信息 */
type CfGetLayerVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetLayerVersionApi(client *core.CtyunClient) *CfGetLayerVersionApi {
	return &CfGetLayerVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/layers/{layerName}/versions/{version}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetLayerVersionApi) Do(ctx context.Context, credential core.Credential, req *CfGetLayerVersionRequest) (*CfGetLayerVersionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("layerName", req.LayerName)
	builder = builder.ReplaceUrl("version", req.Version)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetLayerVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetLayerVersionRequest struct {
	LayerName string `json:"layerName"` /*  层的名称。只能包含字母、数字、下划线和中划线。只能字母开头。长度在 1-64字符之间  */
	Version   string `json:"version"`   /*  层的版本  */
	RegionId  string `json:"regionId"`  /*  资源池id  */
}

type CfGetLayerVersionResponse struct {
	StatusCode *int32                              `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string                             `json:"error"`      /*  错误码  */
	Message    *string                             `json:"message"`    /*  信息  */
	ReturnObj  *CfGetLayerVersionReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfGetLayerVersionReturnObjResponse struct {
	LayerName    *string                                 `json:"layerName"`    /*  层名  */
	Version      *int32                                  `json:"version"`      /*  版本  */
	Description  *string                                 `json:"description"`  /*  版本描述信息  */
	Ctrn         *string                                 `json:"ctrn"`         /*  版本唯一标识  */
	Code         *CfGetLayerVersionReturnObjCodeResponse `json:"code"`         /*  版本代码配置  */
	Codesize     *int32                                  `json:"codesize"`     /*  代码大小  */
	CodeChecksum *string                                 `json:"codeChecksum"` /*  代码校验码  */
	CreateTime   *string                                 `json:"createTime"`   /*  版本创建时间  */
	BuildStatus  *bool                                   `json:"buildStatus"`  /*  版本构建状态，false 表示构建中，true 表示构建成功  */
}

type CfGetLayerVersionReturnObjCodeResponse struct {
	OssBucketName *string `json:"ossBucketName"` /*  oss的bucket  */
	OssObjectName *string `json:"ossObjectName"` /*  oss的name  */
}
