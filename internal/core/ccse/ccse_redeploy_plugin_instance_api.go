package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseRedeployPluginInstanceApi
/* 调用该接口可重新部署/更新集群插件配置
 */type CcseRedeployPluginInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseRedeployPluginInstanceApi(client *core.CtyunClient) *CcseRedeployPluginInstanceApi {
	return &CcseRedeployPluginInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/plugininstance/{pluginName}/redeploy",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseRedeployPluginInstanceApi) Do(ctx context.Context, credential core.Credential, req *CcseRedeployPluginInstanceRequest) (*CcseRedeployPluginInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("pluginName", req.PluginName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("namespace", req.Namespace)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseRedeployPluginInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseRedeployPluginInstanceRequest struct {
	ClusterId  string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	PluginName string /*  插件实例名称  */
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	Namespace    string /*  命名空间名称  */
	ChartName    string `json:"chartName,omitempty"`    /*  插件名称  */
	ChartVersion string `json:"chartVersion,omitempty"` /*  插件版本号  */
	InstanceName string `json:"instanceName,omitempty"` /*  实例名称  */
	Values       string `json:"values,omitempty"`       /*  实例配置参数（YAML格式），与valuesJson二选一  */
	ValuesJson   string `json:"valuesJson,omitempty"`   /*  实例配置参数（JSON格式），与values二选一  */
}

type CcseRedeployPluginInstanceResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求id  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  返回对象  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
