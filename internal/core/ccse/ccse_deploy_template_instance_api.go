package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseDeployTemplateInstanceApi
/* 调用该接口创建模板实例。
 */type CcseDeployTemplateInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseDeployTemplateInstanceApi(client *core.CtyunClient) *CcseDeployTemplateInstanceApi {
	return &CcseDeployTemplateInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/namespaces/{namespaceName}/templateinstance",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseDeployTemplateInstanceApi) Do(ctx context.Context, credential core.Credential, req *CcseDeployTemplateInstanceRequest) (*CcseDeployTemplateInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("namespaceName", req.NamespaceName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseDeployTemplateInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseDeployTemplateInstanceRequest struct {
	ClusterId     string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	NamespaceName string /*  命名空间名称  */
	RegionId      string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	ChartName     string `json:"chartName,omitempty"`     /*  Chart名称  */
	ChartVersion  string `json:"chartVersion,omitempty"`  /*  Chart版本  */
	CrNamespaceId int64  `json:"crNamespaceId,omitempty"` /*  模板存放在镜像服务实例中的命名空间的ID，获取方法参见<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=44&api=10629&data=127&isNormal=1&vid=120" target="_blank">查询HelmChart命名空间列表</a>

	 */
	InstanceName      string `json:"instanceName,omitempty"`      /*  实例名  */
	InstanceValue     string `json:"instanceValue,omitempty"`     /*  实例参数  */
	InstanceValueType string `json:"instanceValueType,omitempty"` /*  实例参数  */
	RepositoryId      int64  `json:"repositoryId,omitempty"`      /*  镜像服务实例仓库ID，获取方法参见<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=44&api=10621&data=127&isNormal=1&vid=120" target="_blank">查询镜像仓库列表</a>  */
	Timeout           string `json:"timeout,omitempty"`           /*  实例安装超时时间  */
	IsSyncMode        *bool  `json:"isSyncMode"`                  /*  是否同步调用  */
}

type CcseDeployTemplateInstanceResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求id  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  返回对象  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
