package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetPluginInstanceApi
/* 调用该接口可查询插件实例的详细信息。
 */type CcseGetPluginInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetPluginInstanceApi(client *core.CtyunClient) *CcseGetPluginInstanceApi {
	return &CcseGetPluginInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/plugininstance/{pluginName}/detail",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetPluginInstanceApi) Do(ctx context.Context, credential core.Credential, req *CcseGetPluginInstanceRequest) (*CcseGetPluginInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("pluginName", req.PluginName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("namespace", req.Namespace)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetPluginInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetPluginInstanceRequest struct {
	ClusterId  string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	PluginName string /*  插件名称  */
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	Namespace string /*  插件所在命名空间名称  */
}

type CcseGetPluginInstanceResponse struct {
	StatusCode int32                                   `json:"statusCode,omitempty"` /*  状态码  */
	RequestId  string                                  `json:"requestId,omitempty"`  /*  请求id  */
	Message    string                                  `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseGetPluginInstanceReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetPluginInstanceReturnObjResponse struct {
	ReleaseHistoryListDTOS []*CcseGetPluginInstanceReturnObjReleaseHistoryListDTOSResponse `json:"releaseHistoryListDTOS"`    /*  插件实例发布历史版本列表  */
	FirstDeployTime        string                                                          `json:"firstDeployTime,omitempty"` /*  首次发布时间。  */
	LastDeployTime         string                                                          `json:"lastDeployTime,omitempty"`  /*  最新发布时间。  */
	Name                   string                                                          `json:"name,omitempty"`            /*  插件实例名称。  */
	Namespace              string                                                          `json:"namespace,omitempty"`       /*  命名空间。  */
	Values                 string                                                          `json:"values,omitempty"`          /*  values参数。  */
	ResourceDTOS           []*CcseGetPluginInstanceReturnObjResourceDTOSResponse           `json:"resourceDTOS"`              /*  资源列表。  */
	ClusterName            string                                                          `json:"clusterName,omitempty"`     /*  集群ID  */
	CreatedTime            string                                                          `json:"createdTime,omitempty"`     /*  创建时间  */
	ChartVersion           string                                                          `json:"chartVersion,omitempty"`    /*  chart版本  */
	Status                 string                                                          `json:"status,omitempty"`          /*  状态  */
	Readme                 string                                                          `json:"readme,omitempty"`          /*  说明  */
	ChartUrl               string                                                          `json:"chartUrl,omitempty"`        /*  chart url  */
	RepositoryId           int64                                                           `json:"repositoryId,omitempty"`    /*  仓库ID  */
	ChartName              string                                                          `json:"chartName,omitempty"`       /*  chart名称  */
	Icon                   string                                                          `json:"icon,omitempty"`            /*  icon名称  */
}

type CcseGetPluginInstanceReturnObjReleaseHistoryListDTOSResponse struct {
	Revision     string `json:"revision,omitempty"`     /*  版本  */
	Updated      string `json:"updated,omitempty"`      /*  更新时间。  */
	Status       string `json:"status,omitempty"`       /*  状态。  */
	Chart        string `json:"chart,omitempty"`        /*  Chart名称和版本。  */
	AppVersion   string `json:"appVersion,omitempty"`   /*  版本。  */
	Description  string `json:"description,omitempty"`  /*  描述。  */
	ClusterId    string `json:"clusterId,omitempty"`    /*  集群ID。  */
	Name         string `json:"name,omitempty"`         /*  实例名称。  */
	Namespace    string `json:"namespace,omitempty"`    /*  命名空间。  */
	ChartName    string `json:"chartName,omitempty"`    /*  Chart名称。  */
	ChartVersion string `json:"chartVersion,omitempty"` /*  Chart版本。  */
}

type CcseGetPluginInstanceReturnObjResourceDTOSResponse struct {
	Name      string `json:"name,omitempty"`      /*  资源名称。  */
	Namespace string `json:"namespace,omitempty"` /*  资源命名空间。  */
	Kind      string `json:"kind,omitempty"`      /*  资源kind。  */
	Yaml      string `json:"yaml,omitempty"`      /*  资源内容。  */
}
