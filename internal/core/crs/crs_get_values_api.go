package crs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CrsGetValuesApi
/* 查询版本values
 */type CrsGetValuesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCrsGetValuesApi(client *core.CtyunClient) *CrsGetValuesApi {
	return &CrsGetValuesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/getValues",
			ContentType:  "application/json",
		},
	}
}

func (a *CrsGetValuesApi) Do(ctx context.Context, credential core.Credential, req *CrsGetValuesRequest) (*CrsGetValuesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("Content-Type", req.ContentType)
	ctReq.AddHeader("regionId", req.RegionIdHeader)
	ctReq.AddParam("regionId", req.RegionIdParam)
	ctReq.AddParam("namespaceName", req.NamespaceName)
	ctReq.AddParam("repositoryName", req.RepositoryName)
	ctReq.AddParam("tagName", req.TagName)
	if req.RawType != "" {
		ctReq.AddParam("type", req.RawType)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CrsGetValuesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CrsGetValuesRequest struct {
	ContentType    string `json:"Content-Type,omitempty"`   /*  类型  */
	RegionIdHeader string `json:"regionId,omitempty"`       /*  资源池编码（资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026765/11005378" target="_blank">容器镜像服务资源池</a>获取）  */
	RegionIdParam  string `json:"regionId,omitempty"`       /*  资源池编码（资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026765/11005378" target="_blank">容器镜像服务资源池</a>获取）  */
	NamespaceName  string `json:"namespaceName,omitempty"`  /*  命名空间名称  */
	RepositoryName string `json:"repositoryName,omitempty"` /*  插件名称  */
	TagName        string `json:"tagName,omitempty"`        /*  版本名称  */
	RawType        string `json:"type,omitempty"`           /*  values的格式，默认为YAML(YAML:YAML格式, JSON:JSON格式)  */
}

type CrsGetValuesResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应码（0为请求成功，其它为失败 ）  */
	Message    string `json:"message,omitempty"`    /*  返回信息  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
	ReturnObj  string `json:"returnObj,omitempty"`  /*  values详情  */
}
