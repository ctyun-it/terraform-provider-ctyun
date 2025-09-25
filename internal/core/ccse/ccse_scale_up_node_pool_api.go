package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseScaleUpNodePoolApi
/* 调用该接口扩容节点池中的节点。
 */type CcseScaleUpNodePoolApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseScaleUpNodePoolApi(client *core.CtyunClient) *CcseScaleUpNodePoolApi {
	return &CcseScaleUpNodePoolApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepool/{nodePoolId}/scaleup",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseScaleUpNodePoolApi) Do(ctx context.Context, credential core.Credential, req *CcseScaleUpNodePoolRequest) (*CcseScaleUpNodePoolResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("nodePoolId", req.NodePoolId)
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
	var resp CcseScaleUpNodePoolResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseScaleUpNodePoolRequest struct {
	ClusterId  string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	NodePoolId string /*  节点池ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	Num int32 `json:"num,omitempty"` /*  扩容节点规模  */
}

type CcseScaleUpNodePoolResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  string `json:"returnObj,omitempty"`  /*  返回对象，伸缩记录的id，可通过伸缩记录详情接口查看伸缩状态等信息  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
