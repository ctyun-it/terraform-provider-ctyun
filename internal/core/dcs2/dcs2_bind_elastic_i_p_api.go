package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2BindElasticIPApi
/* 可将弹性IP（Elastic IP Address，简称EIP）与分布式缓存Redis实例绑定。
 */type Dcs2BindElasticIPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2BindElasticIPApi(client *core.CtyunClient) *Dcs2BindElasticIPApi {
	return &Dcs2BindElasticIPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/eip/bindElasticIP",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2BindElasticIPApi) Do(ctx context.Context, credential core.Credential, req *Dcs2BindElasticIPRequest) (*Dcs2BindElasticIPResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
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
	var resp Dcs2BindElasticIPResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2BindElasticIPRequest struct {
	RegionId       string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	BindObjType    int64  `json:"bindObjType,omitempty"`    /*  绑定对象类型<li>4：VIP  */
	ProdInstId     string `json:"prodInstId,omitempty"`     /*  实例标识  */
	ElasticIp      string `json:"elasticIp,omitempty"`      /*  弹性IP，您可以查看<a href="https://www.ctyun.cn/document/10026753/10026909">产品定义-弹性IP</a>来了解弹性IP<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8652&vid=88">新查询指定地域已创建的EIP</a> eipAddress字段。<br><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=5723&isNormal=1&vid=88">创建EIP</a>  */
	BindByPaasProd int32  `json:"bindByPaasProd,omitempty"` /*  由组件侧绑定弹性IP：<br><li>1：是<br><li>0：否（默认值）<br><br>适用场景：<br><li>仅桌面云ELB绑定弹性IP需传1<br><li>其他场景传0  */
}

type Dcs2BindElasticIPResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2BindElasticIPReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2BindElasticIPReturnObjResponse struct{}
