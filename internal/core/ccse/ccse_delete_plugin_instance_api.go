package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CcseDeletePluginInstanceApi
/* 调用该接口卸载集群插件。
 */type CcseDeletePluginInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseDeletePluginInstanceApi(client *core.CtyunClient) *CcseDeletePluginInstanceApi {
	return &CcseDeletePluginInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/v2/cce/clusters/{clusterId}/plugininstance/{instanceName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseDeletePluginInstanceApi) Do(ctx context.Context, credential core.Credential, req *CcseDeletePluginInstanceRequest) (*CcseDeletePluginInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("instanceName", req.InstanceName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("namespace", req.Namespace)
	if req.IsDeleteRecord != nil {
		ctReq.AddParam("isDeleteRecord", strconv.FormatBool(*req.IsDeleteRecord))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseDeletePluginInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseDeletePluginInstanceRequest struct {
	ClusterId    string `json:"clusterId,omitempty"`    /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	InstanceName string `json:"instanceName,omitempty"` /*  插件实例名称  */
	RegionId     string `json:"regionId,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	Namespace      string `json:"namespace,omitempty"` /*  命名空间名称  */
	IsDeleteRecord *bool  `json:"isDeleteRecord"`      /*  是否保留卸载记录（默认true，保留记录）  */
}

type CcseDeletePluginInstanceResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求id  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  返回对象  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
