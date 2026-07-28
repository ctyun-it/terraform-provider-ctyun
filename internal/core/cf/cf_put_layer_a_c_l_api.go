package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfPutLayerACLApi
/* 设置层的权限 */
type CfPutLayerACLApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfPutLayerACLApi(client *core.CtyunClient) *CfPutLayerACLApi {
	return &CfPutLayerACLApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/layers/{layerName}/acl",
			ContentType:  "application/json",
		},
	}
}

func (a *CfPutLayerACLApi) Do(ctx context.Context, credential core.Credential, req *CfPutLayerACLRequest) (*CfPutLayerACLResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("layerName", req.LayerName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfPutLayerACLRequest
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
	var resp CfPutLayerACLResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfPutLayerACLRequest struct {
	LayerName string `json:"layerName"` /*  层的名称。只能包含字母、数字、下划线和中划线。只能字母开头。长度在 1-64字符之间  */
	RegionId  string `json:"regionId"`  /*  资源池id  */
	IsPublic  bool   `json:"isPublic"`  /*  是否公开  */
}

type CfPutLayerACLResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string `json:"error"`      /*  错误码  */
	Message    *string `json:"message"`    /*  信息  */
}
