package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseUpdateClusterKubeConfigExpirationTimeApi
/* 调用该接口修改主账号下指定子账号的集群KubeConfig过期时间，需要使用主账号ak调用
 */type CcseUpdateClusterKubeConfigExpirationTimeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseUpdateClusterKubeConfigExpirationTimeApi(client *core.CtyunClient) *CcseUpdateClusterKubeConfigExpirationTimeApi {
	return &CcseUpdateClusterKubeConfigExpirationTimeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/certificate/kubeconfig/expire",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseUpdateClusterKubeConfigExpirationTimeApi) Do(ctx context.Context, credential core.Credential, req *CcseUpdateClusterKubeConfigExpirationTimeRequest) (*CcseUpdateClusterKubeConfigExpirationTimeResponse, error) {
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
	var resp CcseUpdateClusterKubeConfigExpirationTimeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseUpdateClusterKubeConfigExpirationTimeRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	UserId      int64 `json:"userId,omitempty"`      /*  用户ID，可在云容器引擎控制台 > 安全管理 > 授权页面查看主账号名下的子账号列表及对应的用户ID  */
	ValidPeriod int32 `json:"validPeriod,omitempty"` /*  证书过期秒数，最长为一年（31536000）  */
}

type CcseUpdateClusterKubeConfigExpirationTimeResponse struct {
	StatusCode int32                                                       `json:"statusCode,omitempty"` /*  响应状态码  */
	Message    string                                                      `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *CcseUpdateClusterKubeConfigExpirationTimeReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	Error      string                                                      `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type CcseUpdateClusterKubeConfigExpirationTimeReturnObjResponse struct{}
