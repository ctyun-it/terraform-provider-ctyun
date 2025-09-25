package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CcseGetClusterKubeConfigApi
/* 调用该接口查看当前ak对应的用户的集群KubeConfig。
 */type CcseGetClusterKubeConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetClusterKubeConfigApi(client *core.CtyunClient) *CcseGetClusterKubeConfigApi {
	return &CcseGetClusterKubeConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/certificate/kubeconfig",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetClusterKubeConfigApi) Do(ctx context.Context, credential core.Credential, req *CcseGetClusterKubeConfigRequest) (*CcseGetClusterKubeConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("isPublic", strconv.FormatBool(req.IsPublic))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetClusterKubeConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetClusterKubeConfigRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	IsPublic bool /*  是否获取公网KubeConfig  */
}

type CcseGetClusterKubeConfigResponse struct {
	StatusCode int32                                      `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                     `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseGetClusterKubeConfigReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                     `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetClusterKubeConfigReturnObjResponse struct {
	ExpireDate string `json:"expireDate,omitempty"` /*  KubeConfig过期时间  */
	KubeConfig string `json:"kubeConfig,omitempty"` /*  KubeConfig  */
}
