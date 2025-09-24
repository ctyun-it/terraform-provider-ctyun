package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseCheckControlPlaneLogEnabledApi
/* 调用该接口查询核心组件日志采集开启情况
 */type CcseCheckControlPlaneLogEnabledApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseCheckControlPlaneLogEnabledApi(client *core.CtyunClient) *CcseCheckControlPlaneLogEnabledApi {
	return &CcseCheckControlPlaneLogEnabledApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/logcenter/controlplane/check",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseCheckControlPlaneLogEnabledApi) Do(ctx context.Context, credential core.Credential, req *CcseCheckControlPlaneLogEnabledRequest) (*CcseCheckControlPlaneLogEnabledResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseCheckControlPlaneLogEnabledResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseCheckControlPlaneLogEnabledRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcseCheckControlPlaneLogEnabledResponse struct {
	StatusCode int32    `json:"statusCode,omitempty"` /*  statusCode  */
	Message    string   `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  []string `json:"returnObj"`            /*  响应对象，已开通的核心组件日志，目前支持的核心组件有：apiserver、etcd、scheduler、controller-manager  */
	Error      string   `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}
