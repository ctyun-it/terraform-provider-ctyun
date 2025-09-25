package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetTemplateInstanceApi
/* 调用该接口查看模板实例详情。
 */type CcseGetTemplateInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetTemplateInstanceApi(client *core.CtyunClient) *CcseGetTemplateInstanceApi {
	return &CcseGetTemplateInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/namespaces/{namespaceName}/templateinstance/{templateInstanceName}/detail",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetTemplateInstanceApi) Do(ctx context.Context, credential core.Credential, req *CcseGetTemplateInstanceRequest) (*CcseGetTemplateInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("namespaceName", req.NamespaceName)
	builder = builder.ReplaceUrl("templateInstanceName", req.TemplateInstanceName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetTemplateInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetTemplateInstanceRequest struct {
	ClusterId            string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	NamespaceName        string /*  命名空间名称  */
	TemplateInstanceName string /*  模板名称  */
	RegionId             string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcseGetTemplateInstanceResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  状态码  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求id  */
	Message    string                                    `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseGetTemplateInstanceReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetTemplateInstanceReturnObjResponse struct {
	ReleaseHistoryListDTOS []*CcseGetTemplateInstanceReturnObjReleaseHistoryListDTOSResponse `json:"releaseHistoryListDTOS"`    /*  发布版本列表  */
	FirstDeployTime        string                                                            `json:"firstDeployTime,omitempty"` /*  首次发布时间  */
	LastDeployTime         string                                                            `json:"lastDeployTime,omitempty"`  /*  最新发布时间  */
	Name                   string                                                            `json:"name,omitempty"`            /*  实例名称  */
	Namespace              string                                                            `json:"namespace,omitempty"`       /*  命名空间  */
	Values                 string                                                            `json:"values,omitempty"`          /*  values参数  */
	ResourceDTOS           []*CcseGetTemplateInstanceReturnObjResourceDTOSResponse           `json:"resourceDTOS"`              /*  资源列表  */
	ClusterId              string                                                            `json:"clusterId,omitempty"`       /*  集群ID  */
	CreatedTime            string                                                            `json:"createdTime,omitempty"`     /*  创建时间  */
	ChartVersion           string                                                            `json:"chartVersion,omitempty"`    /*  Chart版本  */
	Status                 string                                                            `json:"status,omitempty"`          /*  状态  */
	LastEvent              string                                                            `json:"lastEvent,omitempty"`       /*  历史事件  */
	Readme                 string                                                            `json:"readme,omitempty"`          /*  readme  */
	ExternalIps            []string                                                          `json:"externalIps"`               /*  externalIps  */
	ChartUrl               string                                                            `json:"chartUrl,omitempty"`        /*  Chart地址  */
	KubeConfigPath         string                                                            `json:"kubeConfigPath,omitempty"`  /*  kubeConfig路径  */
	RepositoryId           string                                                            `json:"repositoryId,omitempty"`    /*  镜像实例仓库id  */
	ChartName              string                                                            `json:"chartName,omitempty"`       /*  chart名称  */
	Icon                   string                                                            `json:"icon,omitempty"`            /*  icon  */
}

type CcseGetTemplateInstanceReturnObjReleaseHistoryListDTOSResponse struct {
	Revision     string `json:"revision,omitempty"`     /*  版本  */
	Updated      string `json:"updated,omitempty"`      /*  更新时间  */
	Status       string `json:"status,omitempty"`       /*  状态  */
	Chart        string `json:"chart,omitempty"`        /*  chart名称和版本  */
	AppVersion   string `json:"appVersion,omitempty"`   /*  app版本  */
	Description  string `json:"description,omitempty"`  /*  描述  */
	ClusterId    string `json:"clusterId,omitempty"`    /*  集群id  */
	Name         string `json:"name,omitempty"`         /*  实例名称  */
	Namespace    string `json:"namespace,omitempty"`    /*  命名空间  */
	ChartVersion string `json:"chartVersion,omitempty"` /*  chart版本  */
}

type CcseGetTemplateInstanceReturnObjResourceDTOSResponse struct {
	Name      string `json:"name,omitempty"`      /*  名称  */
	Namespace string `json:"namespace,omitempty"` /*  命名空间  */
	Kind      string `json:"kind,omitempty"`      /*  资源类型  */
	Yaml      string `json:"yaml,omitempty"`      /*  YAML内容  */
}
