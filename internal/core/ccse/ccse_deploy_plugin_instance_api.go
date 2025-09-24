package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseDeployPluginInstanceApi
/* 调用该接口可在指定集群安装插件
 */type CcseDeployPluginInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseDeployPluginInstanceApi(client *core.CtyunClient) *CcseDeployPluginInstanceApi {
	return &CcseDeployPluginInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/plugininstance",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseDeployPluginInstanceApi) Do(ctx context.Context, credential core.Credential, req *CcseDeployPluginInstanceRequest) (*CcseDeployPluginInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CcseDeployPluginInstanceRequest
		RegionId  interface{} `json:"regionId,omitempty"`
		ClusterId interface{} `json:"clusterId,omitempty"`
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
	var resp CcseDeployPluginInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseDeployPluginInstanceRequest struct {
	ClusterId string `json:"clusterId,omitempty"` /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string `json:"regionId,omitempty"`  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	ChartName    string `json:"chartName,omitempty"`    /*  插件名称  */
	ChartVersion string `json:"chartVersion,omitempty"` /*  插件版本号，可通过容器镜像服务的接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=44&api=17879&data=127&isNormal=1&vid=120">查询插件市场列表</a>和<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=44&api=18067&data=127&isNormal=1&vid=120">查询版本列表</a>获取可用版本。  */
	Values       string `json:"values,omitempty"`       /*  插件配置参数(YAML格式)，与valuesJson二选一。可通过容器镜像服务的接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=44&api=18132&data=127&isNormal=1&vid=120">查询版本values</a>获取values的模板。  */
	ValuesJson   string `json:"valuesJson,omitempty"`   /*  插件配置参数(JSON格式)，与values二选一。可通过容器镜像服务的接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=44&api=18132&data=127&isNormal=1&vid=120">查询版本values</a>获取values的模板。  */
}

type CcseDeployPluginInstanceResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求id  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  返回对象  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
