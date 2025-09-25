package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2UnBindElasticIPApi
/* 可将弹性IP（Elastic IP Address，简称EIP）与分布式缓存Redis实例解绑。
 */type Dcs2UnBindElasticIPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2UnBindElasticIPApi(client *core.CtyunClient) *Dcs2UnBindElasticIPApi {
	return &Dcs2UnBindElasticIPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/eip/unBindElasticIP",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2UnBindElasticIPApi) Do(ctx context.Context, credential core.Credential, req *Dcs2UnBindElasticIPRequest) (*Dcs2UnBindElasticIPResponse, error) {
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
	var resp Dcs2UnBindElasticIPResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2UnBindElasticIPRequest struct {
	RegionId       string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	BindObjType    int64  `json:"bindObjType,omitempty"`    /*  绑定对象类型<li>4：VIP  */
	ProdInstId     string `json:"prodInstId,omitempty"`     /*  实例ID  */
	ElasticIp      string `json:"elasticIp,omitempty"`      /*  弹性IP<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7739&isNormal=1&vid=270">查询实例网络信息</a> elasticIp字段获取当前实例绑定的弹性IP  */
	BindByPaasProd int32  `json:"bindByPaasProd,omitempty"` /*  由组件侧绑定弹性IP：<br><li>1：是<br><li>0：否（默认值）<br><br>适用场景：<br><li>仅桌面云ELB绑定弹性IP需传1<br><li>其他场景传0  */
}

type Dcs2UnBindElasticIPResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2UnBindElasticIPReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2UnBindElasticIPReturnObjResponse struct{}
