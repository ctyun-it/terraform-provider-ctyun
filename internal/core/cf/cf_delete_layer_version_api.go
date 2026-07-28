package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfDeleteLayerVersionApi
/* 删除层版本 */
type CfDeleteLayerVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfDeleteLayerVersionApi(client *core.CtyunClient) *CfDeleteLayerVersionApi {
	return &CfDeleteLayerVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/openapi/v1/layers/{layerName}/versions/{version}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfDeleteLayerVersionApi) Do(ctx context.Context, credential core.Credential, req *CfDeleteLayerVersionRequest) (*CfDeleteLayerVersionResponse, error) {
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
	var resp CfDeleteLayerVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfDeleteLayerVersionRequest struct {
	LayerName string `json:"layerName"` /*  层的名称  */
	Version   string `json:"version"`   /*  层的版本  */
	RegionId  string `json:"regionId"`  /*  资源池id  */
}

type CfDeleteLayerVersionResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string `json:"error"`      /*  错误码  */
	Message    *string `json:"message"`    /*  信息  */
}
