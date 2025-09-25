package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseUpdateClusterApi
/* 调用该接口修改集群。
 */type CcseUpdateClusterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseUpdateClusterApi(client *core.CtyunClient) *CcseUpdateClusterApi {
	return &CcseUpdateClusterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/v2/cce/clusters/{clusterId}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseUpdateClusterApi) Do(ctx context.Context, credential core.Credential, req *CcseUpdateClusterRequest) (*CcseUpdateClusterResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
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
	var resp CcseUpdateClusterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseUpdateClusterRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	ClusterDesc     string                             `json:"clusterDesc,omitempty"`     /*  集群描述  */
	ClusterAlias    string                             `json:"clusterAlias,omitempty"`    /*  集群显示名称  */
	StartPort       int32                              `json:"startPort,omitempty"`       /*  服务起始端口，取值范围[20106, 32767]  */
	EndPort         int32                              `json:"endPort,omitempty"`         /*  服务结束端口，取值范围[20106, 32767]  */
	SecurityGroupId string                             `json:"securityGroupId,omitempty"` /*  节点默认安全组ID  */
	CustomSan       *CcseUpdateClusterCustomSanRequest `json:"customSan"`                 /*  自定义SAN列表  */
	Cubecni         *CcseUpdateClusterCubecniRequest   `json:"cubecni"`                   /*  cubecni插件网络配置  */
}

type CcseUpdateClusterCustomSanRequest struct {
	Action string   `json:"action,omitempty"` /*  操作类型，支持：overwrite  */
	Values []string `json:"values"`           /*  SAN列表  */
}

type CcseUpdateClusterCubecniRequest struct {
	MinPoolSize   int32    `json:"minPoolSize,omitempty"` /*  最小缓存辅助IP数，取值范围[0, 60]，不大于maxPoolSize  */
	MaxPoolSize   int32    `json:"maxPoolSize,omitempty"` /*  最大缓存辅助IP数，取值范围[0, 60]，不小于minPoolSize  */
	AppendSubnets []string `json:"appendSubnets"`         /*  新增Pod子网ID列表  */
}

type CcseUpdateClusterResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  返回对象  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
