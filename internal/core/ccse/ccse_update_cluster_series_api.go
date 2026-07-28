package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseUpdateClusterSeriesApi
/* 调用该接口变更托管集群控制面规格。
 */type CcseUpdateClusterSeriesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseUpdateClusterSeriesApi(client *core.CtyunClient) *CcseUpdateClusterSeriesApi {
	return &CcseUpdateClusterSeriesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/v2/cce/clusters/{clusterId}/series",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseUpdateClusterSeriesApi) Do(ctx context.Context, credential core.Credential, req *CcseUpdateClusterSeriesRequest) (*CcseUpdateClusterSeriesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CcseUpdateClusterSeriesRequest
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
	var resp CcseUpdateClusterSeriesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseUpdateClusterSeriesRequest struct {
	ClusterId string `json:"clusterId,omitempty"` /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId  string `json:"regionId,omitempty"`  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	SeriesType string `json:"seriesType,omitempty"` /*  控制面规格，managedbase代表单实例，managedpro代表高可用  */
	NodeScale  string `json:"nodeScale,omitempty"`  /*  集群节点规模,managedbase中包含10节点规模，managedpro中包含50、200、1000、2000节点规模  */
}

type CcseUpdateClusterSeriesResponse struct {
	Code       int32                                     `json:"code"`       /*  请求结果编码  */
	RequestId  string                                    `json:"requestId"`  /*  请求id  */
	StatusCode int32                                     `json:"statusCode"` /*  状态码  */
	Message    string                                    `json:"message"`    /*  提示信息  */
	ReturnObj  *CcseUpdateClusterSeriesReturnObjResponse `json:"returnObj"`  /*  订单实例  */
	Error      string                                    `json:"error"`      /*  错误码  */
}

type CcseUpdateClusterSeriesReturnObjResponse struct {
	OrderId    string `json:"orderId"`    /*  订单id  */
	OrderNo    string `json:"orderNo"`    /*  订单编码  */
	ProdInstId string `json:"prodInstId"` /*  集群实例id  */
}
